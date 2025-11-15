// apps/backend/internal/interfaces/http/middleware/rate_limit.go

package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// RateLimiterMiddleware یک Middleware برای محدودسازی تعداد درخواست‌ها با استفاده از Redis است.
func RateLimiterMiddleware(redisClient *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("ratelimit:%s", c.ClientIP())
		
		// استفاده از الگوریتم Fixed Window برای سادگی
		count, err := redisClient.Incr(c.Request.Context(), key).Result()
		if err != nil {
			c.Next()
			return
		}

		if count == 1 {
			redisClient.Expire(c.Request.Context(), key, window)
		}

		if count > int64(limit) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			return
		}

		c.Next()
	}
}