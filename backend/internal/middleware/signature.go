package middleware

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/x-store/backend/internal/config"
	"github.com/x-store/backend/pkg/crypto"
	"github.com/x-store/backend/pkg/response"
)

// AntiReplay 防重放签名验证中间件
// 前端请求需携带 Header: X-Timestamp, X-Nonce, X-Signature
func AntiReplay() gin.HandlerFunc {
	return func(c *gin.Context) {
		timestamp := c.GetHeader("X-Timestamp")
		nonce := c.GetHeader("X-Nonce")
		signature := c.GetHeader("X-Signature")

		if timestamp == "" || nonce == "" || signature == "" {
			response.BadRequest(c, "缺少签名参数")
			c.Abort()
			return
		}

		// 验证时间戳是否在有效窗口内
		ts, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			response.BadRequest(c, "时间戳格式错误")
			c.Abort()
			return
		}
		diff := math.Abs(float64(time.Now().Unix() - ts))
		if diff > float64(config.Global.Signature.ExpireSeconds) {
			response.BadRequest(c, "请求已过期")
			c.Abort()
			return
		}

		// 验证 nonce 唯一性（Redis 去重，防止重放）
		nonceKey := fmt.Sprintf("nonce:%s", nonce)
		ctx := context.Background()
		exists, err := config.RDB.Exists(ctx, nonceKey).Result()
		if err != nil {
			response.ServerError(c, "签名验证服务异常")
			c.Abort()
			return
		}
		if exists > 0 {
			response.BadRequest(c, "重复请求")
			c.Abort()
			return
		}

		// 验证签名
		params := make(map[string]string)
		// 将 query 参数纳入签名
		for k, v := range c.Request.URL.Query() {
			if len(v) > 0 {
				params[k] = v[0]
			}
		}
		if !crypto.VerifySignature(params, timestamp, nonce, config.Global.Signature.Secret, signature) {
			response.BadRequest(c, "签名验证失败")
			c.Abort()
			return
		}

		// 签名验证通过，将 nonce 写入 Redis 并设置过期时间
		config.RDB.Set(ctx, nonceKey, "1", time.Duration(config.Global.Signature.ExpireSeconds*2)*time.Second)

		c.Next()
	}
}
