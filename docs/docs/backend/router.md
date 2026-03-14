---
sidebar_position: 4
title: 路由设计
---

# 路由设计

所有路由在 `backend/internal/router/router.go` 中集中注册，按权限分组。

## 路由分组

```go title="backend/internal/router/router.go"
api := r.Group("/api")

// 公开接口（无需认证）
api.POST("/users/register", userH.Register)
api.POST("/users/login", userH.Login)
api.GET("/oauth/providers", oauthH.ListProviders)
api.GET("/categories", catH.List)
api.GET("/products", prodH.List)
// ...

// 需要 JWT 认证的接口
authed := api.Group("").Use(middleware.JWTAuth())
authed.GET("/user/profile", userH.Profile)
// ...

// 管理员接口（JWT + admin 角色）
admin := api.Group("/admin").Use(middleware.JWTAuth(), middleware.AdminOnly())
admin.GET("/products", prodH.AdminList)
admin.POST("/products", prodH.Create)
// ...

// 防重放签名接口（核心写操作）
signed := api.Group("").Use(middleware.AntiReplay(), middleware.RateLimit())
signed.POST("/orders", orderH.Create)
signed.POST("/orders/:order_no/pay", orderH.Pay)
```

## 完整路由表

### 公开接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/users/register` | 用户注册 |
| POST | `/api/users/login` | 用户登录 |
| GET | `/api/oauth/providers` | 获取已启用的 OAuth 提供商 |
| GET | `/api/oauth/github` | GitHub OAuth 跳转 |
| GET | `/api/oauth/github/callback` | GitHub OAuth 回调 |
| GET | `/api/oauth/google` | Google OAuth 跳转 |
| GET | `/api/oauth/google/callback` | Google OAuth 回调 |
| GET | `/api/categories` | 分类列表 |
| GET | `/api/products` | 商品列表 |
| GET | `/api/products/:id` | 商品详情 |
| GET | `/api/orders/:order_no` | 订单查询 |
| GET | `/api/payment/channels` | 支付渠道列表 |
| POST | `/api/webhooks/payment/:channel_id` | 支付回调 |

### 认证接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/user/profile` | 当前用户信息 |
| GET | `/api/user/orders` | 我的订单 |

### 管理员接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/admin/dashboard/stats` | 仪表盘统计 |
| CRUD | `/api/admin/products` | 商品管理 |
| CRUD | `/api/admin/categories` | 分类管理 |
| GET | `/api/admin/orders` | 订单列表 |
| POST | `/api/admin/orders/:order_no/refund` | 订单退款 |
| POST | `/api/admin/cardkeys/import` | 卡密导入 |
| CRUD | `/api/admin/payment-channels` | 支付渠道管理 |
| GET | `/api/admin/oauth-providers` | OAuth 配置列表 |
| PUT | `/api/admin/oauth-providers/:id` | 更新 OAuth 配置 |
| POST | `/api/admin/oauth-providers/:id/toggle` | 切换 OAuth 启用状态 |

### 防重放接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/orders` | 创建订单 |
| POST | `/api/orders/:order_no/pay` | 订单支付 |

## 路由设计原则

1. **RESTful 风格**：资源用名词，操作用 HTTP 方法
2. **权限分层**：公开 → 认证 → 管理员，层层递进
3. **核心写操作**：加防重放签名，防止刷单
4. **Webhook 独立**：支付回调使用独立路径，不经过认证中间件
