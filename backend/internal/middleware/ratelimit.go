package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/x-store/backend/internal/config"
	"github.com/x-store/backend/pkg/response"
)

// RateLimit 基于 Redis 的 IP 限流中间件
// maxRequests: 时间窗口内最大请求数, window: 时间窗口
func RateLimit(maxRequests int64, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf("ratelimit:%s:%s", c.FullPath(), ip)
		ctx := context.Background()

		current, err := config.RDB.Incr(ctx, key).Result()
		if err != nil {
			response.ServerError(c, "限流服务异常")
			c.Abort()
			return
		}

		// 首次请求时设置过期时间
		if current == 1 {
			config.RDB.Expire(ctx, key, window)
		}

		if current > maxRequests {
			response.Fail(c, 429, 429, "系统繁忙，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}
