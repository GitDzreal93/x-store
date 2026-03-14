---
sidebar_position: 5
title: 支付
---

# 支付 API

## 获取支付渠道列表

```
GET /api/payment/channels
```

**响应**：

```json
{
  "code": 0,
  "data": [
    {
      "id": 1,
      "name": "模拟支付",
      "provider_type": "mock",
      "channel_type": "mock",
      "interaction_mode": "redirect",
      "enabled": true
    },
    {
      "id": 2,
      "name": "支付宝-当面付",
      "provider_type": "alipay",
      "channel_type": "alipay_f2f",
      "interaction_mode": "qrcode",
      "enabled": true
    }
  ]
}
```

## 订单支付

```
POST /api/orders/:order_no/pay
Authorization: Bearer <token>
X-Nonce: <uuid>
X-Timestamp: <timestamp>
X-Signature: <signature>
```

:::caution
该接口需要防重放签名。
:::

**请求体**：

```json
{
  "channel_id": 1
}
```

**响应**：

```json
{
  "code": 0,
  "data": {
    "payment_id": 1,
    "pay_url": "https://pay.example.com/checkout?id=xxx",
    "qr_code": "",
    "trade_no": "PAY20260314001",
    "expire_at": 1710417780
  }
}
```

**字段说明**：

| 字段 | 说明 |
|------|------|
| `pay_url` | 支付链接（跳转或嵌入） |
| `qr_code` | 二维码链接（扫码支付时返回） |
| `trade_no` | 第三方交易号 |
| `expire_at` | 支付过期时间（Unix 时间戳） |

## 支付回调（Webhook）

```
POST /api/webhooks/payment/:channel_id
```

该接口由支付平台主动调用，用于通知支付结果。不同支付提供商的请求格式不同，系统会自动根据 `channel_id` 匹配对应的 Provider 进行验签。

**处理流程**：

1. 根据 `channel_id` 查询支付渠道配置
2. 创建对应的 Provider 实例
3. 调用 `VerifyNotify()` 验签并解析回调数据
4. 更新订单状态为「已支付」
5. 自动发货（分配卡密）
6. 返回成功响应给支付平台

**注意**：
- Webhook 不经过 JWT 认证中间件
- 验签逻辑在各 Provider 内部实现
- 支付成功后会自动触发卡密分配和邮件通知

## 查询支付状态

```
GET /api/payment/:payment_id/status
```

**响应**：

```json
{
  "code": 0,
  "data": {
    "payment_id": 1,
    "status": "paid",
    "amount": 29.90,
    "channel": "mock",
    "trade_no": "PAY20260314001"
  }
}
```

**支付状态**：

| 值 | 说明 |
|----|------|
| `pending` | 待支付 |
| `paid` | 已支付 |
| `failed` | 支付失败 |
| `refunded` | 已退款 |
