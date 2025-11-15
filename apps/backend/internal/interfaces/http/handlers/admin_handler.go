// apps/backend/internal/interfaces/http/handlers/admin_handler.go

package handlers

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	app_admin "narrative-architecture/apps/backend/internal/application/admin"
	"narrative-architecture/apps/backend/internal/domain/user"
)

type AdminHandler struct {
	listUsersUC *app_admin.ListUsersUseCase
	updateStatusUC *app_admin.UpdateUserStatusUseCase
}

func NewAdminHandler(listUC *app_admin.ListUsersUseCase, updateUC *app_admin.UpdateUserStatusUseCase) *AdminHandler {
	return &AdminHandler{listUsersUC: listUC, updateStatusUC: updateUC}
}

// ListUsers لیست کاربران را برای پنل ادمین برمی‌گرداند.
func (h *AdminHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	
	users, total, err := h.listUsersUC.Execute(c.Request.Context(), page, pageSize)
	if err != nil {
		HandleError(c, err)
		return
	}
	
	// ... (تبدیل به DTO و ارسال پاسخ)
	c.JSON(http.StatusOK, gin.H{"data": users, "total": total})
}

// UpdateUserStatus وضعیت یک کاربر را تغییر می‌دهد.
func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	userID := c.Param("id")
	
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	
	useCaseReq := app_admin.UpdateUserStatusRequest{
		UserID: userID,
		Status: user.UserStatus(req.Status),
	}
	
	if err := h.updateStatusUC.Execute(c.Request.Context(), useCaseReq); err != nil {
		HandleError(c, err)
		return
	}
	
	c.Status(http.StatusOK)
}