---
sidebar_position: 1
title: API 概览
---

# API 概览

X-Store 后端提供 RESTful API，所有接口以 `/api` 为前缀。

## 基础信息

| 项目 | 值 |
|------|-----|
| **Base URL** | `http://localhost:8082/api` |
| **协议** | HTTP / HTTPS |
| **数据格式** | JSON |
| **认证方式** | Bearer Token (JWT) |

## 统一响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

### 状态码说明

| code | HTTP Status | 说明 |
|------|------------|------|
| `0` | 200 | 成功 |
| `400` | 400 | 参数错误 |
| `401` | 401 | 未认证 / Token 无效 |
| `403` | 403 | 无权限 |
| `404` | 404 | 资源不存在 |
| `429` | 429 | 请求过于频繁 |
| `500` | 500 | 服务器内部错误 |

## 认证方式

需要认证的接口，在 HTTP Header 中携带 JWT Token：

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

Token 通过登录或注册接口获取，有效期由 `config.yaml` 的 `jwt.expire_hours` 控制。

## 防重放签名

创建订单和支付等核心写接口需要防重放签名，请求头需携带：

| Header | 说明 |
|--------|------|
| `X-Nonce` | 随机字符串（UUID），每次请求唯一 |
| `X-Timestamp` | 当前时间戳（秒） |
| `X-Signature` | 签名 = MD5(nonce + timestamp + secret) |

```bash
# 示例
NONCE=$(uuidgen)
TIMESTAMP=$(date +%s)
SIGNATURE=$(echo -n "${NONCE}${TIMESTAMP}${SECRET}" | md5sum | cut -d' ' -f1)

curl -X POST http://localhost:8082/api/orders \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Nonce: $NONCE" \
  -H "X-Timestamp: $TIMESTAMP" \
  -H "X-Signature: $SIGNATURE" \
  -H "Content-Type: application/json" \
  -d '{"product_id": 1, "quantity": 1}'
```

## 接口分类

| 类别 | 前缀 | 认证 | 说明 |
|------|------|------|------|
| [用户认证](/docs/api/auth) | `/api/users/*`, `/api/oauth/*` | ❌ | 注册、登录、OAuth |
| [商品](/docs/api/products) | `/api/products/*`, `/api/categories/*` | ❌ | 商品和分类查询 |
| [订单](/docs/api/orders) | `/api/orders/*` | 部分 | 订单创建和查询 |
| [支付](/docs/api/payment) | `/api/payment/*`, `/api/webhooks/*` | 部分 | 支付和回调 |
| [管理后台](/docs/api/admin) | `/api/admin/*` | ✅ Admin | 后台管理接口 |
