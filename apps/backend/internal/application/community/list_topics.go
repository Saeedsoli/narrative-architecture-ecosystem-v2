// apps/backend/internal/application/community/list_topics.go

package community

import (
	"context"
	"narrative-architecture/apps/backend/internal/domain/community"
)

type ListTopicsResponse struct {
	Topics     []*community.Topic
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

type ListTopicsUseCase struct {
	repo community.Repository
}

func NewListTopicsUseCase(repo community.Repository) *ListTopicsUseCase {
	return &ListTopicsUseCase{repo: repo}
}

func (uc *ListTopicsUseCase) Execute(ctx context.Context, locale string, page, pageSize int) (*ListTopicsResponse, error) {
	if page <= 0 { page = 1 }
	if pageSize <= 0 { pageSize = 10 }

	topics, total, err := uc.repo.ListTopics(ctx, locale, page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := 0
	if total > 0 {
		totalPages = (int(total) + pageSize - 1) / pageSize
	}

	return &ListTopicsResponse{
		Topics:     topics,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}