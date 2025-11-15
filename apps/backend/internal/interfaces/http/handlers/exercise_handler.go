// apps/backend/internal/interfaces/http/handlers/exercise_handler.go

package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	app_exercise "narrative-architecture/apps/backend/internal/application/exercise"
)

type ExerciseHandler struct {
	getExerciseUC *app_exercise.GetExerciseUseCase
	listExercisesUC *app_exercise.ListExercisesUseCase
}

func NewExerciseHandler(getUC *app_exercise.GetExerciseUseCase, listUC *app_exercise.ListExercisesUseCase) *ExerciseHandler {
	return &ExerciseHandler{getExerciseUC: getUC, listExercisesUC: listUC}
}

// GetExercise جزئیات یک تمرین را برمی‌گرداند.
func (h *ExerciseHandler) GetExercise(c *gin.Context) {
	exerciseID := c.Param("id")
	exercise, err := h.getExerciseUC.Execute(c.Request.Context(), exerciseID)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, exercise)
}

// ListExercises لیست تمرینات یک فصل را برمی‌گرداند.
func (h *ExerciseHandler) ListExercises(c *gin.Context) {
	chapterID := c.Query("chapterId")
	if chapterID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "chapterId query parameter is required"})
		return
	}
	
	exercises, err := h.listExercisesUC.Execute(c.Request.Context(), chapterID)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, exercises)
}