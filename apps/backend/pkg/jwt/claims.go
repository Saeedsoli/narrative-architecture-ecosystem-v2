// apps/backend/pkg/jwt/claims.go

package jwt

import "github.com/golang-jwt/jwt/v5"

// Claims اطلاعاتی است که در Access Token JWT ذخیره می‌شود.
type Claims struct {
	UserID string   `json:"user_id"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}