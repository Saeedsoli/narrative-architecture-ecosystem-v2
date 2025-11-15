// apps/backend/internal/domain/moderation/moderation.go

package moderation

import "time"

type QueueStatus string

const (
	StatusPending  QueueStatus = "pending"
	StatusApproved QueueStatus = "approved"
	StatusRejected QueueStatus = "rejected"
)

type QueueItem struct {
	ID        string
	TargetType string
	TargetID   string
	Reason    string
	Flags     map[string]interface{}
	Status    QueueStatus
	CreatedAt time.Time
}