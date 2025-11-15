// apps/backend/internal/interfaces/http/middleware/error_handler.go

package middleware

import (
	"errors"
	"net/http"
	"github.com/gin-gonic/gin"
)

// فرض می‌کنیم این خطاها در یک پکیج domain/errors تعریف شده‌اند.
var (
	ErrNotFound         = errors.New("resource not found")
	ErrPermissionDenied = errors.New("permission denied")
	ErrInvalidInput     = errors.New("invalid input")
	ErrConflict         = errors.New("resource conflict") // e.g., email already exists
)

// HandleError یک تابع کمکی برای مپ کردن خطاهای دامنه به کدهای وضعیت HTTP است.
func HandleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, ErrPermissionDenied):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, ErrInvalidInput):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, ErrConflict):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		// برای خطاهای ناشناخته، جزئیات را لاگ گرفته و یک پیام عمومی برمی‌گردانیم.
		c.Error(err) // این خطا توسط middleware لاگ‌گیری قابل پردازش است.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected internal error occurred"})
	}
}