package model

// ProductDetail 商品详情表（长文本抽离，避免拖慢主表查询）
type ProductDetail struct {
	BaseModel
	ProductID   uint   `json:"product_id" gorm:"uniqueIndex;not null"`
	Description string `json:"description" gorm:"type:text"`       // Markdown 富文本商品说明
	Notice      string `json:"notice" gorm:"type:text"`            // 购买须知
	Extra       string `json:"extra" gorm:"type:text"`             // 额外 JSON 扩展字段
}

func (ProductDetail) TableName() string { return "product_details" }
