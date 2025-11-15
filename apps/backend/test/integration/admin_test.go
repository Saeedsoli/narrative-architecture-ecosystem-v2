// apps/backend/test/integration/admin_test.go

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	
	app_admin "narrative-architecture/apps/backend/internal/application/admin"
	infra_postgres "narrative-architecture/apps/backend/internal/infrastructure/database/postgres"
	"narrative-architecture/apps/backend/internal/interfaces/http/handlers"
	"narrative-architecture/apps/backend/internal/interfaces/http/middleware"
)

func setupAdminRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// --- Dependency Injection ---
	userRepo := infra_postgres.NewUserRepository(pgDB)
	adminRepo := infra_postgres.NewAdminRepository(pgDB)

	listUsersUC := app_admin.NewListUsersUseCase(adminRepo)
	updateUserStatusUC := app_admin.NewUpdateUserStatusUseCase(userRepo)
	
	adminHandler := handlers.NewAdminHandler(listUsersUC, updateUserStatusUC)

	// --- Router Setup ---
	v1 := r.Group("/api/v1")
	{
		adminRoutes := v1.Group("/admin")
		// برای تست، فرض می‌کنیم کاربر ادمین است
		adminRoutes.Use(func(c *gin.Context) {
			c.Set("userRoles", []string{"admin"})
			c.Next()
		})
		{
			adminRoutes.GET("/users", adminHandler.ListUsers)
			adminRoutes.PUT("/users/:id/status", adminHandler.UpdateUserStatus)
		}
	}
	return r
}

func TestAdminIntegration(t *testing.T) {
	// --- Setup ---
	_, err := pgDB.Exec(`TRUNCATE TABLE users, user_profiles CASCADE;`)
	assert.NoError(t, err)
	// ایجاد یک کاربر تستی
	_, err = pgDB.Exec(`
		INSERT INTO users (id, email, password_hash) VALUES ('01HUSERIDTEST', 'test@user.com', 'hashed');
		INSERT INTO user_profiles (user_id, full_name, username) VALUES ('01HUSERIDTEST', 'Test User', 'testuser');
	`)
	assert.NoError(t, err)

	router := setupAdminRouter()

	// --- 1. تست لیست کردن کاربران ---
	req, _ := http.NewRequest("GET", "/api/v1/admin/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	
	var listResponse struct {
		Data []interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &listResponse)
	assert.Len(t, listResponse.Data, 1)

	// --- 2. تست تغییر وضعیت کاربر ---
	updatePayload := gin.H{"status": "suspended"}
	body, _ := json.Marshal(updatePayload)
	req, _ = http.NewRequest("PUT", "/api/v1/admin/users/01HUSERIDTEST/status", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// بررسی اینکه وضعیت در دیتابیس تغییر کرده است
	var status string
	err = pgDB.QueryRow("SELECT status FROM users WHERE id = '01HUSERIDTEST'").Scan(&status)
	assert.NoError(t, err)
	assert.Equal(t, "suspended", status)
}