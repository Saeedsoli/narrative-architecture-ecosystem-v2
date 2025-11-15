// apps/backend/internal/domain/user/repository.go

package user

import (
	"context"
	"time"
)

// UserRepository رابطی برای دسترسی به داده‌های کاربران است.
type Repository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
	Save(ctx context.Context, u *User) error
}

// TokenRepository رابطی برای مدیریت توکن‌های احراز هویت است.
type TokenRepository interface {
	// SaveRefreshToken یک Refresh Token جدید را برای یک کاربر ذخیره می‌کند.
	SaveRefreshToken(ctx context.Context, userID, token string, ttl time.Duration) error

	// ValidateRefreshToken بررسی می‌کند که آیا توکن ارائه شده معتبر و قابل استفاده است یا خیر.
	ValidateRefreshToken(ctx context.Context, userID, token string) error

	// RotateRefreshToken توکن قدیمی را باطل کرده و توکن جدید را جایگزین آن می‌کند.
	RotateRefreshToken(ctx context.Context, userID, oldToken, newToken string, ttl time.Duration) error

	// InvalidateRefreshToken یک Refresh Token را باطل می‌کند (برای Logout).
	InvalidateRefreshToken(ctx context.Context, token string) error
}