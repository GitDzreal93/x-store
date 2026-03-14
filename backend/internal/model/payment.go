package model

import "time"

// Payment 支付记录表（一个订单可对应多次支付尝试）
type Payment struct {
	BaseModel
	OrderID     uint      `json:"order_id" gorm:"not null;index"`
	TradeNo     string    `json:"trade_no" gorm:"type:varchar(128);index"`  // 第三方支付流水号
	PayMethod   string    `json:"pay_method" gorm:"type:varchar(32);not null"`
	Amount      float64   `json:"amount" gorm:"type:decimal(10,2);not null"`
	Status      int       `json:"status" gorm:"default:0;not null"`          // 0=处理中 1=成功 2=失败
	RawNotify   string    `json:"raw_notify" gorm:"type:text"`              // 网关原始回调数据
	CompletedAt *time.Time `json:"completed_at"`
}

func (Payment) TableName() string { return "payments" }

const (
	PaymentStatusProcessing = 0
	PaymentStatusSuccess    = 1
	PaymentStatusFailed     = 2
)
