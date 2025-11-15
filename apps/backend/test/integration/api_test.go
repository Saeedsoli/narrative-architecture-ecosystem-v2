// apps/backend/test/integration/api_test.go

package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"narrative-architecture/apps/backend/internal/interfaces/http/handlers"
)

// setupTestRouter یک روتر ساده برای تست‌های عمومی ایجاد می‌کند.
func setupTestRouter() *gin.Engine {
	r := gin.New()
	healthHandler := handlers.NewHealthHandler()
	r.GET("/health", healthHandler.CheckHealth)
	return r
}

func TestHealthCheck(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"status":"ok"`)
}