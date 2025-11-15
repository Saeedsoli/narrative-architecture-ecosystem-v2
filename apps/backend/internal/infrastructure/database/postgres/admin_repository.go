// apps/backend/internal/infrastructure/database/postgres/admin_repository.go

package postgres

import (
	"context"
	"database/sql"
	"narrative-architecture/apps/backend/internal/domain/user"
)

type AdminRepository struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

// FindAllUsers لیستی از تمام کاربران را با صفحه‌بندی برمی‌گرداند.
func (r *AdminRepository) FindAllUsers(ctx context.Context, page, pageSize int) ([]*user.User, int64, error) {
	var total int64
	countQuery := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
        SELECT u.id, u.email, u.email_verified, up.full_name, up.username, u.status, u.created_at, u.updated_at
        FROM users u
        LEFT JOIN user_profiles up ON u.id = up.user_id
        WHERE u.deleted_at IS NULL
        ORDER BY u.created_at DESC
        LIMIT $1 OFFSET $2
    `
	offset := (page - 1) * pageSize
	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*user.User
	for rows.Next() {
		var u user.User
		if err := rows.Scan(
			&u.ID, &u.Email, &u.EmailVerified, &u.FullName, &u.Username, &u.Status, &u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		users = append(users, &u)
	}

	return users, total, nil
}

// UpdateUserStatus وضعیت یک کاربر را آپدیت می‌کند.
func (r *AdminRepository) UpdateUserStatus(ctx context.Context, userID string, status user.UserStatus) error {
	query := `UPDATE users SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, userID)
	return err
}