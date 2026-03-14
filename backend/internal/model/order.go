package model

import "time"

// Order 订单表（与 payments 分离，一个订单可对应多次支付尝试）
type Order struct {
	BaseModel
	OrderNo    string    `json:"order_no" gorm:"type:varchar(64);uniqueIndex;not null"` // 订单编号 XS20260314001
	UserID     *uint     `json:"user_id" gorm:"index"`                                  // 可为空（游客下单）
	ProductID  uint      `json:"product_id" gorm:"not null;index"`
	Email      string    `json:"email" gorm:"type:varchar(128);not null"`                // 接收卡密的邮箱
	Amount     float64   `json:"amount" gorm:"type:decimal(10,2);not null"`
	Status     int       `json:"status" gorm:"default:0;not null;index"`                 // 0=待支付 1=已支付 2=已发货 3=已完成 4=已取消 5=已退款
	PayMethod  string    `json:"pay_method" gorm:"type:varchar(32)"`                     // wechat | alipay | stripe
	PaidAt     *time.Time `json:"paid_at"`
	ExpireAt   time.Time  `json:"expire_at" gorm:"not null"`                             // 库存锁定过期时间

	// 关联
	Product  *Product  `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	CardKeys []CardKey `json:"card_keys,omitempty" gorm:"foreignKey:OrderID"`
}

func (Order) TableName() string { return "orders" }

// 订单状态常量
const (
	OrderStatusPending   = 0
	OrderStatusPaid      = 1
	OrderStatusDelivered = 2
	OrderStatusCompleted = 3
	OrderStatusCancelled = 4
	OrderStatusRefunded  = 5
)
