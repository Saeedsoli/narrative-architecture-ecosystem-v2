// apps/backend/test/integration/submission_test.go

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/oklog/ulid/v2"

	app_submission "narrative-architecture/apps/backend/internal/application/submission"
	"narrative-architecture/apps/backend/internal/domain/user"
	infra_postgres "narrative-architecture/apps/backend/internal/infrastructure/database/postgres"
	"narrative-architecture/apps/backend/internal/interfaces/http/handlers"
	"narrative-architecture/apps/backend/internal/mocks"
)

// MockAIClient یک نسخه Mock از کلاینت AI برای استفاده در تست است.
// این ساختار را می‌توان به یک فایل مشترک در mocks منتقل کرد.
type MockAIClient struct {
	mock.Mock
}

func (m *MockAIClient) AnalyzeText(ctx context.Context, text, context string) (*app_submission.AIAnalysisResult, error) {
	args := m.Called(ctx, text, context)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*app_submission.AIAnalysisResult), args.Error(1)
}

// setupSubmissionRouter یک روتر Gin با وابستگی‌های لازم برای ماژول Submission ایجاد می‌کند.
func setupSubmissionRouter(aiClient *MockAIClient) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// --- Dependency Injection ---
	userRepo := infra_postgres.NewUserRepository(pgDB) // از pgDB که در main_test.go تعریف شده
	submissionRepo := infra_postgres.NewSubmissionRepository(pgDB)

	analyzeUC := app_submission.NewAnalyzeSubmissionUseCase(submissionRepo, aiClient)
	submitUC := app_submission.NewSubmitExerciseUseCase(submissionRepo, userRepo, analyzeUC)
	getUserSubmissionsUC := app_submission.NewGetUserSubmissionsUseCase(submissionRepo)
	gradeSubmissionUC := app_submission.NewGradeSubmissionUseCase(submissionRepo, userRepo)

	submissionHandler := handlers.NewSubmissionHandler(
		submitUC,
		analyzeUC,
		getUserSubmissionsUC,
		gradeSubmissionUC,
	)

	// --- Router Setup ---
	v1 := r.Group("/api/v1")
	{
		protected := v1.Group("")
		protected.Use(func(c *gin.Context) {
			c.Set("userID", "01HUSERIDTEST")
			c.Next()
		})
		{
			submissions := protected.Group("/submissions")
			{
				submissions.POST("", submissionHandler.SubmitExercise)
				submissions.POST("/:id/analyze", submissionHandler.AnalyzeSubmission)
				// TODO: Add GET and grading routes
			}
		}
	}
	return r
}

// TestSubmissionIntegration تست یکپارچه برای جریان کامل Submission است.
func TestSubmissionIntegration(t *testing.T) {
	// --- Setup ---
	// تمیز کردن جداول و ایجاد داده‌های اولیه
	_, err := pgDB.Exec(`
		TRUNCATE TABLE submissions, exercises, ai_logs, users, user_profiles CASCADE;
		
		-- ایجاد یک کاربر تستی
		INSERT INTO users (id, email, password_hash) VALUES ('01HUSERIDTEST', 'test@user.com', 'hashed_password');
		INSERT INTO user_profiles (user_id, full_name, username) VALUES ('01HUSERIDTEST', 'Test User', 'testuser');

		-- ایجاد یک تمرین تستی
		INSERT INTO exercises (id, chapter_id, type, difficulty, points, content)
		VALUES ('01HEXERCISEID', 'chapter-01', 'essay', 'intermediate', 10, '{}');
	`)
	assert.NoError(t, err)

	mockAI := new(MockAIClient)
	router := setupSubmissionRouter(mockAI)
	
	// --- 1. تست ارسال موفق یک پاسخ ---
	submissionPayload := gin.H{
		"exerciseId": "01HEXERCISEID",
		"answer": gin.H{
			"text": "این متن پاسخ من به تمرین مقاله نویسی است.",
		},
	}
	body, _ := json.Marshal(submissionPayload)
	req, _ := http.NewRequest("POST", "/api/v1/submissions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var submissionResponse struct {
		SubmissionID string `json:"submissionId"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &submissionResponse)
	assert.NoError(t, err)
	submissionID := submissionResponse.SubmissionID
	assert.NotEmpty(t, submissionID)

	// --- 2. تست درخواست تحلیل AI ---
	// تنظیم Mock برای پاسخ سرویس AI
	mockAI.On("AnalyzeText", mock.Anything, "این متن پاسخ من به تمرین مقاله نویسی است.", mock.Anything).Return(&app_submission.AIAnalysisResult{
		Analysis: "تحلیل شما عالی بود! ساختار متن شما بسیار منسجم است.",
		FullResponse: map[string]interface{}{"model": "gpt-4-turbo", "usage": 120},
	}, nil).Once()

	req, _ = http.NewRequest("POST", "/api/v1/submissions/"+submissionID+"/analyze", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusAccepted, w.Code)
	
	// --- 3. بررسی اینکه تحلیل در دیتابیس ذخیره شده است ---
	// چون تحلیل در یک goroutine انجام می‌شود، باید کمی صبر کنیم.
	time.Sleep(100 * time.Millisecond)
	
	var aiSummary sql.NullString
	var aiLogID sql.NullString
	err = pgDB.QueryRow("SELECT ai_summary, ai_log_id FROM submissions WHERE id = $1", submissionID).Scan(&aiSummary, &aiLogID)
	assert.NoError(t, err)
	assert.True(t, aiSummary.Valid)
	assert.Equal(t, "تحلیل شما عالی بود! ساختار متن شما بسیار منسجم است.", aiSummary.String)
	assert.True(t, aiLogID.Valid)
	
	// اطمینان از اینکه Mock AI Client فراخوانی شده است
	mockAI.AssertExpectations(t)
}