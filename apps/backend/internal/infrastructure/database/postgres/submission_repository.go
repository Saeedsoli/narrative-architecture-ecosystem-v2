// apps/backend/internal/infrastructure/database/postgres/submission_repository.go

package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"narrative-architecture/apps/backend/internal/domain/submission"
)

// SubmissionRepository پیاده‌سازی رابط برای عملیات مربوط به ارسال‌ها در PostgreSQL است.
type SubmissionRepository struct {
	db *sql.DB
}

// NewSubmissionRepository یک نمونه جدید از SubmissionRepository ایجاد می‌کند.
func NewSubmissionRepository(db *sql.DB) *SubmissionRepository {
	return &SubmissionRepository{db: db}
}

// Create یک Submission جدید را در دیتابیس ایجاد می‌کند.
func (r *SubmissionRepository) Create(ctx context.Context, sub *submission.Submission) error {
	answerJSON, err := json.Marshal(sub.Answer)
	if err != nil {
		return fmt.Errorf("failed to marshal submission answer: %w", err)
	}

	query := `
        INSERT INTO submissions (id, exercise_id, user_id, status, answer, submitted_at, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
    `
	_, err = r.db.ExecContext(ctx, query, sub.ID, sub.ExerciseID, sub.UserID, sub.Status, answerJSON, sub.SubmittedAt)
	return err
}

// FindByID یک Submission را بر اساس ID پیدا می‌کند.
func (r *SubmissionRepository) FindByID(ctx context.Context, id string) (*submission.Submission, error) {
	var s submission.Submission
	var answerJSON []byte
	var score sql.NullInt64
	var feedback, aiSummary sql.NullString
	var gradedAt, deletedAt sql.NullTime
	var aiLogID sql.NullString

	query := `
        SELECT id, exercise_id, user_id, status, answer, score, feedback, ai_summary, ai_log_id,
               submitted_at, graded_at, created_at, updated_at, deleted_at
        FROM submissions
        WHERE id = $1 AND deleted_at IS NULL
    `
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID, &s.ExerciseID, &s.UserID, &s.Status, &answerJSON, &score, &feedback, &aiSummary, &aiLogID,
		&s.SubmittedAt, &gradedAt, &s.CreatedAt, &s.UpdatedAt, &deletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("submission not found")
		}
		return nil, err
	}

	if err := json.Unmarshal(answerJSON, &s.Answer); err != nil {
		return nil, fmt.Errorf("failed to unmarshal answer: %w", err)
	}

	if score.Valid {
		scoreValue := int(score.Int64)
		s.Score = &scoreValue
	}
	if feedback.Valid {
		s.Feedback = &feedback.String
	}
	if aiSummary.Valid {
		s.AISummary = &aiSummary.String
	}
	if gradedAt.Valid {
		s.GradedAt = &gradedAt.Time
	}
	if deletedAt.Valid {
		s.DeletedAt = &deletedAt.Time
	}

	return &s, nil
}

// FindByUserAndExercise تمام ارسال‌های یک کاربر برای یک تمرین را برمی‌گرداند.
func (r *SubmissionRepository) FindByUserAndExercise(ctx context.Context, userID, exerciseID string) ([]*submission.Submission, error) {
	query := `
        SELECT id, exercise_id, user_id, status, answer, score, feedback, ai_summary, ai_log_id,
               submitted_at, graded_at, created_at, updated_at, deleted_at
        FROM submissions
        WHERE user_id = $1 AND exercise_id = $2 AND deleted_at IS NULL
        ORDER BY submitted_at DESC
    `
	rows, err := r.db.QueryContext(ctx, query, userID, exerciseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []*submission.Submission
	for rows.Next() {
		var s submission.Submission
		var answerJSON []byte
		var score sql.NullInt64
		var feedback, aiSummary sql.NullString
		var gradedAt, deletedAt sql.NullTime
		var aiLogID sql.NullString

		if err := rows.Scan(
			&s.ID, &s.ExerciseID, &s.UserID, &s.Status, &answerJSON, &score, &feedback, &aiSummary, &aiLogID,
			&s.SubmittedAt, &gradedAt, &s.CreatedAt, &s.UpdatedAt, &deletedAt,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(answerJSON, &s.Answer); err != nil {
			return nil, fmt.Errorf("failed to unmarshal answer for submission %s: %w", s.ID, err)
		}

		if score.Valid {
			scoreValue := int(score.Int64)
			s.Score = &scoreValue
		}
		if feedback.Valid {
			s.Feedback = &feedback.String
		}
		if aiSummary.Valid {
			s.AISummary = &aiSummary.String
		}
		if gradedAt.Valid {
			s.GradedAt = &gradedAt.Time
		}
		if deletedAt.Valid {
			s.DeletedAt = &deletedAt.Time
		}
		
		submissions = append(submissions, &s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return submissions, nil
}


// UpdateAISummary خلاصه تحلیل AI و شناسه لاگ را آپدیت می‌کند.
func (r *SubmissionRepository) UpdateAISummary(ctx context.Context, submissionID, summary, aiLogID string) error {
	query := `
        UPDATE submissions SET ai_summary = $1, ai_log_id = $2, updated_at = NOW()
        WHERE id = $3
    `
	_, err := r.db.ExecContext(ctx, query, summary, aiLogID, submissionID)
	return err
}

// UpdateStatus وضعیت یک Submission را آپدیت می‌کند.
func (r *SubmissionRepository) UpdateStatus(ctx context.Context, id string, status submission.SubmissionStatus) error {
	query := `UPDATE submissions SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

// Update یک Submission موجود را آپدیت می‌کند (برای امتیازدهی).
func (r *SubmissionRepository) Update(ctx context.Context, sub *submission.Submission) error {
	answerJSON, err := json.Marshal(sub.Answer)
	if err != nil {
		return fmt.Errorf("failed to marshal submission answer: %w", err)
	}

	query := `
        UPDATE submissions
        SET
            status = $1,
            answer = $2,
            score = $3,
            feedback = $4,
            ai_summary = $5,
            ai_log_id = $6,
            graded_at = $7,
            updated_at = NOW()
        WHERE id = $8
    `
	_, err = r.db.ExecContext(ctx, query,
		sub.Status, answerJSON, sub.Score, sub.Feedback, sub.AISummary, sub.AILogID,
		sub.GradedAt, sub.ID,
	)
	return err
}

// GetExercise یک تمرین را بر اساس ID آن پیدا می‌کند.
func (r *SubmissionRepository) GetExercise(ctx context.Context, id string) (*submission.Exercise, error) {
	var ex submission.Exercise
	var contentJSON []byte
	var deletedAt sql.NullTime

	query := `
        SELECT id, chapter_id, type, difficulty, points, content, created_at, updated_at, deleted_at
        FROM exercises
        WHERE id = $1 AND deleted_at IS NULL
    `
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&ex.ID, &ex.ChapterID, &ex.Type, &ex.Difficulty, &ex.Points, &contentJSON,
		&ex.CreatedAt, &ex.UpdatedAt, &deletedAt,
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
	if deletedAt.Valid {
		ex.DeletedAt = &deletedAt.Time
	}

	return &ex, nil
}