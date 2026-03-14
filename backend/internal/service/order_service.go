package service

import (
	"context"
	"fmt"
	"time"

	"github.com/x-store/backend/internal/config"
	"github.com/x-store/backend/internal/model"
	"github.com/x-store/backend/internal/repo"
	"github.com/x-store/backend/pkg/crypto"
	"github.com/x-store/backend/pkg/email"
)

type OrderService struct {
	orderRepo    *repo.OrderRepo
	cardKeyRepo  *repo.CardKeyRepo
	productRepo  *repo.ProductRepo
	stockMgr     *crypto.StockManager
	emailService *email.EmailService
}

func NewOrderService() *OrderService {
	return &OrderService{
		orderRepo:    repo.NewOrderRepo(),
		cardKeyRepo:  repo.NewCardKeyRepo(),
		productRepo:  repo.NewProductRepo(),
		stockMgr:     crypto.NewStockManager(config.RDB),
		emailService: email.NewEmailService(&config.Global.Email),
	}
}

type CreateOrderReq struct {
	ProductID uint   `json:"product_id" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	UserID    *uint  `json:"user_id"` // 可选，游客下单时为空
}

type CreateOrderResp struct {
	OrderNo  string    `json:"order_no"`
	Amount   float64   `json:"amount"`
	ExpireAt time.Time `json:"expire_at"`
}

// CreateOrder 创建订单（锁定库存 + 锁定卡密）
func (s *OrderService) CreateOrder(ctx context.Context, req CreateOrderReq) (*CreateOrderResp, error) {
	// 1. 查询商品
	product, err := s.productRepo.GetByID(req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("商品不存在")
	}
	if product.Status != 1 {
		return nil, fmt.Errorf("商品已下架")
	}

	// 2. Redis 扣减库存（原子性防超卖）
	if err := s.stockMgr.DeductStock(ctx, req.ProductID, 1); err != nil {
		return nil, err
	}

	// 3. 数据库锁定卡密（仅自动发货商品）
	var lockedKeys []model.CardKey
	if product.DeliveryType == "auto" {
		lockedKeys, err = s.cardKeyRepo.LockAvailable(req.ProductID, 1)
		if err != nil {
			// 回滚 Redis 库存
			s.stockMgr.ReleaseStock(ctx, req.ProductID, 1)
			return nil, fmt.Errorf("卡密库存不足")
		}
	}

	// 4. 创建订单
	orderNo := s.generateOrderNo()
	expireAt := time.Now().Add(15 * time.Minute) // 15分钟锁定期

	order := &model.Order{
		OrderNo:   orderNo,
		UserID:    req.UserID,
		ProductID: req.ProductID,
		Email:     req.Email,
		Amount:    product.Price,
		Status:    model.OrderStatusPending,
		ExpireAt:  expireAt,
	}

	if err := s.orderRepo.Create(order); err != nil {
		// 回滚库存和卡密
		s.stockMgr.ReleaseStock(ctx, req.ProductID, 1)
		if len(lockedKeys) > 0 {
			ids := make([]uint, len(lockedKeys))
			for i, k := range lockedKeys {
				ids[i] = k.ID
			}
			s.cardKeyRepo.ReleaseLockedKeys(ids)
		}
		return nil, fmt.Errorf("创建订单失败: %w", err)
	}

	return &CreateOrderResp{
		OrderNo:  orderNo,
		Amount:   order.Amount,
		ExpireAt: expireAt,
	}, nil
}

// GetByOrderNo 根据订单号查询
func (s *OrderService) GetByOrderNo(orderNo string) (*model.Order, error) {
	return s.orderRepo.GetByOrderNo(orderNo)
}

// PayOrder 支付订单（模拟支付成功）
func (s *OrderService) PayOrder(ctx context.Context, orderNo string, payMethod string) error {
	order, err := s.orderRepo.GetByOrderNo(orderNo)
	if err != nil {
		return fmt.Errorf("订单不存在")
	}

	if order.Status != model.OrderStatusPending {
		return fmt.Errorf("订单状态异常")
	}

	if time.Now().After(order.ExpireAt) {
		return fmt.Errorf("订单已过期")
	}

	// 更新订单状态
	now := time.Now()
	order.Status = model.OrderStatusPaid
	order.PayMethod = payMethod
	order.PaidAt = &now

	if err := s.orderRepo.Update(order); err != nil {
		return err
	}

	// 自动发货
	if err := s.DeliverOrder(ctx, order.ID); err != nil {
		return fmt.Errorf("发货失败: %w", err)
	}

	return nil
}

// DeliverOrder 发货（标记卡密为已售出）
func (s *OrderService) DeliverOrder(ctx context.Context, orderID uint) error {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return err
	}

	// 获取锁定的卡密
	keys, err := s.cardKeyRepo.GetByOrderID(orderID)
	if err != nil || len(keys) == 0 {
		// 如果没有预锁定的卡密，现在分配
		product, _ := s.productRepo.GetByID(order.ProductID)
		if product != nil && product.DeliveryType == "auto" {
			keys, err = s.cardKeyRepo.LockAvailable(order.ProductID, 1)
			if err != nil {
				return fmt.Errorf("卡密不足")
			}
		}
	}

	// 标记为已售出
	if len(keys) > 0 {
		ids := make([]uint, len(keys))
		for i, k := range keys {
			ids[i] = k.ID
		}
		if err := s.cardKeyRepo.MarkAsSold(ids, orderID); err != nil {
			return err
		}
	}

	// 更新订单状态
	order.Status = model.OrderStatusDelivered
	if err := s.orderRepo.Update(order); err != nil {
		return err
	}

	// 发送邮件通知
	go s.sendOrderEmail(order, keys)

	return nil
}

// CancelOrder 取消订单（释放库存和卡密）
func (s *OrderService) CancelOrder(ctx context.Context, orderNo string) error {
	order, err := s.orderRepo.GetByOrderNo(orderNo)
	if err != nil {
		return err
	}

	if order.Status != model.OrderStatusPending {
		return fmt.Errorf("订单状态不允许取消")
	}

	// 释放 Redis 库存
	s.stockMgr.ReleaseStock(ctx, order.ProductID, 1)

	// 释放锁定的卡密
	keys, _ := s.cardKeyRepo.GetByOrderID(order.ID)
	if len(keys) > 0 {
		ids := make([]uint, len(keys))
		for i, k := range keys {
			ids[i] = k.ID
		}
		s.cardKeyRepo.ReleaseLockedKeys(ids)
	}

	// 更新订单状态
	order.Status = model.OrderStatusCancelled
	return s.orderRepo.Update(order)
}

// ProcessExpiredOrders 处理过期订单（定时任务调用）
func (s *OrderService) ProcessExpiredOrders(ctx context.Context) error {
	orders, err := s.orderRepo.ListExpired(100)
	if err != nil {
		return err
	}

	for _, order := range orders {
		s.CancelOrder(ctx, order.OrderNo)
	}
	return nil
}

// ListAll 管理后台获取全部订单
func (s *OrderService) ListAll(page, size int, status *int, keyword string) ([]model.Order, int64, error) {
	return s.orderRepo.ListAll(page, size, status, keyword)
}

// RefundOrder 退款订单（释放卡密、更新状态）
func (s *OrderService) RefundOrder(ctx context.Context, orderNo string) error {
	order, err := s.orderRepo.GetByOrderNo(orderNo)
	if err != nil {
		return fmt.Errorf("订单不存在")
	}

	// 只有已支付或已发货的订单才能退款
	if order.Status != model.OrderStatusPaid && order.Status != model.OrderStatusDelivered {
		return fmt.Errorf("订单状态不允许退款")
	}

	// 释放卡密（如果已发货）
	if order.Status == model.OrderStatusDelivered {
		keys, _ := s.cardKeyRepo.GetByOrderID(order.ID)
		if len(keys) > 0 {
			ids := make([]uint, len(keys))
			for i, k := range keys {
				ids[i] = k.ID
			}
			// 将已售出的卡密释放回可用状态
			s.cardKeyRepo.ReleaseLockedKeys(ids)
		}
	}

	// 释放 Redis 库存
	s.stockMgr.ReleaseStock(ctx, order.ProductID, 1)

	// 更新订单状态为已退款
	order.Status = model.OrderStatusCancelled // 使用取消状态表示退款
	return s.orderRepo.Update(order)
}

// sendOrderEmail 发送订单邮件通知
func (s *OrderService) sendOrderEmail(order *model.Order, keys []model.CardKey) {
	product, _ := s.productRepo.GetByID(order.ProductID)
	if product == nil {
		return
	}

	cardKeyContents := make([]string, len(keys))
	for i, k := range keys {
		cardKeyContents[i] = k.Content
	}

	emailData := email.OrderEmailData{
		OrderNo:      order.OrderNo,
		Email:        order.Email,
		ProductTitle: product.Title,
		Amount:       order.Amount,
		CardKeys:     cardKeyContents,
		CreatedAt:    order.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if err := s.emailService.SendOrderNotification(emailData); err != nil {
		// 邮件发送失败不影响主流程，仅记录日志
		fmt.Printf("[Email] Failed to send order notification: %v\n", err)
	}
}

// generateOrderNo 生成订单号 XS20260314001
func (s *OrderService) generateOrderNo() string {
	return fmt.Sprintf("XS%s%03d", time.Now().Format("20060102"), time.Now().UnixNano()%1000)
}
