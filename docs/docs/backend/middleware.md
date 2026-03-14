---
sidebar_position: 5
title: 中间件
---

# 中间件

X-Store 后端使用多个自定义中间件保障安全和性能。

## JWT 认证

```go title="backend/internal/middleware/auth.go"
func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 从 Authorization Header 获取 Token
        token := c.GetHeader("Authorization")
        token = strings.TrimPrefix(token, "Bearer ")

        // 2. 解析 Token
        claims, err := crypto.ParseToken(config.Global.JWT.Secret, token)
        if err != nil {
            response.Unauthorized(c, "认证失败")
            c.Abort()
            return
        }

        // 3. 将用户信息存入 Context
        c.Set("userID", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("role", claims.Role)
        c.Next()
    }
}
```

**使用方式**：在 Handler 中通过 `c.GetUint("userID")` 获取当前用户 ID。

## 管理员权限

```go title="backend/internal/middleware/admin.go"
func AdminOnly() gin.HandlerFunc {
    return func(c *gin.Context) {
        role := c.GetString("role")
        if role != "admin" {
            response.Forbidden(c, "需要管理员权限")
            c.Abort()
            return
        }
        c.Next()
    }
}
```

必须搭配 `JWTAuth()` 使用，确保 `role` 已被设置。

## 防重放攻击

```go title="backend/internal/middleware/anti_replay.go"
func AntiReplay() gin.HandlerFunc {
    return func(c *gin.Context) {
        nonce := c.GetHeader("X-Nonce")
        timestamp := c.GetHeader("X-Timestamp")
        sign := c.GetHeader("X-Signature")

        // 1. 验证时间戳有效期
        // 2. 验证签名 = MD5(nonce + timestamp + secret)
        // 3. 检查 nonce 是否已使用（Redis SETNX）
    }
}
```

**原理**：
1. 客户端每次请求生成唯一 `nonce`
2. 加上当前时间戳和密钥计算签名
3. 服务端验签并检查 nonce 是否已消费
4. 已消费的 nonce 存入 Redis，设置过期时间

## 接口限流

```go title="backend/internal/middleware/rate_limit.go"
func RateLimit(maxRequests int, window time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        key := "rate:" + c.ClientIP() + ":" + c.FullPath()
        // 基于 Redis 的滑动窗口限流
        count := redis.Incr(key)
        if count == 1 {
            redis.Expire(key, window)
        }
        if count > maxRequests {
            response.TooManyRequests(c, "请求过于频繁")
            c.Abort()
            return
        }
        c.Next()
    }
}
```

## CORS 跨域

```go title="backend/internal/middleware/cors.go"
func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
        c.Header("Access-Control-Allow-Headers",
            "Authorization,Content-Type,X-Nonce,X-Timestamp,X-Signature")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    }
}
```

## 中间件执行顺序

```
请求进入
  → CORS（跨域处理）
    → JWTAuth（认证，可选）
      → AdminOnly（权限，可选）
        → AntiReplay（防重放，可选）
          → RateLimit（限流，可选）
            → Handler（业务处理）
```
