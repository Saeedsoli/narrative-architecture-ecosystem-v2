// apps/backend/internal/domain/submission/repository.go

package submission

import (
	"context"
)

type Repository interface {
	// Create یک Submission جدید در دیتابیس ایجاد می‌کند.
	Create(ctx context.Context, sub *Submission) error

	// FindByID یک Submission را بر اساس ID پیدا می‌کند.
	FindByID(ctx context.Context, id string) (*Submission, error)

	// FindByUserAndExercise لیستی از تمام ارسال‌های یک کاربر برای یک تمرین خاص را برمی‌گرداند.
	FindByUserAndExercise(ctx context.Context, userID, exerciseID string) ([]*Submission, error)

	// UpdateAISummary خلاصه تحلیل AI و شناسه لاگ مربوطه را برای یک Submission آپدیت می‌کند.
	UpdateAISummary(ctx context.Context, id, summary, aiLogID string) error

	// UpdateStatus وضعیت یک Submission را آپدیت می‌کند.
	UpdateStatus(ctx context.Context, id string, status SubmissionStatus) error

	// GetExercise یک تمرین را بر اساس ID آن از دیتابیس می‌خواند.
	GetExercise(ctx context.Context, id string) (*Exercise, error)
}