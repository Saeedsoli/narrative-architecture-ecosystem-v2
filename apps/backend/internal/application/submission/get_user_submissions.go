// apps/backend/internal/application/submission/get_user_submissions.go

package submission

import (
	"context"
	"errors"

	"narrative-architecture/apps/backend/internal/domain/submission"
)

type GetUserSubmissionsRequest struct {
	UserID     string
	ExerciseID string
}

type GetUserSubmissionsUseCase struct {
	repo submission.Repository
}

func NewGetUserSubmissionsUseCase(repo submission.Repository) *GetUserSubmissionsUseCase {
	return &GetUserSubmissionsUseCase{repo: repo}
}

func (uc *GetUserSubmissionsUseCase) Execute(ctx context.Context, req GetUserSubmissionsRequest) ([]*submission.Submission, error) {
	if req.UserID == "" || req.ExerciseID == "" {
		return nil, errors.New("userID and exerciseID are required")
	}

	return uc.repo.FindByUserAndExercise(ctx, req.UserID, req.ExerciseID)
}