// apps/backend/internal/domain/article/repository.go

package article

import "context"

// Repository رابطی برای دسترسی به داده‌های مقالات است.
type Repository interface {
	FindByID(ctx context.Context, id string) (*Article, error)
	FindBySlug(ctx context.Context, slug string) (*Article, error)
	Find(ctx context.Context, filter Filter) ([]*Article, int64, error)
	FindTranslations(ctx context.Context, contentGroupID string) ([]Translation, error)
	Save(ctx context.Context, art *Article) error
	Delete(ctx context.Context, id string) error
}

// Filter ساختاری برای فیلتر کردن و صفحه‌بندی مقالات است.
type Filter struct {
	Tags     []string
	Category string
	AuthorID string
	Status   string
	Locale   string
	Page     int
	PageSize int
}