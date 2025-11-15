// apps/backend/internal/application/community/list_posts.go

package community

import (
	"context"
	"narrative-architecture/apps/backend/internal/domain/community"
)

type ListPostsResponse struct {
	Posts      []*community.Post
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

type ListPostsUseCase struct {
	repo community.Repository
}

func NewListPostsUseCase(repo community.Repository) *ListPostsUseCase {
	return &ListPostsUseCase{repo: repo}
}

func (uc *ListPostsUseCase) Execute(ctx context.Context, topicID string, page, pageSize int) (*ListPostsResponse, error) {
	if page <= 0 { page = 1 }
	if pageSize <= 0 { pageSize = 20 }

	posts, total, err := uc.repo.FindPostsByTopicID(ctx, topicID, page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := 0
	if total > 0 {
		totalPages = (int(total) + pageSize - 1) / pageSize
	}

	return &ListPostsResponse{
		Posts:      posts,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}