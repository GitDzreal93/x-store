package model

// Category 商品分类
type Category struct {
	BaseModel
	Name     string `json:"name" gorm:"type:varchar(64);not null"`
	Icon     string `json:"icon" gorm:"type:varchar(64)"`           // emoji 或图标名
	Sort     int    `json:"sort" gorm:"default:0;not null"`         // 排序权重，越大越靠前
	ParentID *uint  `json:"parent_id" gorm:"index"`                 // 支持二级分类
	Status   int    `json:"status" gorm:"default:1;not null;index"` // 1=启用 0=禁用
}

func (Category) TableName() string { return "categories" }
