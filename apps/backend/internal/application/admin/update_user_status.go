// apps/backend/internal/application/admin/update_user_status.go

package admin

import (
	"context"
	"errors"
	"narrative-architecture/apps/backend/internal/domain/user"
)

type UpdateUserStatusRequest struct {
	UserID string
	Status user.UserStatus
}

type UpdateUserStatusUseCase struct {
	userRepo user.Repository
}

func NewUpdateUserStatusUseCase(repo user.Repository) *UpdateUserStatusUseCase {
	return &UpdateUserStatusUseCase{userRepo: repo}
}

func (uc *UpdateUserStatusUseCase) Execute(ctx context.Context, req UpdateUserStatusRequest) error {
	u, err := uc.userRepo.FindByID(ctx, req.UserID)
	if err != nil {
		return err
	}

	// اطمینان از اینکه وضعیت جدید معتبر است
	switch req.Status {
	case user.StatusActive, user.StatusSuspended, user.StatusDeleted:
		u.Status = req.Status
		if req.Status == user.StatusDeleted && u.DeletedAt == nil {
			now := time.Now()
			u.DeletedAt = &now
		}
	default:
		return errors.New("invalid user status")
	}

	return uc.userRepo.Save(ctx, u)
}