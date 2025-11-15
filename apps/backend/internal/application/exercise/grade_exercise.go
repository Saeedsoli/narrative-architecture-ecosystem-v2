// apps/backend/internal/application/exercise/get_exercise.go

package exercise

import (
	"context"
	"narrative-architecture/apps/backend/internal/domain/submission"
)

type GetExerciseUseCase struct {
	repo submission.Repository
}

func NewGetExerciseUseCase(repo submission.Repository) *GetExerciseUseCase {
	return &GetExerciseUseCase{repo: repo}
}

func (uc *GetExerciseUseCase) Execute(ctx context.Context, exerciseID string) (*submission.Exercise, error) {
	return uc.repo.GetExercise(ctx, exerciseID)
}