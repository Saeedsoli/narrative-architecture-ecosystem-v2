// apps/backend/internal/interfaces/http/middleware/auth.go
// نسخه کامل و نهایی - بدون نیاز به تغییر

package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"narrative-architecture/apps/backend/pkg/jwt"
)

// AuthMiddleware از APIهای نیازمند احراز هویت محافظت می‌کند.
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Bearer token format required"})
			return
		}

		claims, err := jwt.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// افزودن اطلاعات کاربر به context برای استفاده در Handlerهای بعدی
		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)
		c.Set("userRoles", claims.Roles)

		c.Next()
	}
}