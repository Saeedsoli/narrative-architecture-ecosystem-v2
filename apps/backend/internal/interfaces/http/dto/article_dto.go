package dto

import (
	"time"
	"narrative-architecture/apps/backend/internal/domain/article"
)

// ArticleResponse ساختار پاسخ API برای یک مقاله است.
type ArticleResponse struct {
	ID           string               `json:"id"`
	Locale       string               `json:"locale"`
	Slug         string               `json:"slug"`
	Title        string               `json:"title"`
	Excerpt      string               `json:"excerpt"`
	Content      string               `json:"content"`
	CoverImage   interface{}          `json:"coverImage"`
	Author       interface{}          `json:"author"`
	Metadata     interface{}          `json:"metadata"`
	PublishedAt  *time.Time           `json:"publishedAt"`
	UpdatedAt    time.Time            `json:"updatedAt"`
	Translations []article.Translation `json:"translations"`
}

// ToArticleResponse موجودیت Domain را به DTO پاسخ تبدیل می‌کند.
func ToArticleResponse(art *article.Article) *ArticleResponse {
	return &ArticleResponse{
		ID:           art.ID,
		Locale:       art.Locale,
		Slug:         art.Slug,
		Title:        art.Title,
		Excerpt:      art.Excerpt,
		Content:      art.Content,
		CoverImage:   art.CoverImage,
		Author:       art.Author,
		Metadata:     art.Metadata,
		PublishedAt:  art.PublishedAt,
		UpdatedAt:    art.UpdatedAt,
		Translations: art.Translations,
	}
}