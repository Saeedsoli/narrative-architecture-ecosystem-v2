// apps/backend/internal/domain/article/bookmark.go

package article

import (
	"context"
	"time"
)

// Bookmark نشان‌دهنده یک بوکمارک است که توسط یک کاربر برای یک مقاله ایجاد شده.
type Bookmark struct {
	ID        string
	UserID    string
	ArticleID string
	CreatedAt time.Time
}

// NewBookmark یک نمونه جدید از Bookmark ایجاد می‌کند.
func NewBookmark(id, userID, articleID string) *Bookmark {
	return &Bookmark{
		ID:        id,
		UserID:    userID,
		ArticleID: articleID,
		CreatedAt: time.Now(),
	}
}

// BookmarkRepository رابطی برای عملیات مربوط به بوکمارک‌ها است.
type BookmarkRepository interface {
	Create(ctx context.Context, b *Bookmark) error
	Delete(ctx context.Context, userID, articleID string) error
	FindByUser(ctx context.Context, userID string) ([]*Bookmark, error)
	Exists(ctx context.Context, userID, articleID string) (bool, error)
}