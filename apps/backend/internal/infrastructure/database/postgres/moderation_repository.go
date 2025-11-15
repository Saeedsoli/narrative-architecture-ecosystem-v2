// apps/backend/internal/infrastructure/database/postgres/moderation_repository.go

package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
	"narrative-architecture/apps/backend/internal/domain/moderation"
)

type ModerationRepository struct {
	db *sql.DB
}

func NewModerationRepository(db *sql.DB) *ModerationRepository {
	return &ModerationRepository{db: db}
}

// FindQueueItems لیستی از آیتم‌های در صف بررسی را برمی‌گرداند.
func (r *ModerationRepository) FindQueueItems(ctx context.Context, status string, page, pageSize int) ([]*moderation.QueueItem, int64, error) {
	// ... (پیاده‌سازی مشابه FindAllUsers)
	return nil, 0, nil
}

// FindQueueItemByID یک آیتم را بر اساس ID پیدا می‌کند.
func (r *ModerationRepository) FindQueueItemByID(ctx context.Context, id string) (*moderation.QueueItem, error) {
	// ... (پیاده‌سازی کوئری SELECT)
	return nil, nil
}

// UpdateQueueItemStatus وضعیت یک آیتم را آپدیت می‌کند.
func (r *ModerationRepository) UpdateQueueItemStatus(ctx context.Context, id, moderatorID, reason string, status moderation.QueueStatus) error {
	query := `
        UPDATE moderation_queue
        SET status = $1, assigned_to = $2, resolution_reason = $3, resolved_at = $4, updated_at = NOW()
        WHERE id = $5
    `
	_, err := r.db.ExecContext(ctx, query, status, moderatorID, reason, time.Now(), id)
	return err
}