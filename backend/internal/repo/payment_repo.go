package repo

import (
	"github.com/x-store/backend/internal/config"
	"github.com/x-store/backend/internal/model"
	"gorm.io/gorm"
)

type PaymentRepo struct {
	db *gorm.DB
}

func NewPaymentRepo() *PaymentRepo {
	return &PaymentRepo{db: config.DB}
}

func (r *PaymentRepo) Create(p *model.Payment) error {
	return r.db.Create(p).Error
}

func (r *PaymentRepo) GetByID(id uint) (*model.Payment, error) {
	var p model.Payment
	if err := r.db.First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PaymentRepo) GetByOrderID(orderID uint) ([]model.Payment, error) {
	var list []model.Payment
	if err := r.db.Where("order_id = ?", orderID).Order("id DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *PaymentRepo) Update(p *model.Payment) error {
	return r.db.Save(p).Error
}

// GetLatestByOrderID 获取订单最新的支付记录
func (r *PaymentRepo) GetLatestByOrderID(orderID uint) (*model.Payment, error) {
	var p model.Payment
	if err := r.db.Where("order_id = ?", orderID).Order("id DESC").First(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

// --- PaymentChannel Repo ---

type PaymentChannelRepo struct {
	db *gorm.DB
}

func NewPaymentChannelRepo() *PaymentChannelRepo {
	return &PaymentChannelRepo{db: config.DB}
}

func (r *PaymentChannelRepo) Create(c *model.PaymentChannel) error {
	return r.db.Create(c).Error
}

func (r *PaymentChannelRepo) GetByID(id uint) (*model.PaymentChannel, error) {
	var c model.PaymentChannel
	if err := r.db.First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *PaymentChannelRepo) List(isActive *bool) ([]model.PaymentChannel, error) {
	var list []model.PaymentChannel
	q := r.db.Order("sort DESC, id ASC")
	if isActive != nil {
		q = q.Where("is_active = ?", *isActive)
	}
	if err := q.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *PaymentChannelRepo) Update(c *model.PaymentChannel) error {
	return r.db.Save(c).Error
}

func (r *PaymentChannelRepo) Delete(id uint) error {
	return r.db.Delete(&model.PaymentChannel{}, id).Error
}
