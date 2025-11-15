// apps/backend/internal/interfaces/http/dto/submission_dto.go

package dto

import (
	"time"
	"narrative-architecture/apps/backend/internal/domain/submission"
)

type SubmissionResponse struct {
	ID          string                 `json:"id"`
	ExerciseID  string                 `json:"exerciseId"`
	UserID      string                 `json:"userId"`
	Status      string                 `json:"status"`
	Answer      map[string]interface{} `json:"answer"`
	Score       *int                   `json:"score,omitempty"`
	Feedback    *string                `json:"feedback,omitempty"`
	AISummary   *string                `json:"aiSummary,omitempty"`
	SubmittedAt time.Time              `json:"submittedAt"`
	GradedAt    *time.Time             `json:"gradedAt,omitempty"`
}

func ToSubmissionResponse(s *submission.Submission) *SubmissionResponse {
	return &SubmissionResponse{
		ID:          s.ID,
		ExerciseID:  s.ExerciseID,
		UserID:      s.UserID,
		Status:      string(s.Status),
		Answer:      s.Answer,
		Score:       s.Score,
		Feedback:    s.Feedback,
		AISummary:   s.AISummary,
		SubmittedAt: s.SubmittedAt,
		GradedAt:    s.GradedAt,
	}
}