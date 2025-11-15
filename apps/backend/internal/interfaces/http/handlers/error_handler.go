package handlers

import (
	"errors"
	"net/http"
	"github.com/gin-gonic/gin"
)

// فرض می‌کنیم خطاهای Domain تعریف شده‌اند
var (
	ErrNotFound      = errors.New("resource not found")
	ErrPermissionDenied = errors.New("permission denied")
	ErrInvalidInput  = errors.New("invalid input")
)

func HandleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, ErrPermissionDenied):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, ErrInvalidInput):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
	}
}