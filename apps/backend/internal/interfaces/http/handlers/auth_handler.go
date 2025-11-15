// apps/backend/internal/interfaces/http/handlers/auth_handler.go

package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"narrative-architecture/apps/backend/internal/application/auth"
	"narrative-architecture/apps/backend/internal/interfaces/http/dto"
)

type AuthHandler struct {
	registerUC *auth.RegisterUserUseCase
	loginUC    *auth.LoginUserUseCase
	logoutUC   *auth.LogoutUserUseCase
	refreshUC  *auth.RefreshTokenUseCase
}

func NewAuthHandler(registerUC *auth.RegisterUserUseCase, loginUC *auth.LoginUserUseCase, logoutUC *auth.LogoutUserUseCase, refreshUC *auth.RefreshTokenUseCase) *AuthHandler {
	return &AuthHandler{
		registerUC: registerUC,
		loginUC:    loginUC,
		logoutUC:   logoutUC,
		refreshUC:  refreshUC,
	}
}

// Register یک کاربر جدید ثبت‌نام می‌کند.
func (h *AuthHandler) Register(c *gin.Context) {
	var req auth.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	res, err := h.registerUC.Execute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("refreshToken", res.RefreshToken, int(30*24*time.Hour.Seconds()), "/", "localhost", true, true)

	c.JSON(http.StatusCreated, gin.H{
		"user":        dto.ToUserResponse(res.User),
		"accessToken": res.AccessToken,
	})
}

// Login کاربر را وارد سیستم می‌کند.
func (h *AuthHandler) Login(c *gin.Context) {
	var req auth.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	res, err := h.loginUC.Execute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("refreshToken", res.RefreshToken, int(30*24*time.Hour.Seconds()), "/", "localhost", true, true)

	c.JSON(http.StatusOK, gin.H{
		"user":        dto.ToUserResponse(res.User),
		"accessToken": res.AccessToken,
	})
}

// Refresh یک Access Token جدید با استفاده از Refresh Token تولید می‌کند.
func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not found in cookie"})
		return
	}

	res, err := h.refreshUC.Execute(c.Request.Context(), refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("refreshToken", res.RefreshToken, int(30*24*time.Hour.Seconds()), "/", "localhost", true, true)

	c.JSON(http.StatusOK, gin.H{
		"accessToken": res.AccessToken,
	})
}

// Logout کاربر را از سیستم خارج می‌کند.
func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, _ := c.Cookie("refreshToken")
	if refreshToken != "" {
		// ما اهمیتی به خطای اینجا نمی‌دهیم، چون در هر صورت کوکی را پاک می‌کنیم
		_ = h.logoutUC.Execute(c.Request.Context(), refreshToken)
	}

	// کوکی را در مرورگر پاک می‌کند
	c.SetCookie("refreshToken", "", -1, "/", "localhost", true, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}