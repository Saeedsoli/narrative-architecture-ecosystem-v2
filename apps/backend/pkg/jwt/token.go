// apps/backend/pkg/jwt/token.go

package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateTokens یک جفت Access Token و Refresh Token تولید می‌کند.
func GenerateTokens(userID, email string, roles []string, accessSecret, refreshSecret string, accessTTL, refreshTTL time.Duration) (string, string, error) {
	accessToken, err := generateAccessToken(userID, email, roles, accessSecret, accessTTL)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := generateRefreshToken(userID, refreshSecret, refreshTTL)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// generateAccessToken یک توکن دسترسی با عمر کوتاه تولید می‌کند.
func generateAccessToken(userID, email string, roles []string, secret string, ttl time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "narrative-arch",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// generateRefreshToken یک توکن تازه‌سازی با عمر طولانی تولید می‌کند.
func generateRefreshToken(userID, secret string, ttl time.Duration) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "narrative-arch",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateAccessToken یک Access Token را اعتبارسنجی کرده و Claims آن را برمی‌گرداند.
func ValidateAccessToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid access token")
}

// ValidateRefreshToken یک Refresh Token را اعتبارسنجی کرده و RegisteredClaims آن را برمی‌گرداند.
func ValidateRefreshToken(tokenString, secret string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid refresh token")
}