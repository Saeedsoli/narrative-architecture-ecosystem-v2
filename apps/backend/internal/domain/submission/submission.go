// apps/backend/internal/domain/submission/submission.go

package submission

import (
	"time"
	"narrative-architecture/apps/backend/internal/domain/entitlement"
)

// SubmissionStatus وضعیت پاسخ تمرین است.
type SubmissionStatus string

const (
	StatusPending  SubmissionStatus = "pending"
	StatusGraded   SubmissionStatus = "graded"
	StatusReviewed SubmissionStatus = "reviewed"
)

// Submission موجودیت اصلی برای ذخیره پاسخ‌های تمرینی کاربران است.
type Submission struct {
	ID           string
	ExerciseID   string
	UserID       string
	Status       SubmissionStatus
	Answer       map[string]interface{} // ساختار انعطاف‌پذیر برای انواع تمرین (متنی، صوتی و ...)
	Score        *int
	Feedback     *string
	AIAnalysis   *string
	AIReport     map[string]interface{} // گزارش تحلیل سیستم AI
	SubmittedAt  time.Time
	GradedAt     *time.Time
	ReviewedAt   *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

// ExerciseType نوع تمرین را مشخص می‌کند.
type ExerciseType string

const (
	ExerciseTypeMCQ      ExerciseType = "mcq"
	ExerciseTypeEssay    ExerciseType = "essay"
	ExerciseTypePractical ExerciseType = "practical"
)

// Exercise ساختار تمرین است.
type Exercise struct {
	ID          string
	ChapterID   string
	Type        ExerciseType
	Difficulty  string
	Points      int
	Content     map[string]interface{} // محتوای تمرین
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
	Rubric      *Rubric
}

// Rubric معیارهای ارزیابی تمرین است.
type Rubric struct {
	ID       string
	ExerciseID string
	Criteria []RubricCriteria
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// RubricCriteria معیارهای خاص ارزیابی است.
type RubricCriteria struct {
	Name        string
	Points      int
	Description string
}

// ExerciseResult نتیجه نهایی تمرین را نشان می‌دهد.
type ExerciseResult struct {
	ExerciseID string
	UserID     string
	Score      int
	MaxScore   int
	Feedback   string
	AIAnalysis string
}