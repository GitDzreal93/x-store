---
sidebar_position: 2
title: 用户认证
---

# 用户认证 API

## 用户注册

```
POST /api/users/register
```

**请求体**：

```json
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "123456",
  "nickname": "测试用户"
}
```

**响应**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "nickname": "测试用户",
      "role": "buyer",
      "status": 1
    }
  }
}
```

**字段校验**：
- `username`: 3-32 字符，唯一
- `email`: 合法邮箱格式，唯一
- `password`: 至少 6 位

## 用户登录

```
POST /api/users/login
```

**请求体**：

```json
{
  "username": "testuser",
  "password": "123456"
}
```

**响应**：同注册接口。

## 获取当前用户信息

```
GET /api/user/profile
Authorization: Bearer <token>
```

**响应**：

```json
{
  "code": 0,
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "nickname": "测试用户",
    "avatar": "",
    "role": "buyer",
    "status": 1
  }
}
```

## OAuth 登录

### 获取已启用的 OAuth 提供商

```
GET /api/oauth/providers
```

**响应**：

```json
{
  "code": 0,
  "data": [
    {
      "name": "github",
      "label": "GitHub",
      "auth_url": "https://github.com/login/oauth/authorize?client_id=..."
    },
    {
      "name": "google",
      "label": "Google",
      "auth_url": "https://accounts.google.com/o/oauth2/v2/auth?..."
    }
  ]
}
```

如果没有启用任何 OAuth，`data` 为空数组。

### GitHub 登录

```
GET /api/oauth/github
```

重定向到 GitHub 授权页。授权成功后回调 `/api/oauth/github/callback`，最终重定向到前端 `/auth/callback?token=xxx`。

### Google 登录

```
GET /api/oauth/google
```

流程同 GitHub。
