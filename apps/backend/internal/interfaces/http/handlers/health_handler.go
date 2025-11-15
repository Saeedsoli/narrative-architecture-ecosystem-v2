// apps/backend/internal/interfaces/http/handlers/health_handler.go

package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// HealthHandler کنترلر HTTP برای بررسی سلامت سرویس است.
type HealthHandler struct {
	// می‌توان وابستگی‌هایی مانند کلاینت دیتابیس را برای بررسی‌های عمیق‌تر اضافه کرد.
}

// NewHealthHandler یک نمونه جدید از HealthHandler ایجاد می‌کند.
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// CheckHealth وضعیت سلامت سرویس را برمی‌گرداند.
func (h *HealthHandler) CheckHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}