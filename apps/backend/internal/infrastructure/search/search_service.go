// apps/backend/internal/infrastructure/search/search_service.go

package search

import (
	"context"
	"narrative-architecture/apps/backend/internal/domain/article"
)

type Service interface {
	SearchArticles(ctx context.Context, query, locale string, page, pageSize int) ([]*article.Article, int64, error)
}