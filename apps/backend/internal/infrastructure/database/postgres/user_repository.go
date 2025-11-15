// apps/backend/internal/infrastructure/database/postgres/user_repository.go

package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"narrative-architecture/apps/backend/internal/domain/user"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Save یک کاربر جدید را ایجاد یا کاربر موجود را آپدیت می‌کند.
func (r *UserRepository) Save(ctx context.Context, u *user.User) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Upsert در جدول users
	userQuery := `
        INSERT INTO users (id, email, email_verified, password_hash, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (id) DO UPDATE SET
            email = EXCLUDED.email,
            password_hash = EXCLUDED.password_hash,
            status = EXCLUDED.status,
            updated_at = NOW()
    `
	_, err = tx.ExecContext(ctx, userQuery, u.ID, u.Email, u.EmailVerified, u.PasswordHash, u.Status, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return err
	}

	// Upsert در جدول user_profiles
	profileQuery := `
        INSERT INTO user_profiles (user_id, full_name, username)
        VALUES ($1, $2, $3)
        ON CONFLICT (user_id) DO UPDATE SET
            full_name = EXCLUDED.full_name,
            username = EXCLUDED.username,
            updated_at = NOW()
    `
	_, err = tx.ExecContext(ctx, profileQuery, u.ID, u.FullName, u.Username)
	if err != nil {
		return err
	}
	
	// همگام‌سازی نقش‌ها (حذف نقش‌های قدیمی، افزودن نقش‌های جدید)
	_, err = tx.ExecContext(ctx, `DELETE FROM user_roles WHERE user_id = $1`, u.ID)
	if err != nil {
		return err
	}

	if len(u.Roles) > 0 {
		roleQuery := `
            INSERT INTO user_roles (user_id, role_id)
            SELECT $1, id FROM roles WHERE name = ANY($2::text[])
        `
		_, err = tx.ExecContext(ctx, roleQuery, u.ID, u.Roles)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// FindByID یک کاربر را بر اساس ID پیدا می‌کند.
func (r *UserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	query := `
        SELECT u.id, u.email, u.email_verified, u.password_hash, up.username, up.full_name, u.status, u.created_at, u.updated_at
        FROM users u
        LEFT JOIN user_profiles up ON u.id = up.user_id
        WHERE u.id = $1 AND u.deleted_at IS NULL
    `
	u, err := r.scanUser(r.db.QueryRowContext(ctx, query, id))
	if err != nil {
		return nil, err
	}
	
	roles, err := r.loadRoles(ctx, u.ID)
	if err != nil {
		return nil, err
	}
	u.Roles = roles

	return u, nil
}

// FindByEmail یک کاربر را بر اساس ایمیل پیدا می‌کند.
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
        SELECT u.id, u.email, u.email_verified, u.password_hash, up.username, up.full_name, u.status, u.created_at, u.updated_at
        FROM users u
        LEFT JOIN user_profiles up ON u.id = up.user_id
        WHERE u.email = $1 AND u.deleted_at IS NULL
    `
	u, err := r.scanUser(r.db.QueryRowContext(ctx, query, email))
	if err != nil {
		return nil, err
	}

	roles, err := r.loadRoles(ctx, u.ID)
	if err != nil {
		return nil, err
	}
	u.Roles = roles
	
	return u, nil
}

// scanUser یک تابع کمکی برای اسکن کردن ردیف‌ها است.
func (r *UserRepository) scanUser(row *sql.Row) (*user.User, error) {
	var u user.User
	err := row.Scan(
		&u.ID, &u.Email, &u.EmailVerified, &u.PasswordHash,
		&u.Username, &u.FullName, &u.Status, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &u, nil
}

// loadRoles نقش‌های یک کاربر را از دیتابیس می‌خواند.
func (r *UserRepository) loadRoles(ctx context.Context, userID string) ([]string, error) {
	query := `SELECT r.name FROM roles r JOIN user_roles ur ON r.id = ur.role_id WHERE ur.user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var roleName string
		if err := rows.Scan(&roleName); err != nil {
			return nil, err
		}
		roles = append(roles, roleName)
	}
	return roles, nil
}