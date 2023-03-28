package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"net/http"
	"time"
)

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(fillInterval time.Duration, cap, quantum int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(fillInterval, cap, quantum)
	return func(c *gin.Context) {
		// 如果取不到令牌就中断本次请求返回
		if bucket.TakeAvailable(1) < 1 {
			c.String(http.StatusForbidden, "接口限速中，请稍后尝试......")
			c.Abort()
			return
		}
		c.Next()
	}
}
