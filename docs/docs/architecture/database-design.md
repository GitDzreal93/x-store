---
sidebar_position: 4
title: 数据库设计
---

# 数据库设计

X-Store 使用 PostgreSQL 作为主数据库，通过 GORM AutoMigrate 自动建表。

## ER 关系图

```
┌──────────┐     ┌───────────┐     ┌──────────┐
│  users   │     │  orders   │     │ card_keys│
│──────────│     │───────────│     │──────────│
│ id       │◄──┐ │ id        │  ┌──►│ id       │
│ username │   │ │ user_id ──┼──┘  │ product_id│
│ email    │   │ │ order_no  │     │ order_id │
│ password │   └─┤ product_id│     │ content  │
│ role     │     │ status    │     │ status   │
│ oauth_*  │     │ amount    │     └──────────┘
└──────────┘     └─────┬─────┘
                       │
                 ┌─────▼─────┐
                 │ payments  │
                 │───────────│
                 │ id        │
                 │ order_id  │
                 │ channel_id│
                 │ amount    │
                 │ status    │
                 └───────────┘

┌────────────┐   ┌──────────────────┐   ┌────────────────┐
│ categories │   │ payment_channels │   │ oauth_providers│
│────────────│   │──────────────────│   │────────────────│
│ id         │   │ id               │   │ id             │
│ name       │   │ name             │   │ provider       │
│ icon       │   │ provider_type    │   │ name           │
│ sort       │   │ channel_type     │   │ enabled        │
└────────────┘   │ config_json      │   │ client_id      │
                 └──────────────────┘   │ client_secret  │
┌────────────┐                          └────────────────┘
│ products   │
│────────────│
│ id         │
│ category_id│
│ title      │
│ price      │
│ stock      │
│ status     │
└────────────┘
```

## 核心表结构

### users（用户表）

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | BIGINT PK | 主键 |
| `username` | VARCHAR(64) | 用户名，唯一索引 |
| `email` | VARCHAR(128) | 邮箱，唯一索引 |
| `password` | VARCHAR(256) | bcrypt 加密密码 |
| `nickname` | VARCHAR(64) | 昵称 |
| `avatar` | VARCHAR(512) | 头像 URL |
| `role` | VARCHAR(16) | 角色：`buyer` / `admin` |
| `status` | INT | 状态：1=正常 0=禁用 |
| `oauth_provider` | VARCHAR(32) | OAuth 提供商：github / google |
| `oauth_id` | VARCHAR(128) | OAuth 平台用户 ID |

### products（商品表）

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | BIGINT PK | 主键 |
| `category_id` | BIGINT FK | 分类 ID |
| `title` | VARCHAR(128) | 商品名称 |
| `cover` | VARCHAR(512) | 封面图 |
| `price` | DECIMAL(10,2) | 售价 |
| `original_price` | DECIMAL(10,2) | 原价 |
| `stock` | INT | 库存 |
| `sales` | INT | 销量 |
| `delivery_type` | VARCHAR(16) | 发货方式：`auto` / `manual` |
| `tags` | VARCHAR(256) | 标签，逗号分隔 |
| `status` | INT | 状态：1=上架 0=下架 |

### orders（订单表）

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | BIGINT PK | 主键 |
| `order_no` | VARCHAR(64) | 订单号，唯一索引 |
| `user_id` | BIGINT FK | 买家用户 ID |
| `product_id` | BIGINT FK | 商品 ID |
| `quantity` | INT | 数量 |
| `amount` | DECIMAL(10,2) | 总金额 |
| `status` | INT | 状态：0=待支付 1=已支付 2=已发货 3=已完成 |
| `expire_at` | TIMESTAMP | 订单过期时间 |

### card_keys（卡密表）

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | BIGINT PK | 主键 |
| `product_id` | BIGINT FK | 所属商品 |
| `order_id` | BIGINT | 关联订单（卖出后填充） |
| `content` | TEXT | 卡密内容 |
| `status` | INT | 状态：0=未售 1=已售 |

### payment_channels（支付渠道表）

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | BIGINT PK | 主键 |
| `name` | VARCHAR(64) | 渠道名称 |
| `provider_type` | VARCHAR(32) | 提供商类型 |
| `channel_type` | VARCHAR(32) | 渠道类型 |
| `config_json` | TEXT | JSON 格式配置 |
| `enabled` | BOOL | 是否启用 |

### oauth_providers（OAuth 配置表）

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | BIGINT PK | 主键 |
| `provider` | VARCHAR(32) | 提供商标识：github / google |
| `name` | VARCHAR(64) | 显示名称 |
| `enabled` | BOOL | 是否启用 |
| `client_id` | VARCHAR(256) | OAuth Client ID |
| `client_secret` | VARCHAR(512) | OAuth Client Secret |
| `redirect_url` | VARCHAR(512) | 回调地址 |

## 索引设计

- `users`: `username` UNIQUE, `email` UNIQUE, `oauth_provider + oauth_id` INDEX
- `orders`: `order_no` UNIQUE, `user_id` INDEX, `status` INDEX
- `products`: `category_id` INDEX, `status` INDEX
- `card_keys`: `product_id + status` 复合索引
- `payment_channels`: `provider_type` INDEX
- `oauth_providers`: `provider` UNIQUE
