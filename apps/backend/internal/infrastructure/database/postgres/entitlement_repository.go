// apps/backend/internal/infrastructure/database/postgres/entitlement_repository.go

package postgres

import (
	"context"
	"database/sql"
)

type EntitlementRepository struct {
	db *sql.DB
}

func NewEntitlementRepository(db *sql.DB) *EntitlementRepository {
	return &EntitlementRepository{db: db}
}

// HasAccess بررسی می‌کند که آیا یک کاربر حق دسترسی به یک منبع خاص را دارد یا خیر.
func (r *EntitlementRepository) HasAccess(ctx context.Context, userID string, resourceType string, resourceID string) (bool, error) {
	var exists bool
	query := `
        SELECT EXISTS (
            SELECT 1 FROM entitlements
            WHERE user_id = $1
              AND resource_type = $2
              AND resource_id = $3
              AND status = 'active'
        )
    `
	err := r.db.QueryRowContext(ctx, query, userID, resourceType, resourceID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}