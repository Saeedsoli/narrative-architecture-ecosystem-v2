// apps/backend/internal/infrastructure/database/postgres/token_repository.go

package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/sha3"
)

// TokenRepository پیاده‌سازی رابط برای مدیریت توکن‌ها در PostgreSQL است.
type TokenRepository struct {
	db *sql.DB
}

// NewTokenRepository یک نمونه جدید از TokenRepository ایجاد می‌کند.
func NewTokenRepository(db *sql.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

// hashToken یک توکن را برای ذخیره‌سازی امن در دیتابیس هش می‌کند.
func hashToken(token string) string {
	hasher := sha3.New256()
	hasher.Write([]byte(token))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

// SaveRefreshToken یک Refresh Token جدید را در دیتابیس ذخیره می‌کند.
func (r *TokenRepository) SaveRefreshToken(ctx context.Context, userID, token string, ttl time.Duration) error {
	tokenHash := hashToken(token)
	expiresAt := time.Now().Add(ttl)

	query := `
        INSERT INTO auth_tokens (id, user_id, type, token_hash, expires_at)
        VALUES ($1, $2, 'refresh', $3, $4)
    `
	_, err := r.db.ExecContext(ctx, query, ulid.New().String(), userID, tokenHash, expiresAt)
	return err
}

// ValidateRefreshToken بررسی می‌کند که آیا Refresh Token معتبر است یا خیر.
// یک توکن معتبر است اگر در دیتابیس وجود داشته باشد و هنوز استفاده نشده باشد (used_at IS NULL).
func (r *TokenRepository) ValidateRefreshToken(ctx context.Context, userID, token string) error {
	tokenHash := hashToken(token)
	var usedAt sql.NullTime

	query := `
        SELECT used_at FROM auth_tokens
        WHERE user_id = $1 AND token_hash = $2 AND type = 'refresh' AND expires_at > NOW()
    `
	err := r.db.QueryRowContext(ctx, query, userID, tokenHash).Scan(&usedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("invalid or expired refresh token")
		}
		return err
	}

	if usedAt.Valid {
		return errors.New("refresh token has already been used")
	}

	return nil
}

// RotateRefreshToken توکن قدیمی را باطل کرده و توکن جدید را ذخیره می‌کند (در یک تراکنش).
func (r *TokenRepository) RotateRefreshToken(ctx context.Context, userID, oldToken, newToken string, ttl time.Duration) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // در صورت بروز خطا، تراکنش را بازگردانی کن

	// 1. توکن قدیمی را به‌عنوان "استفاده شده" علامت بزن
	oldTokenHash := hashToken(oldToken)
	res, err := tx.ExecContext(ctx,
		`UPDATE auth_tokens SET used_at = NOW() WHERE user_id = $1 AND token_hash = $2 AND used_at IS NULL`,
		userID, oldTokenHash,
	)
	if err != nil {
		return err
	}

	// اگر هیچ ردیفی آپدیت نشد، یعنی توکن قبلاً استفاده شده یا نامعتبر است.
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("failed to rotate token: old token is invalid or already used")
	}

	// 2. توکن جدید را ذخیره کن
	newTokenHash := hashToken(newToken)
	expiresAt := time.Now().Add(ttl)
	_, err = tx.ExecContext(ctx,
		`INSERT INTO auth_tokens (id, user_id, type, token_hash, expires_at) VALUES ($1, $2, 'refresh', $3, $4)`,
		ulid.New().String(), userID, newTokenHash, expiresAt,
	)
	if err != nil {
		return err
	}

	// 3. تراکنش را نهایی کن
	return tx.Commit()
}

// InvalidateRefreshToken یک Refresh Token را باطل می‌کند (برای Logout).
func (r *TokenRepository) InvalidateRefreshToken(ctx context.Context, token string) error {
	tokenHash := hashToken(token)
	query := `UPDATE auth_tokens SET used_at = NOW() WHERE token_hash = $1 AND type = 'refresh' AND used_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, tokenHash)
	return err
}