package postgres

import (
	"context"
	"database/sql"
)

type EntitlementRepository struct {
	db *sql.DB
}

func (r *EntitlementRepository) HasAccess(ctx context.Context, userID, resourceType, resourceID string) (bool, error) {
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