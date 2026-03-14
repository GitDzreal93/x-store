---
sidebar_position: 3
title: 数据模型
---

# 数据模型

所有数据模型位于 `backend/internal/model/` 目录，使用 GORM 标签定义表结构。

## BaseModel

所有模型继承 `BaseModel`，包含通用字段：

```go title="backend/internal/model/base.go"
type BaseModel struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

## User（用户）

```go title="backend/internal/model/user.go"
type User struct {
    BaseModel
    Username      string `json:"username" gorm:"type:varchar(64);uniqueIndex;not null"`
    Email         string `json:"email" gorm:"type:varchar(128);uniqueIndex;not null"`
    Password      string `json:"-" gorm:"type:varchar(256)"`
    Nickname      string `json:"nickname" gorm:"type:varchar(64)"`
    Avatar        string `json:"avatar" gorm:"type:varchar(512)"`
    Role          string `json:"role" gorm:"type:varchar(16);default:buyer;not null;index"`
    Status        int    `json:"status" gorm:"default:1;not null"`
    OAuthProvider string `json:"oauth_provider,omitempty" gorm:"type:varchar(32);index"`
    OAuthID       string `json:"oauth_id,omitempty" gorm:"type:varchar(128);index"`
}
```

**要点**：
- `Password` 使用 `json:"-"` 避免序列化到响应
- `Role` 区分 `buyer`（买家）和 `admin`（管理员）
- OAuth 字段支持第三方登录绑定

## Product（商品）

```go title="backend/internal/model/product.go"
type Product struct {
    BaseModel
    CategoryID    uint    `json:"category_id" gorm:"index;not null"`
    Title         string  `json:"title" gorm:"type:varchar(128);not null"`
    Cover         string  `json:"cover" gorm:"type:varchar(512)"`
    Price         float64 `json:"price" gorm:"type:decimal(10,2);not null"`
    OriginalPrice float64 `json:"original_price" gorm:"type:decimal(10,2)"`
    Stock         int     `json:"stock" gorm:"default:0"`
    Sales         int     `json:"sales" gorm:"default:0"`
    DeliveryType  string  `json:"delivery_type" gorm:"type:varchar(16);default:auto"`
    Tags          string  `json:"tags" gorm:"type:varchar(256)"`
    Sort          int     `json:"sort" gorm:"default:0"`
    Status        int     `json:"status" gorm:"default:1;not null;index"`
    IsNew         bool    `json:"is_new" gorm:"default:false"`
    IsHot         bool    `json:"is_hot" gorm:"default:false"`
}
```

## Order（订单）

```go title="backend/internal/model/order.go"
type Order struct {
    BaseModel
    OrderNo    string    `json:"order_no" gorm:"type:varchar(64);uniqueIndex;not null"`
    UserID     uint      `json:"user_id" gorm:"index;not null"`
    ProductID  uint      `json:"product_id" gorm:"not null"`
    Quantity   int       `json:"quantity" gorm:"default:1;not null"`
    Amount     float64   `json:"amount" gorm:"type:decimal(10,2);not null"`
    Status     int       `json:"status" gorm:"default:0;not null;index"`
    ExpireAt   time.Time `json:"expire_at"`
}
```

**订单状态流转**：

```
0 (待支付) → 1 (已支付) → 2 (已发货) → 3 (已完成)
     ↓
 -1 (已取消/超时)
```

## CardKey（卡密）

```go
type CardKey struct {
    BaseModel
    ProductID uint   `json:"product_id" gorm:"index;not null"`
    OrderID   uint   `json:"order_id" gorm:"index"`
    Content   string `json:"content" gorm:"type:text;not null"`
    Status    int    `json:"status" gorm:"default:0;not null;index"`
}
```

状态值：`0` = 未售出，`1` = 已售出

## PaymentChannel（支付渠道）

```go
type PaymentChannel struct {
    BaseModel
    Name            string  `json:"name" gorm:"type:varchar(64);not null"`
    ProviderType    string  `json:"provider_type" gorm:"type:varchar(32);not null;index"`
    ChannelType     string  `json:"channel_type" gorm:"type:varchar(32);not null"`
    InteractionMode string  `json:"interaction_mode" gorm:"type:varchar(16)"`
    ConfigJSON      string  `json:"config_json" gorm:"type:text"`
    FeeRate         float64 `json:"fee_rate" gorm:"type:decimal(5,2)"`
    MinAmount       float64 `json:"min_amount" gorm:"type:decimal(10,2)"`
    Enabled         bool    `json:"enabled" gorm:"default:false"`
    Sort            int     `json:"sort" gorm:"default:0"`
}
```

## OAuthProvider（OAuth 配置）

```go
type OAuthProvider struct {
    BaseModel
    Provider     string `json:"provider" gorm:"type:varchar(32);uniqueIndex;not null"`
    Name         string `json:"name" gorm:"type:varchar(64);not null"`
    Enabled      bool   `json:"enabled" gorm:"default:false;not null;index"`
    ClientID     string `json:"client_id" gorm:"type:varchar(256);not null"`
    ClientSecret string `json:"client_secret,omitempty" gorm:"type:varchar(512);not null"`
    RedirectURL  string `json:"redirect_url" gorm:"type:varchar(512);not null"`
    Sort         int    `json:"sort" gorm:"default:0;not null"`
}
```
