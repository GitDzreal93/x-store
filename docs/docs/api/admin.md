---
sidebar_position: 6
title: 管理后台
---

# 管理后台 API

所有管理后台接口需要 JWT 认证且用户角色为 `admin`。

```
Authorization: Bearer &lt;admin_token&gt;
```

## 仪表盘统计

```
GET /api/admin/dashboard/stats
```

**响应**：

```json
{
  "code": 0,
  "data": {
    "today_gmv": 1299.50,
    "today_orders": 42,
    "total_users": 156,
    "total_products": 25,
    "order_trend": [
      { "date": "2026-03-08", "count": 15, "amount": 450.00 },
      { "date": "2026-03-09", "count": 23, "amount": 689.70 }
    ],
    "top_products": [
      { "id": 1, "title": "ChatGPT Plus", "sales": 256, "revenue": 7654.40 }
    ]
  }
}
```

## 商品管理

### 商品列表（管理）

```
GET /api/admin/products?page=1&page_size=10
```

### 创建商品

```
POST /api/admin/products
```

```json
{
  "category_id": 1,
  "title": "ChatGPT Plus 一个月",
  "cover": "https://example.com/cover.jpg",
  "price": 29.90,
  "original_price": 49.90,
  "delivery_type": "auto",
  "tags": "热销,AI",
  "status": 1
}
```

### 更新商品

```
PUT /api/admin/products/:id
```

请求体同创建，字段可选更新。

### 删除商品

```
DELETE /api/admin/products/:id
```

## 分类管理

### 创建分类

```
POST /api/admin/categories
```

```json
{
  "name": "AI 大模型",
  "icon": "🤖",
  "sort": 100,
  "status": 1
}
```

### 更新分类

```
PUT /api/admin/categories/:id
```

### 删除分类

```
DELETE /api/admin/categories/:id
```

## 订单管理

### 订单列表

```
GET /api/admin/orders?page=1&page_size=10&status=0
```

### 订单退款

```
POST /api/admin/orders/:order_no/refund
```

### 手动发货

```
POST /api/admin/orders/:order_no/deliver
```

## 卡密管理

### 批量导入卡密

```
POST /api/admin/cardkeys/import
```

```json
{
  "product_id": 1,
  "content": "KEY-001-ABCD\nKEY-002-EFGH\nKEY-003-IJKL"
}
```

每行一个卡密，自动按换行符分割。

### 查询可用卡密数量

```
GET /api/admin/cardkeys/count/:product_id
```

```json
{
  "code": 0,
  "data": {
    "product_id": 1,
    "available": 45,
    "sold": 256
  }
}
```

## 支付渠道管理

### 渠道列表

```
GET /api/admin/payment-channels
```

### 更新渠道

```
PUT /api/admin/payment-channels/:id
```

### 切换渠道启用状态

```
POST /api/admin/payment-channels/:id/toggle
```

## OAuth 提供商管理

### 获取所有配置

```
GET /api/admin/oauth-providers
```

响应中 `client_secret` 显示为 `******`（脱敏）。

### 获取完整配置

```
GET /api/admin/oauth-providers/:id
```

返回包含完整 `client_secret` 的配置。

### 更新配置

```
PUT /api/admin/oauth-providers/:id
```

```json
{
  "enabled": true,
  "client_id": "your_github_client_id",
  "client_secret": "your_github_client_secret",
  "redirect_url": "http://localhost:8082/api/oauth/github/callback"
}
```

### 切换启用状态

```
POST /api/admin/oauth-providers/:id/toggle
```

一键切换 `enabled` 状态。
