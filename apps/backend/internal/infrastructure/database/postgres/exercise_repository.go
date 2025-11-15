// apps/backend/internal/infrastructure/database/postgres/exercise_repository.go

package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"narrative-architecture/apps/backend/internal/domain/submission"
)

type ExerciseRepository struct {
	db *sql.DB
}

func NewExerciseRepository(db *sql.DB) *ExerciseRepository {
	return &ExerciseRepository{db: db}
}

// GetExercise یک تمرین را بر اساس ID آن پیدا می‌کند.
func (r *ExerciseRepository) GetExercise(ctx context.Context, id string) (*submission.Exercise, error) {
	var ex submission.Exercise
	var contentJSON []byte

	query := `
        SELECT id, chapter_id, type, difficulty, points, content, created_at, updated_at, deleted_at
        FROM exercises
        WHERE id = $1 AND deleted_at IS NULL
    `
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&ex.ID, &ex.ChapterID, &ex.Type, &ex.Difficulty, &ex.Points, &contentJSON,
		&ex.CreatedAt, &ex.UpdatedAt, &ex.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("exercise not found")
		}
		return nil, err
	}

	if err := json.Unmarshal(contentJSON, &ex.Content); err != nil {
		return nil, err
	}

	// (اختیاری) خواندن Rubric مربوطه
	rubric, err := r.GetRubric(ctx, id)
	if err == nil {
		ex.Rubric = rubric
	}
	
	return &ex, nil
}

// GetRubric معیارهای نمره‌دهی یک تمرین را پیدا می‌کند.
func (r *ExerciseRepository) GetRubric(ctx context.Context, exerciseID string) (*submission.Rubric, error) {
	var ru submission.Rubric
	var criteriaJSON []byte

	query := `SELECT id, exercise_id, criteria, created_at FROM exercise_rubrics WHERE exercise_id = $1`
	err := r.db.QueryRowContext(ctx, query, exerciseID).Scan(
		&ru.ID, &ru.ExerciseID, &criteriaJSON, &ru.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(criteriaJSON, &ru.Criteria); err != nil {
		return nil, err
	}

	return &ru, nil
}