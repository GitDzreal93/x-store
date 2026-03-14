package model

// Product 商品主表（轻量，加速列表查询）
type Product struct {
	BaseModel
	CategoryID   uint    `json:"category_id" gorm:"not null;index"`
	Title        string  `json:"title" gorm:"type:varchar(256);not null"`
	Cover        string  `json:"cover" gorm:"type:varchar(512)"`                        // 封面图 URL
	Price        float64 `json:"price" gorm:"type:decimal(10,2);not null"`              // 售价
	OriginalPrice float64 `json:"original_price" gorm:"type:decimal(10,2);default:0"`   // 原价（用于划线价）
	Stock        int     `json:"stock" gorm:"default:0;not null"`                       // 数据库库存（Redis 中也维护一份）
	Sales        int     `json:"sales" gorm:"default:0;not null"`                       // 累计销量
	DeliveryType string  `json:"delivery_type" gorm:"type:varchar(16);default:auto"`    // auto=自动发货 | manual=人工代充
	Tags         string  `json:"tags" gorm:"type:varchar(512)"`                         // JSON 数组字符串 ["自动发货","独享资源"]
	Sort         int     `json:"sort" gorm:"default:0;not null"`
	Status       int     `json:"status" gorm:"default:1;not null;index"`                // 1=上架 0=下架
	IsNew        bool    `json:"is_new" gorm:"default:false"`
	IsHot        bool    `json:"is_hot" gorm:"default:false"`

	// 关联
	Category *Category      `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Detail   *ProductDetail `json:"detail,omitempty" gorm:"foreignKey:ProductID"`
}

func (Product) TableName() string { return "products" }
