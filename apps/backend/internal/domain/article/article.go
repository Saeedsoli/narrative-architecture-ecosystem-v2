package article

import "time"

// ArticleStatus نشان‌دهنده وضعیت یک مقاله است.
type ArticleStatus string

const (
	StatusDraft     ArticleStatus = "draft"
	StatusReview    ArticleStatus = "review"
	StatusPublished ArticleStatus = "published"
	StatusScheduled ArticleStatus = "scheduled"
)

// Translation اطلاعات یک نسخه ترجمه شده را نگه می‌دارد.
type Translation struct {
	Locale string `json:"locale"`
	Slug   string `json:"slug"`
}

// Article موجودیت اصلی مقاله در لایه Domain است.
type Article struct {
	ID             string
	ContentGroupID string
	Locale         string
	Slug           string
	Title          string
	Excerpt        string
	Content        string
	CoverImage     struct {
		URL string
		Alt string
	}
	Author struct {
		ID     string
		Name   string
		Avatar string
	}
	Metadata struct {
		Tags       []string
		Category   string
		ReadTime   int
		Difficulty string
	}
	Status      ArticleStatus
	PublishedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
	Translations []Translation
}

// IsPublished بررسی می‌کند که آیا مقاله برای نمایش عمومی آماده است یا خیر.
func (a *Article) IsPublished() bool {
	return a.Status == StatusPublished && a.PublishedAt != nil && a.PublishedAt.Before(time.Now())
}