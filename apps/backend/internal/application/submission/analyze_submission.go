// apps/backend/internal/application/submission/analyze_submission.go

package submission

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"narrative-architecture/apps/backend/internal/domain/submission"
	"github.com/oklog/ulid/v2"
)

// AIClient رابطی برای ارتباط با سرویس AI است.
type AIClient interface {
	AnalyzeText(ctx context.Context, text, context string) (*AIAnalysisResult, error)
}

// AIAnalysisResult ساختار پاسخی است که از AIClient دریافت می‌شود.
type AIAnalysisResult struct {
	Analysis      string
	FullResponse  map[string]interface{} // پاسخ کامل برای لاگ کردن
}

// AnalyzeSubmissionUseCase منطق تحلیل یک پاسخ تمرین را مدیریت می‌کند.
type AnalyzeSubmissionUseCase struct {
	submissionRepo submission.Repository
	aiClient       AIClient
}

// NewAnalyzeSubmissionUseCase یک نمونه جدید از Use Case ایجاد می‌کند.
func NewAnalyzeSubmissionUseCase(repo submission.Repository, client AIClient) *AnalyzeSubmissionUseCase {
	return &AnalyzeSubmissionUseCase{
		submissionRepo: repo,
		aiClient:       client,
	}
}

// Execute متد اصلی برای اجرای Use Case است.
func (uc *AnalyzeSubmissionUseCase) Execute(ctx context.Context, submissionID string, userID string) error {
	// 1. دریافت Submission از دیتابیس
	sub, err := uc.submissionRepo.FindByID(ctx, submissionID)
	if err != nil {
		return fmt.Errorf("submission not found: %w", err)
	}

	// 2. بررسی مالکیت
	if sub.UserID != userID {
		return errors.New("user does not have permission to analyze this submission")
	}
	
	// 3. بررسی اینکه آیا قبلاً تحلیل شده است
	if sub.AISummary != "" {
		log.Printf("INFO: Submission %s already analyzed. Skipping.", submissionID)
		return nil // خطایی وجود ندارد، صرفاً عملیات انجام نمی‌شود
	}
	
	// 4. استخراج متن پاسخ برای ارسال به AI
	textToAnalyze, ok := sub.Answer["text"].(string)
	if !ok || textToAnalyze == "" {
		return errors.New("submission contains no text to analyze")
	}

	// 5. فراخوانی سرویس AI در یک goroutine جداگانه
	go func() {
		// یک context جدید با timeout برای این عملیات پس‌زمینه ایجاد می‌کنیم
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		
		log.Printf("INFO: Starting AI analysis for submission %s", sub.ID)
		
		// فراخوانی کلاینت AI
		analysisResult, err := uc.aiClient.AnalyzeText(bgCtx, textToAnalyze, fmt.Sprintf("exercise:%s", sub.ExerciseID))
		if err != nil {
			log.Printf("ERROR: AI analysis failed for submission %s: %v", sub.ID, err)
			// TODO: وضعیت submission را به 'analysis_failed' تغییر دهید
			return
		}
		
		// 6. ذخیره لاگ کامل AI (این متد باید در Repository پیاده‌سازی شود)
		aiLogID := ulid.New().String()
		// err = uc.submissionRepo.SaveAILog(bgCtx, aiLogID, sub.ID, analysisResult.FullResponse)
		// if err != nil {
		// 	log.Printf("ERROR: Failed to save AI log for submission %s: %v", sub.ID, err)
		// }

		// 7. ذخیره خلاصه تحلیل در Submission
		err = uc.submissionRepo.UpdateAISummary(bgCtx, sub.ID, analysisResult.Analysis, aiLogID)
		if err != nil {
			log.Printf("ERROR: Failed to save AI analysis summary for submission %s: %v", sub.ID, err)
			return
		}
		
		log.Printf("SUCCESS: AI analysis completed for submission %s", sub.ID)
		
		// TODO: ارسال یک نوتیفیکیشن به کاربر (با WebSocket یا ایمیل)
		// eventBus.Publish("submission.analyzed", { submissionId: sub.ID, userId: sub.UserID })
	}()

	return nil
}