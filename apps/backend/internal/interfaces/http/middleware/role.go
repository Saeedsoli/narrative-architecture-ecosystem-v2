// apps/backend/internal/interfaces/http/middleware/role.go

package middleware

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// RoleMiddleware یک Middleware برای بررسی اینکه آیا کاربر حداقل یکی از نقش‌های مجاز را دارد، ایجاد می‌کند.
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// این مقدار باید توسط AuthMiddleware در context قرار داده شده باشد.
		userRoles, exists := c.Get("userRoles")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied: user roles not found in context"})
			return
		}

		roles, ok := userRoles.([]string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error: user roles are not in the expected format"})
			return
		}

		// ساخت یک map برای جستجوی سریع نقش‌های مجاز
		allowedMap := make(map[string]struct{}, len(allowedRoles))
		for _, role := range allowedRoles {
			allowedMap[role] = struct{}{}
		}

		// بررسی اینکه آیا کاربر حداقل یکی از نقش‌های مجاز را دارد
		isAllowed := false
		for _, userRole := range roles {
			if _, ok := allowedMap[userRole]; ok {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied: insufficient permissions"})
			return
		}

		c.Next()
	}
}