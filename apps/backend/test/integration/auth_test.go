// apps/backend/test/integration/auth_test.go

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"narrative-architecture/apps/backend/internal/application/auth"
	"narrative-architecture/apps/backend/internal/infrastructure/database/postgres"
	"narrative-architecture/apps/backend/internal/interfaces/http/handlers"
	"narrative-architecture/apps/backend/internal/interfaces/http/middleware"
)

func setupAuthRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// --- Dependency Injection for Auth ---
	userRepo := postgres.NewUserRepository(db)
	tokenRepo := postgres.NewTokenRepository(db)

	accessTTL := 15 * time.Minute
	refreshTTL := 24 * time.Hour
	accessSecret := "test_access_secret"
	refreshSecret := "test_refresh_secret"

	registerUC := auth.NewRegisterUserUseCase(userRepo, tokenRepo, accessSecret, refreshSecret, accessTTL, refreshTTL)
	loginUC := auth.NewLoginUserUseCase(userRepo, tokenRepo, accessSecret, refreshSecret, accessTTL, refreshTTL)
	logoutUC := auth.NewLogoutUserUseCase(tokenRepo)
	refreshUC := auth.NewRefreshTokenUseCase(userRepo, tokenRepo, accessSecret, refreshSecret, accessTTL, refreshTTL)
	
	authHandler := handlers.NewAuthHandler(registerUC, loginUC, logoutUC, refreshUC)
	
	// --- Router Setup ---
	v1 := r.Group("/api/v1")
	{
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/login", authHandler.Login)
		}
		
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(accessSecret))
		{
			protected.POST("/auth/logout", authHandler.Logout)
		}
	}
	return r
}

func TestAuthIntegration(t *testing.T) {
	// هر تست باید با یک دیتابیس تمیز شروع شود
	_, err := db.Exec("TRUNCATE TABLE users, user_profiles, user_roles, auth_tokens, roles CASCADE")
	assert.NoError(t, err)
	// Seed roles
	_, err = db.Exec(`INSERT INTO roles (id, name) VALUES ('01H00000000000000000000005', 'user') ON CONFLICT DO NOTHING`)
	assert.NoError(t, err)

	router := setupAuthRouter()

	// --- 1. تست ثبت‌نام موفق ---
	registerPayload := gin.H{
		"email":    "integration@test.com",
		"password": "Password123!",
		"username": "integ_test",
		"fullName": "Integration Test",
	}
	body, _ := json.Marshal(registerPayload)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var registerResponse struct {
		User struct {
			Email string `json:"email"`
		} `json:"user"`
		AccessToken string `json:"accessToken"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &registerResponse)
	assert.NoError(t, err)
	assert.Equal(t, "integration@test.com", registerResponse.User.Email)
	assert.NotEmpty(t, registerResponse.AccessToken)
	
	// بررسی کوکی Refresh Token
	refreshTokenCookie := w.Result().Cookies()[0]
	assert.Equal(t, "refreshToken", refreshTokenCookie.Name)
	assert.NotEmpty(t, refreshTokenCookie.Value)
	assert.True(t, refreshTokenCookie.HttpOnly)

	// --- 2. تست ورود موفق ---
	loginPayload := gin.H{
		"email":    "integration@test.com",
		"password": "Password123!",
	}
	body, _ = json.Marshal(loginPayload)
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	// --- 3. تست ورود با رمز عبور اشتباه ---
	loginPayload["password"] = "wrongpassword"
	body, _ = json.Marshal(loginPayload)
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}