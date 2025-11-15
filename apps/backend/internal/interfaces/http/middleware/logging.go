// apps/backend/internal/interfaces/http/middleware/logging.go

package middleware

import (
	"time"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// LoggingMiddleware یک Middleware برای لاگ کردن ساختاریافته درخواست‌ها است.
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		c.Next() // پردازش درخواست

		latency := time.Since(start)
		status := c.Writer.Status()
		
		logEvent := log.Info()
		if status >= 500 {
			logEvent = log.Error().Err(c.Errors.Last())
		} else if status >= 400 {
			logEvent = log.Warn()
		}

		logEvent.
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", status).
			Str("latency", latency.String()).
			Str("client_ip", c.ClientIP()).
			Msg("request processed")
	}
}