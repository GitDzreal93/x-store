package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/x-store/backend/internal/model"
	"github.com/x-store/backend/internal/repo"
	"github.com/x-store/backend/pkg/payment"
)

type PaymentService struct {
	paymentRepo *repo.PaymentRepo
	channelRepo *repo.PaymentChannelRepo
	orderSvc    *OrderService
}

func NewPaymentService(orderSvc *OrderService) *PaymentService {
	return &PaymentService{
		paymentRepo: repo.NewPaymentRepo(),
		channelRepo: repo.NewPaymentChannelRepo(),
		orderSvc:    orderSvc,
	}
}

type CreatePaymentReq struct {
	OrderNo   string `json:"order_no" binding:"required"`
	ChannelID uint   `json:"channel_id" binding:"required"`
}

type CreatePaymentResp struct {
	PaymentID uint   `json:"payment_id"`
	PayURL    string `json:"pay_url"`
	QRCode    string `json:"qr_code"`
	TradeNo   string `json:"trade_no"`
	ExpireAt  int64  `json:"expire_at"` // Unix 时间戳
}

// CreatePayment 创建支付单
func (s *PaymentService) CreatePayment(req CreatePaymentReq) (*CreatePaymentResp, error) {
	// 1. 查询订单
	order, err := s.orderSvc.GetByOrderNo(req.OrderNo)
	if err != nil {
		return nil, fmt.Errorf("订单不存在")
	}
	if order.Status != model.OrderStatusPending {
		return nil, fmt.Errorf("订单状态异常")
	}
	if time.Now().After(order.ExpireAt) {
		return nil, fmt.Errorf("订单已过期")
	}

	// 2. 查询支付渠道
	channel, err := s.channelRepo.GetByID(req.ChannelID)
	if err != nil || channel == nil {
		return nil, fmt.Errorf("支付渠道不存在")
	}
	if !channel.IsActive {
		return nil, fmt.Errorf("支付渠道已禁用")
	}

	// 3. 检查是否已有待支付的支付单（避免重复创建）
	existingPayment, _ := s.paymentRepo.GetLatestByOrderID(order.ID)
	if existingPayment != nil && existingPayment.Status == model.PaymentStatusProcessing {
		// 复用现有支付单
		return &CreatePaymentResp{
			PaymentID: existingPayment.ID,
			PayURL:    existingPayment.RawNotify, // 临时存储在这里
			QRCode:    existingPayment.RawNotify,
			TradeNo:   existingPayment.TradeNo,
			ExpireAt:  order.ExpireAt.Unix(),
		}, nil
	}

	// 4. 计算手续费
	feeAmount := order.Amount * channel.FeeRate / 100
	if channel.FixedFee > 0 {
		feeAmount += channel.FixedFee
	}
	totalAmount := order.Amount + feeAmount

	// 5. 创建支付记录
	now := time.Now()
	paymentRecord := &model.Payment{
		OrderID:      order.ID,
		TradeNo:      "",
		PayMethod:    channel.ChannelType,
		Amount:       totalAmount,
		Status:       model.PaymentStatusProcessing,
		RawNotify:    "",
		CompletedAt:  nil,
		BaseModel:    model.BaseModel{CreatedAt: now, UpdatedAt: now},
	}

	if err := s.paymentRepo.Create(paymentRecord); err != nil {
		return nil, fmt.Errorf("创建支付记录失败")
	}

	// 6. 创建支付提供商并调用
	provider, err := s.createProvider(channel)
	if err != nil {
		return nil, err
	}

	input := payment.CreateInput{
		OrderNo:   order.OrderNo,
		PaymentID: paymentRecord.ID,
		Amount:    paymentRecord.Amount,
		Subject:   fmt.Sprintf("订单 %s", order.OrderNo),
		NotifyURL: "", // 使用渠道配置中的默认值
		ReturnURL: "",
	}

	result, err := provider.CreatePayment(input)
	if err != nil {
		return nil, fmt.Errorf("创建支付失败: %w", err)
	}

	// 7. 更新支付记录
	rawNotifyJSON, _ := json.Marshal(result.RawResponse)
	paymentRecord.TradeNo = result.TradeNo
	paymentRecord.RawNotify = string(rawNotifyJSON)
	paymentRecord.UpdatedAt = time.Now()
	if err := s.paymentRepo.Update(paymentRecord); err != nil {
		return nil, fmt.Errorf("更新支付记录失败")
	}

	return &CreatePaymentResp{
		PaymentID: paymentRecord.ID,
		PayURL:    result.PayURL,
		QRCode:    result.QRCode,
		TradeNo:   result.TradeNo,
		ExpireAt:  order.ExpireAt.Unix(),
	}, nil
}

// createProvider 根据渠道配置创建支付提供商
func (s *PaymentService) createProvider(channel *model.PaymentChannel) (payment.Provider, error) {
	switch channel.ProviderType {
	case model.ProviderTypeMock:
		return payment.NewMockPayProviderWrapper(channel.ConfigJSON)

	case model.ProviderTypeStripe:
		return payment.NewStripeProviderWrapper(channel.ConfigJSON)

	case model.ProviderTypeAlipay:
		return payment.NewAlipayProvider(channel.ConfigJSON, channel.ChannelType)

	case model.ProviderTypeWechat:
		return payment.NewWechatPayProvider(channel.ConfigJSON, channel.ChannelType)

	case model.ProviderTypePayPal:
		return payment.NewPayPalProvider(channel.ConfigJSON)

	case model.ProviderTypeEpusdt:
		return payment.NewEpusdtProvider(channel.ConfigJSON)

	case model.ProviderTypeYipay:
		return payment.NewYipayProvider(channel.ConfigJSON, channel.ChannelType)

	case model.ProviderTypePayjs:
		return payment.NewPayjsProvider(channel.ConfigJSON)

	case model.ProviderTypeXunhupay:
		return payment.NewXunhupayProvider(channel.ConfigJSON, channel.ChannelType)

	case model.ProviderTypeVmqpay:
		return payment.NewVmqpayProvider(channel.ConfigJSON, channel.ChannelType)

	default:
		return nil, fmt.Errorf("不支持的支付方式: %s", channel.ProviderType)
	}
}

