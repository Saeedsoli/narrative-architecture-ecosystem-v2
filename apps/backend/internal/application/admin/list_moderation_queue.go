// apps/backend/internal/application/admin/list_moderation_queue.go

package admin

import (
	"context"
	"narrative-architecture/apps/backend/internal/domain/moderation" // فرض بر وجود این پکیج
)

type ListModerationQueueUseCase struct {
	repo moderation.Repository
}

func NewListModerationQueueUseCase(repo moderation.Repository) *ListModerationQueueUseCase {
	return &ListModerationQueueUseCase{repo: repo}
}

func (uc *ListModerationQueueUseCase) Execute(ctx context.Context, status string, page, pageSize int) ([]*moderation.QueueItem, int64, error) {
	return uc.repo.FindQueueItems(ctx, status, page, pageSize)
}