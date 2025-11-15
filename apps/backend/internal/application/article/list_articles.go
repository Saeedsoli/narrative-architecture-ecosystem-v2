package article

import (
	"context"
	"narrative-architecture/apps/backend/internal/domain/article"
)

type ListArticlesResponse struct {
	Articles   []*article.Article
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

type ListArticlesUseCase struct {
	repo article.Repository
}

func NewListArticlesUseCase(repo article.Repository) *ListArticlesUseCase {
	return &ListArticlesUseCase{repo: repo}
}

func (uc *ListArticlesUseCase) Execute(ctx context.Context, filter article.Filter) (*ListArticlesResponse, error) {
	if filter.Page <= 0 { filter.Page = 1 }
	if filter.PageSize <= 0 { filter.PageSize = 10 }

	articles, total, err := uc.repo.Find(ctx, filter)
	if err != nil { return nil, err }

	totalPages := 0
	if total > 0 { totalPages = (int(total) + filter.PageSize - 1) / filter.PageSize }

	return &ListArticlesResponse{
		Articles:   articles,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}, nil
}