package model

// OAuthProvider OAuth 提供商配置表
type OAuthProvider struct {
	BaseModel
	Provider     string `json:"provider" gorm:"type:varchar(32);uniqueIndex;not null"`      // github | google
	Name         string `json:"name" gorm:"type:varchar(64);not null"`                      // GitHub | Google
	Enabled      bool   `json:"enabled" gorm:"default:false;not null;index"`                // 是否启用
	ClientID     string `json:"client_id" gorm:"type:varchar(256);not null"`                // OAuth 应用 ID
	ClientSecret string `json:"client_secret,omitempty" gorm:"type:varchar(512);not null"`  // OAuth 应用密钥（加密存储）
	RedirectURL  string `json:"redirect_url" gorm:"type:varchar(512);not null"`             // 回调地址
	Sort         int    `json:"sort" gorm:"default:0;not null"`                             // 排序权重
}

func (OAuthProvider) TableName() string { return "oauth_providers" }
