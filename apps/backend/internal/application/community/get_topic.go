// apps/backend/internal/application/community/get_topic.go

package community

import (
	"context"
	"narrative-architecture/apps/backend/internal/domain/community"
)

type GetTopicUseCase struct {
	repo community.Repository
}

func NewGetTopicUseCase(repo community.Repository) *GetTopicUseCase {
	return &GetTopicUseCase{repo: repo}
}

func (uc *GetTopicUseCase) Execute(ctx context.Context, topicID string) (*community.Topic, error) {
	return uc.repo.FindTopicByID(ctx, topicID)
}