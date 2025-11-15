// apps/backend/internal/application/submission/submit_exercise.go

package submission

import (
	"context"
	"errors"
	"time"

	"github.com/oklog/ulid/v2"
	"narrative-architecture/apps/backend/internal/domain/submission"
	"narrative-architecture/apps/backend/internal/domain/user"
)

// SubmitExerciseRequest ساختار درخواست برای ارسال پاسخ است.
type SubmitExerciseRequest struct {
	ExerciseID string
	Answer     map[string]interface{}
	UserID     string
}

// SubmitExerciseUseCase منطق تجاری برای ارسال پاسخ را کپسوله می‌کند.
type SubmitExerciseUseCase struct {
	submissionRepo submission.Repository
	userRepo       user.Repository
	analyzeUC      *AnalyzeSubmissionUseCase
}

// NewSubmitExerciseUseCase یک نمونه جدید ایجاد می‌کند.
func NewSubmitExerciseUseCase(
	submissionRepo submission.Repository,
	userRepo user.Repository,
	analyzeUC *AnalyzeSubmissionUseCase,
) *SubmitExerciseUseCase {
	return &SubmitExerciseUseCase{
		submissionRepo: submissionRepo,
		userRepo:       userRepo,
		analyzeUC:      analyzeUC,
	}
}

// Execute متد اصلی برای اجرای Use Case است.
func (uc *SubmitExerciseUseCase) Execute(ctx context.Context, req SubmitExerciseRequest) (*submission.Submission, error) {
	// 1. بررسی وجود تمرین
	_, err := uc.submissionRepo.GetExercise(ctx, req.ExerciseID)
	if err != nil {
		return nil, errors.New("exercise not found")
	}

	// 2. بررسی وجود کاربر
	_, err = uc.userRepo.FindByID(ctx, req.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 3. ایجاد سند ارسال
	newSubmission := &submission.Submission{
		ID:          ulid.New().String(),
		ExerciseID:  req.ExerciseID,
		UserID:      req.UserID,
		Status:      submission.StatusPending,
		Answer:      req.Answer,
		SubmittedAt: time.Now(),
	}

	// 4. ذخیره در دیتابیس
	if err := uc.repo.Create(ctx, newSubmission); err != nil {
		return nil, err
	}

	// 5. آغاز فرآیند تحلیل AI (به صورت Async)
	go func() {
		// ایجاد یک context جدید برای goroutine
		bgCtx := context.Background()
		// ما اهمیتی به خطای این بخش نمی‌دهیم، چون یک عملیات پس‌زمینه است.
		// خطاها در داخل خود Use Case لاگ می‌شوند.
		_ = uc.analyzeUC.Execute(bgCtx, newSubmission.ID, req.UserID)
	}()

	return newSubmission, nil
}