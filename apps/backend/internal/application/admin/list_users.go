// apps/backend/internal/application/admin/list_users.go

package admin

import (
	"context"
	"narrative-architecture/apps/backend/internal/domain/user"
)

type ListUsersUseCase struct {
	userRepo user.Repository
}

func NewListUsersUseCase(repo user.Repository) *ListUsersUseCase {
	return &ListUsersUseCase{userRepo: repo}
}

// Execute لیستی از کاربران را با صفحه‌بندی برمی‌گرداند.
func (uc *ListUsersUseCase) Execute(ctx context.Context, page, pageSize int) ([]*user.User, int64, error) {
	// این متد باید در UserRepository پیاده‌سازی شود.
	return uc.userRepo.FindAll(ctx, page, pageSize)
}