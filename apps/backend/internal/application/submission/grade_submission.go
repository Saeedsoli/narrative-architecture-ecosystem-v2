// apps/backend/internal/application/submission/grade_submission.go

package submission

import (
	"context"
	"errors"
	"time"

	"narrative-architecture/apps/backend/internal/domain/submission"
	"narrative-architecture/apps/backend/internal/domain/user"
)

// GradeSubmissionRequest ساختار درخواست برای امتیازدهی است.
type GradeSubmissionRequest struct {
	SubmissionID string
	GraderID     string
	Score        int
	Feedback     string
}

// GradeSubmissionUseCase منطق تجاری برای امتیازدهی به یک پاسخ را کپسوله می‌کند.
type GradeSubmissionUseCase struct {
	submissionRepo submission.Repository
	userRepo       user.Repository
}

// NewGradeSubmissionUseCase یک نمونه جدید ایجاد می‌کند.
func NewGradeSubmissionUseCase(submissionRepo submission.Repository, userRepo user.Repository) *GradeSubmissionUseCase {
	return &GradeSubmissionUseCase{
		submissionRepo: submissionRepo,
		userRepo:       userRepo,
	}
}

// Execute متد اصلی برای اجرای Use Case است.
func (uc *GradeSubmissionUseCase) Execute(ctx context.Context, req GradeSubmissionRequest) (*submission.Submission, error) {
	// 1. بررسی دسترسی Grader (باید ادمین یا مدرس باشد)
	grader, err := uc.userRepo.FindByID(ctx, req.GraderID)
	if err != nil {
		return nil, errors.New("grader not found")
	}

	isAllowed := false
	for _, role := range grader.Roles {
		if role == "admin" || role == "moderator" {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		return nil, errors.New("user does not have permission to grade submissions")
	}

	// 2. دریافت Submission
	sub, err := uc.submissionRepo.FindByID(ctx, req.SubmissionID)
	if err != nil {
		return nil, err
	}
	
	// 3. دریافت اطلاعات تمرین برای بررسی حداکثر امتیاز
	exercise, err := uc.submissionRepo.GetExercise(ctx, sub.ExerciseID)
	if err != nil {
		return nil, err
	}

	if req.Score < 0 || req.Score > exercise.Points {
		return nil, errors.New("score is out of valid range")
	}

	// 4. آپدیت اطلاعات Submission
	sub.Score = &req.Score
	sub.Feedback = &req.Feedback
	sub.Status = submission.StatusGraded
	now := time.Now()
	sub.GradedAt = &now
	
	// 5. ذخیره در دیتابیس (این متد باید در Repository پیاده‌سازی شود)
	if err := uc.submissionRepo.Update(ctx, sub); err != nil {
		return nil, err
	}
	
	// TODO: ارسال نوتیفیکیشن به کاربر مبنی بر اینکه پاسخ او امتیازدهی شده است.
	
	return sub, nil
}