// apps/backend/internal/domain/admin/repository.go

package admin

import (
	"context"
	"narrative-architecture/apps/backend/internal/domain/moderation"
	"narrative-architecture/apps/backend/internal/domain/user"
)

type UserRepository interface {
	FindAll(ctx context.Context, page, pageSize int) ([]*user.User, int64, error)
	UpdateUserStatus(ctx context.Context, userID string, status user.UserStatus) error
}

type ModerationRepository interface {
	FindQueueItems(ctx context.Context, status string, page, pageSize int) ([]*moderation.QueueItem, int64, error)
	FindQueueItemByID(ctx context.Context, id string) (*moderation.QueueItem, error)
	UpdateQueueItemStatus(ctx context.Context, id, moderatorID, reason string, status moderation.QueueStatus) error
}