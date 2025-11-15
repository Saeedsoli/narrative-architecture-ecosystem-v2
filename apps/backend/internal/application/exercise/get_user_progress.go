// apps/backend/internal/application/submission/get_user_progress.go

package submission

import (
	"context"
	"narrative-architecture/apps/backend/internal/domain/submission"
)

// ProgressRepository رابطی برای دسترسی به داده‌های پیشرفت است.
type ProgressRepository interface {
	FindProgressByUser(ctx context.Context, userID string) ([]*submission.Progress, error)
}

// GetUserProgressUseCase منطق تجاری برای دریافت پیشرفت کاربر را کپسوله می‌کند.
type GetUserProgressUseCase struct {
	repo ProgressRepository
}

// NewGetUserProgressUseCase یک نمونه جدید ایجاد می‌کند.
func NewGetUserProgressUseCase(repo ProgressRepository) *GetUserProgressUseCase {
	return &GetUserProgressUseCase{repo: repo}
}

// Execute متد اصلی برای اجرای Use Case است.
func (uc *GetUserProgressUseCase) Execute(ctx context.Context, userID string) ([]*submission.Progress, error) {
	return uc.repo.FindProgressByUser(ctx, userID)
}