// createMockPayment 创建模拟支付
func (s *PaymentService) createMockPayment(channel *model.PaymentChannel, order *model.Order, paymentRecord *model.Payment) (*payment.MockPayCreateResult, error) {
	// 解析配置
	var cfg payment.MockPayConfig
	if err := json.Unmarshal([]byte(channel.ConfigJSON), &cfg); err != nil {
		return nil, fmt.Errorf("支付渠道配置错误")
	}

	// 调用模拟支付网关
	result, err := payment.CreateMockPayment(cfg, payment.MockPayCreateInput{
		OrderNo:   order.OrderNo,
		PaymentID: paymentRecord.ID,
		Amount:    paymentRecord.Amount,
		Subject:   fmt.Sprintf("订单 %s", order.OrderNo),
	})
	if err != nil {
		return nil, fmt.Errorf("创建模拟支付失败: %w", err)
	}

	return result, nil
}

// HandleCallback 处理支付回调
type PaymentCallbackReq struct {
	TradeNo   string  `json:"trade_no"`
	OrderNo   string  `json:"order_no"`
	PaymentID uint    `json:"payment_id"`
	Amount    float64 `json:"amount"`
	Status    string  `json:"status"`
	PaidAt    string  `json:"paid_at"`
	Sign      string  `json:"sign"`
}

func (s *PaymentService) HandleCallback(channelID uint, req PaymentCallbackReq) error {
	// 1. 查询支付渠道
	channel, err := s.channelRepo.GetByID(channelID)
	if err != nil || channel == nil {
		return fmt.Errorf("支付渠道不存在")
	}

	// 2. 验证签名
	if channel.ProviderType == model.ProviderTypeMock {
		var cfg payment.MockPayConfig
		if err := json.Unmarshal([]byte(channel.ConfigJSON), &cfg); err != nil {
			return fmt.Errorf("支付渠道配置错误")
		}

		notifyData := &payment.MockPayNotifyData{
			TradeNo:   req.TradeNo,
			OrderNo:   req.OrderNo,
			PaymentID: req.PaymentID,
			Amount:    req.Amount,
			Status:    req.Status,
			PaidAt:    req.PaidAt,
			Sign:      req.Sign,
		}

		if !payment.VerifyMockPaySign(cfg.Secret, notifyData) {
			return fmt.Errorf("签名验证失败")
		}
	}

	// 3. 查询支付记录
	paymentRecord, err := s.paymentRepo.GetByID(req.PaymentID)
	if err != nil || paymentRecord == nil {
		return fmt.Errorf("支付记录不存在")
	}

	// 4. 幂等处理：已成功的不再处理
	if paymentRecord.Status == model.PaymentStatusSuccess {
		return nil
	}

	// 5. 验证金额
	if req.Amount != paymentRecord.Amount {
		return fmt.Errorf("支付金额不匹配")
	}

	// 6. 更新支付状态
	status := payment.ParseMockPayStatus(req.Status)
	paymentRecord.Status = mapPaymentStatus(status)
	
	if status == "success" {
		paidAt, _ := time.Parse(time.RFC3339, req.PaidAt)
		paymentRecord.CompletedAt = &paidAt
		
		// 7. 更新订单状态并发货
		if err := s.orderSvc.PayOrder(nil, req.OrderNo, channel.ChannelType); err != nil {
			return fmt.Errorf("订单支付处理失败: %w", err)
		}
	}

	paymentRecord.UpdatedAt = time.Now()
	if err := s.paymentRepo.Update(paymentRecord); err != nil {
		return fmt.Errorf("更新支付记录失败")
	}

	return nil
}

// HandleWebhook 统一处理所有渠道的支付回调（使用 Provider 接口）
func (s *PaymentService) HandleWebhook(channelID uint, body []byte, headers map[string]string) error {
	// 1. 查询支付渠道
	channel, err := s.channelRepo.GetByID(channelID)
	if err != nil || channel == nil {
		return fmt.Errorf("支付渠道不存在")
	}

	// 2. 创建 Provider 并验证回调
	provider, err := s.createProvider(channel)
	if err != nil {
		return fmt.Errorf("创建支付提供商失败: %w", err)
	}

	result, err := provider.VerifyNotify(body, headers)
	if err != nil {
		return fmt.Errorf("回调验证失败: %w", err)
	}

	if result.Status != "success" {
		return fmt.Errorf("支付未成功: %s", result.Status)
	}

	// 3. 更新订单状态并发货
	if err := s.orderSvc.PayOrder(nil, result.OrderNo, channel.ChannelType); err != nil {
		return fmt.Errorf("订单支付处理失败: %w", err)
	}

	return nil
}

// GetPaymentStatus 查询支付状态
func (s *PaymentService) GetPaymentStatus(paymentID uint) (*model.Payment, error) {
	return s.paymentRepo.GetByID(paymentID)
}

// mapPaymentStatus 映射支付状态
func mapPaymentStatus(status string) int {
	switch status {
	case "success":
		return model.PaymentStatusSuccess
	case "failed":
		return model.PaymentStatusFailed
	default:
		return model.PaymentStatusProcessing
	}
}

// ListChannels 获取可用支付渠道列表
func (s *PaymentService) ListChannels() ([]model.PaymentChannel, error) {
	active := true
	return s.channelRepo.List(&active)
}
