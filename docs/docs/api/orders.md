---
sidebar_position: 4
title: 订单
---

# 订单 API

## 创建订单

```
POST /api/orders
Authorization: Bearer <token>
X-Nonce: <uuid>
X-Timestamp: <timestamp>
X-Signature: <signature>
```

:::caution
该接口需要防重放签名，请参考 [API 概览](/docs/api/overview) 中的签名说明。
:::

**请求体**：

```json
{
  "product_id": 1,
  "quantity": 1
}
```

**响应**：

```json
{
  "code": 0,
  "data": {
    "order_no": "ORD20260314175300001",
    "product_id": 1,
    "quantity": 1,
    "amount": 29.90,
    "status": 0,
    "expire_at": "2026-03-14T18:23:00Z"
  }
}
```

**订单状态值**：

| 值 | 说明 |
|----|------|
| `0` | 待支付 |
| `1` | 已支付 |
| `2` | 已发货 |
| `3` | 已完成 |
| `-1` | 已取消 / 已退款 |

## 查询订单

```
GET /api/orders/:order_no
```

**响应**：

```json
{
  "code": 0,
  "data": {
    "order_no": "ORD20260314175300001",
    "user_id": 1,
    "product_id": 1,
    "quantity": 1,
    "amount": 29.90,
    "status": 2,
    "card_keys": [
      {
        "content": "CHATGPT-PLUS-DEMO-001-ABCD1234"
      }
    ]
  }
}
```

订单状态为已发货（2）或已完成（3）时，会返回 `card_keys` 字段。

## 我的订单列表

```
GET /api/user/orders
Authorization: Bearer <token>
```

**查询参数**：

| 参数 | 类型 | 说明 |
|------|------|------|
| `status` | int | 按状态筛选 |
| `page` | int | 页码 |
| `page_size` | int | 每页数量 |

**响应**：

```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "order_no": "ORD20260314175300001",
        "product_id": 1,
        "product_title": "ChatGPT Plus 一个月",
        "quantity": 1,
        "amount": 29.90,
        "status": 2,
        "created_at": "2026-03-14T17:53:00Z"
      }
    ],
    "total": 5,
    "page": 1,
    "page_size": 10
  }
}
```
