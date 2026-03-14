package repo

import (
	"time"

	"github.com/x-store/backend/internal/config"
	"github.com/x-store/backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderRepo struct {
	db *gorm.DB
}

func NewOrderRepo() *OrderRepo {
	return &OrderRepo{db: config.DB}
}

func (r *OrderRepo) Create(order *model.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepo) GetByID(id uint) (*model.Order, error) {
	var order model.Order
	if err := r.db.Preload("Product").Preload("CardKeys").First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepo) GetByOrderNo(orderNo string) (*model.Order, error) {
	var order model.Order
	if err := r.db.Preload("Product").Preload("CardKeys").Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepo) Update(order *model.Order) error {
	return r.db.Save(order).Error
}

// UpdateStatus 更新订单状态
func (r *OrderRepo) UpdateStatus(id uint, status int) error {
	return r.db.Model(&model.Order{}).Where("id = ?", id).Update("status", status).Error
}

// ListByUser 获取用户订单列表
func (r *OrderRepo) ListByUser(userID uint, page, size int) ([]model.Order, int64, error) {
	var list []model.Order
	var total int64

	q := r.db.Model(&model.Order{}).Where("user_id = ?", userID).Preload("Product")
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := q.Order("id DESC").Offset(offset).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// ListAll 管理后台获取全部订单（分页+筛选）
func (r *OrderRepo) ListAll(page, size int, status *int, keyword string) ([]model.Order, int64, error) {
	var list []model.Order
	var total int64

	q := r.db.Model(&model.Order{}).Preload("Product").Preload("CardKeys")
	if status != nil {
		q = q.Where("status = ?", *status)
	}
	if keyword != "" {
		q = q.Where("order_no LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * size
	if err := q.Order("id DESC").Offset(offset).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// ListExpired 获取已过期未支付的订单
func (r *OrderRepo) ListExpired(limit int) ([]model.Order, error) {
	var list []model.Order
	err := r.db.Where("status = ? AND expire_at < ?", model.OrderStatusPending, time.Now()).
		Limit(limit).
		Find(&list).Error
	return list, err
}

// --- CardKey Repo ---

type CardKeyRepo struct {
	db *gorm.DB
}

func NewCardKeyRepo() *CardKeyRepo {
	return &CardKeyRepo{db: config.DB}
}

// BatchCreate 批量创建卡密
func (r *CardKeyRepo) BatchCreate(keys []model.CardKey) error {
	return r.db.Create(&keys).Error
}

// LockAvailable 锁定可用卡密（预占）
func (r *CardKeyRepo) LockAvailable(productID uint, quantity int) ([]model.CardKey, error) {
	var keys []model.CardKey
	// 使用事务 + 行锁
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// FOR UPDATE 行锁，防止并发
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("product_id = ? AND status = ?", productID, model.CardKeyStatusAvailable).
			Limit(quantity).
			Find(&keys).Error; err != nil {
			return err
		}

		if len(keys) < quantity {
			return gorm.ErrRecordNotFound
		}

		// 更新为预占状态
		ids := make([]uint, len(keys))
		for i, k := range keys {
			ids[i] = k.ID
		}
		return tx.Model(&model.CardKey{}).
			Where("id IN ?", ids).
			Update("status", model.CardKeyStatusLocked).Error
	})
	return keys, err
}

// MarkAsSold 标记卡密为已售出
func (r *CardKeyRepo) MarkAsSold(ids []uint, orderID uint) error {
	return r.db.Model(&model.CardKey{}).
		Where("id IN ?", ids).
		Updates(map[string]interface{}{
			"status":   model.CardKeyStatusSold,
			"order_id": orderID,
		}).Error
}

// ReleaseLockedKeys 释放锁定的卡密（订单取消）
func (r *CardKeyRepo) ReleaseLockedKeys(ids []uint) error {
	return r.db.Model(&model.CardKey{}).
		Where("id IN ? AND status = ?", ids, model.CardKeyStatusLocked).
		Updates(map[string]interface{}{
			"status":   model.CardKeyStatusAvailable,
			"order_id": nil,
		}).Error
}

// GetByOrderID 获取订单的卡密
func (r *CardKeyRepo) GetByOrderID(orderID uint) ([]model.CardKey, error) {
	var keys []model.CardKey
	err := r.db.Where("order_id = ?", orderID).Find(&keys).Error
	return keys, err
}

// CountAvailable 统计可用卡密数量
func (r *CardKeyRepo) CountAvailable(productID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.CardKey{}).
		Where("product_id = ? AND status = ?", productID, model.CardKeyStatusAvailable).
		Count(&count).Error
	return count, err
}
