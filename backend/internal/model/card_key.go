package model

// CardKey 卡密表
type CardKey struct {
	BaseModel
	ProductID uint   `json:"product_id" gorm:"not null;index"`
	OrderID   *uint  `json:"order_id" gorm:"index"`                                  // 关联订单（已售出时）
	Content   string `json:"content" gorm:"type:varchar(512);not null"`               // 卡密内容
	Status    int    `json:"status" gorm:"default:0;not null;index"`                  // 0=有效 1=预占(已锁定) 2=已售出
}

func (CardKey) TableName() string { return "card_keys" }

const (
	CardKeyStatusAvailable = 0
	CardKeyStatusLocked    = 1
	CardKeyStatusSold      = 2
)
