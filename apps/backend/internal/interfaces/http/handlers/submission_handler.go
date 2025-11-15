// apps/backend/internal/interfaces/http/handlers/submission_handler.go

package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	app_submission "narrative-architecture/apps/backend/internal/application/submission"
	"narrative-architecture/apps/backend/internal/interfaces/http/dto"
)

// SubmissionHandler کنترلر HTTP برای تمام عملیات مربوط به ارسال تمرینات است.
type SubmissionHandler struct {
	submitExerciseUC *app_submission.SubmitExerciseUseCase
	analyzeUC        *app_submission.AnalyzeSubmissionUseCase
	getSubmissionsUC *app_submission.GetUserSubmissionsUseCase
	gradeSubmissionUC *app_submission.GradeSubmissionUseCase
}

// NewSubmissionHandler یک نمونه جدید از SubmissionHandler ایجاد می‌کند.
func NewSubmissionHandler(
	submitUC *app_submission.SubmitExerciseUseCase, 
	analyzeUC *app_submission.AnalyzeSubmissionUseCase,
	getUC *app_submission.GetUserSubmissionsUseCase,
	gradeUC *app_submission.GradeSubmissionUseCase,
) *SubmissionHandler {
	return &SubmissionHandler{
		submitExerciseUC: submitUC,
		analyzeUC:        analyzeUC,
		getSubmissionsUC: getUC,
		gradeSubmissionUC: gradeUC,
	}
}

// SubmitExercise یک پاسخ جدید برای یک تمرین ثبت می‌کند.
func (h *SubmissionHandler) SubmitExercise(c *gin.Context) {
	var req struct {
		ExerciseID string                 `json:"exerciseId" binding:"required"`
		Answer     map[string]interface{} `json:"answer" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	userID := c.GetString("userID")

	useCaseReq := app_submission.SubmitExerciseRequest{
		ExerciseID: req.ExerciseID,
		Answer:     req.Answer,
		UserID:     userID,
	}

	submission, err := h.submitExerciseUC.Execute(c.Request.Context(), useCaseReq)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"submissionId": submission.ID,
		"status":       submission.Status,
	})
}

// GetUserSubmissions لیست ارسال‌های یک کاربر برای یک تمرین را برمی‌گرداند.
func (h *SubmissionHandler) GetUserSubmissions(c *gin.Context) {
	userID := c.GetString("userID")
	exerciseID := c.Query("exerciseId")

	if exerciseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "exerciseId query parameter is required"})
		return
	}

	req := app_submission.GetUserSubmissionsRequest{
		UserID:     userID,
		ExerciseID: exerciseID,
	}

	submissions, err := h.getSubmissionsUC.Execute(c.Request.Context(), req)
	if err != nil {
		HandleError(c, err)
		return
	}

	var response []*dto.SubmissionResponse
	for _, sub := range submissions {
		response = append(response, dto.ToSubmissionResponse(sub))
	}
	
	c.JSON(http.StatusOK, response)
}

// AnalyzeSubmission فرآیند تحلیل AI برای یک پاسخ را آغاز می‌کند.
func (h *SubmissionHandler) AnalyzeSubmission(c *gin.Context) {
	submissionID := c.Param("id")
	userID := c.GetString("userID")

	err := h.analyzeUC.Execute(c.Request.Context(), submissionID, userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Analysis request accepted. You will be notified when it's complete.",
	})
}

// GradeSubmission به یک پاسخ امتیاز می‌دهد (فقط برای ادمین/مدرس).
func (h *SubmissionHandler) GradeSubmission(c *gin.Context) {
	submissionID := c.Param("id")
	graderID := c.GetString("userID")
	
	var req struct {
		Score    int    `json:"score" binding:"required,gte=0"`
		Feedback string `json:"feedback"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	useCaseReq := app_submission.GradeSubmissionRequest{
		SubmissionID: submissionID,
		GraderID:     graderID,
		Score:        req.Score,
		Feedback:     req.Feedback,
	}

	updatedSub, err := h.gradeSubmissionUC.Execute(c.Request.Context(), useCaseReq)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToSubmissionResponse(updatedSub))
}