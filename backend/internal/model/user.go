package model

// User 用户表（买家 & 管理员共用）
type User struct {
	BaseModel
	Username      string `json:"username" gorm:"type:varchar(64);uniqueIndex;not null"`
	Email         string `json:"email" gorm:"type:varchar(128);uniqueIndex;not null"`
	Password      string `json:"-" gorm:"type:varchar(256)"`
	Nickname      string `json:"nickname" gorm:"type:varchar(64)"`
	Avatar        string `json:"avatar" gorm:"type:varchar(512)"`
	Role          string `json:"role" gorm:"type:varchar(16);default:buyer;not null;index"` // buyer | admin
	Status        int    `json:"status" gorm:"default:1;not null"`                          // 1=正常 0=禁用
	OAuthProvider string `json:"oauth_provider,omitempty" gorm:"type:varchar(32);index"`    // github | google | ""
	OAuthID       string `json:"oauth_id,omitempty" gorm:"type:varchar(128);index"`         // 第三方平台用户 ID
}

func (User) TableName() string { return "users" }
