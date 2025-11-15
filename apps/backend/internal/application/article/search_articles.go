// apps/backend/internal/application/article/search_articles.go

package article

import (
	"context"
)

// SearchService رابطی برای موتور جستجو است.
type SearchService interface {
	SearchArticles(ctx context.Context, query string, locale string, page, pageSize int) ([]*Article, int64, error)
}

// SearchArticlesUseCase منطق تجاری برای جستجوی مقالات را کپسوله می‌کند.
type SearchArticlesUseCase struct {
	searcher SearchService
}

// NewSearchArticlesUseCase یک نمونه جدید از SearchArticlesUseCase ایجاد می‌کند.
func NewSearchArticlesUseCase(searcher SearchService) *SearchArticlesUseCase {
	return &SearchArticlesUseCase{searcher: searcher}
}

// Execute متد اصلی برای اجرای جستجو است.
func (uc *SearchArticlesUseCase) Execute(ctx context.Context, query, locale string, page, pageSize int) (*ListArticlesResponse, error) {
	if page <= 0 { page = 1 }
	if pageSize <= 0 { pageSize = 10 }

	articles, total, err := uc.searcher.SearchArticles(ctx, query, locale, page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := 0
	if total > 0 {
		totalPages = (int(total) + pageSize - 1) / pageSize
	}

	return &ListArticlesResponse{
		Articles:   articles,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}