// apps/backend/internal/infrastructure/database/postgres/article_repository.go

package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"narrative-architecture/apps/backend/internal/domain/article"
	"github.com/oklog/ulid/v2"
)

// ArticleRepository پیاده‌سازی رابط برای عملیات مربوط به مقالات در PostgreSQL است.
// این Repository عمدتاً برای مدیریت داده‌های رابطه‌ای مانند بوکمارک‌ها استفاده می‌شود.
type ArticleRepository struct {
	db *sql.DB
}

// NewArticleRepository یک نمونه جدید از ArticleRepository ایجاد می‌کند.
func NewArticleRepository(db *sql.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

// Create یک بوکمارک جدید در دیتابیس ذخیره می‌کند.
// از ON CONFLICT DO NOTHING برای جلوگیری از خطای تکراری استفاده می‌شود.
func (r *ArticleRepository) Create(ctx context.Context, b *article.Bookmark) error {
	query := `
        INSERT INTO article_bookmarks (id, user_id, article_id, created_at)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (user_id, article_id) DO NOTHING
    `
	_, err := r.db.ExecContext(ctx, query, ulid.New().String(), b.UserID, b.ArticleID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to add bookmark: %w", err)
	}
	return nil
}

// Delete یک بوکمارک را بر اساس کاربر و شناسه مقاله حذف می‌کند.
func (r *ArticleRepository) Delete(ctx context.Context, userID, articleID string) error {
	query := `DELETE FROM article_bookmarks WHERE user_id = $1 AND article_id = $2`
	res, err := r.db.ExecContext(ctx, query, userID, articleID)
	if err != nil {
		return fmt.Errorf("failed to remove bookmark: %w", err)
	}
	
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows after delete: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("bookmark not found or already deleted")
	}
	return nil
}

// FindByUser لیستی از تمام بوکمارک‌های یک کاربر را برمی‌گرداند.
// این لیست بر اساس تاریخ ایجاد (جدیدترین) مرتب شده است.
func (r *ArticleRepository) FindByUser(ctx context.Context, userID string) ([]*article.Bookmark, error) {
	query := `
        SELECT id, user_id, article_id, created_at
        FROM article_bookmarks
        WHERE user_id = $1
        ORDER BY created_at DESC
    `
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user bookmarks: %w", err)
	}
	defer rows.Close()

	var bookmarks []*article.Bookmark
	for rows.Next() {
		var b article.Bookmark
		if err := rows.Scan(&b.ID, &b.UserID, &b.ArticleID, &b.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan bookmark row: %w", err)
		}
		bookmarks = append(bookmarks, &b)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return bookmarks, nil
}

// Exists بررسی می‌کند که آیا یک کاربر یک مقاله خاص را بوکمارک کرده است یا خیر.
func (r *ArticleRepository) Exists(ctx context.Context, userID, articleID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM article_bookmarks WHERE user_id = $1 AND article_id = $2)`
	err := r.db.QueryRowContext(ctx, query, userID, articleID).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil // Should not happen with EXISTS, but for safety
		}
		return false, fmt.Errorf("failed to check if bookmark exists: %w", err)
	}
	return exists, nil
}