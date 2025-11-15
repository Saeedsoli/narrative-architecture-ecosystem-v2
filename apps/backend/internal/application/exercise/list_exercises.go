// apps/backend/internal/application/exercise/list_exercises.go

package exercise

import (
	"context"
	"narrative-architecture/apps/backend/internal/domain/submission"
)

type ListExercisesUseCase struct {
	repo submission.Repository
}

func NewListExercisesUseCase(repo submission.Repository) *ListExercisesUseCase {
	return &ListExercisesUseCase{repo: repo}
}

func (uc *ListExercisesUseCase) Execute(ctx context.Context, chapterID string) ([]*submission.Exercise, error) {
	// این متد باید در Repository پیاده‌سازی شود
	return uc.repo.FindExercisesByChapter(ctx, chapterID)
}