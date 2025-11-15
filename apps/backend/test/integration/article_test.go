// apps/backend/test/integration/article_test.go

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	
	app_article "narrative-architecture/apps/backend/internal/application/article"
	domain_article "narrative-architecture/apps/backend/internal/domain/article"
	infra_cache "narrative-architecture/apps/backend/internal/infrastructure/cache/redis"
	infra_mongo "narrative-architecture/apps/backend/internal/infrastructure/database/mongodb"
	infra_postgres "narrative-architecture/apps/backend/internal/infrastructure/database/postgres"
	"narrative-architecture/apps/backend/internal/interfaces/http/handlers"
	"narrative-architecture/apps/backend/internal/interfaces/http/middleware"
)

// setupArticleRouter یک روتر Gin با تمام وابستگی‌های لازم برای ماژول Article ایجاد می‌کند.
func setupArticleRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// --- Dependency Injection ---
	// Redis client for cache (can be mocked or connected to a real test instance)
	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	articleCache := infra_cache.NewArticleCache(redisClient)

	// Repositories
	pgArticleRepo := infra_postgres.NewArticleRepository(pgDB)
	mongoArticleRepo := infra_mongo.NewArticleRepository(mongoDB)
	
	// Domain Service
	articleDomainService := domain_article.NewService()

	// Use Cases
	getArticleUC := app_article.NewGetArticleUseCase(mongoArticleRepo, articleCache)
	listArticlesUC := app_article.NewListArticlesUseCase(mongoArticleRepo)
	createArticleUC := app_article.NewCreateArticleUseCase(mongoArticleRepo, nil) // userRepo is nil for now
	updateArticleUC := app_article.NewUpdateArticleUseCase(mongoArticleRepo, articleCache)
	deleteArticleUC := app_article.NewDeleteArticleUseCase(mongoArticleRepo, articleCache)
	publishArticleUC := app_article.NewPublishArticleUseCase(mongoArticleRepo, articleDomainService, articleCache)
	addBookmarkUC := app_article.NewAddBookmarkUseCase(pgArticleRepo)

	// Handler
	articleHandler := handlers.NewArticleHandler(
		getArticleUC, createArticleUC, updateArticleUC,
		deleteArticleUC, listArticlesUC, publishArticleUC, addBookmarkUC,
	)

	// --- Router Setup ---
	v1 := r.Group("/api/v1")
	{
		// Public routes
		v1.GET("/articles", articleHandler.ListArticles)
		v1.GET("/articles/:slug", articleHandler.GetArticle)
		
		// Protected routes
		protected := v1.Group("")
		// For testing, we use a mock middleware that sets a user ID.
		protected.Use(func(c *gin.Context) {
			c.Set("userID", "01HUSERIDTEST")
			c.Next()
		})
		{
			protectedArticles := protected.Group("/articles")
			{
				protectedArticles.POST("", articleHandler.CreateArticle)
				protectedArticles.PUT("/:id", articleHandler.UpdateArticle)
				protectedArticles.DELETE("/:id", articleHandler.DeleteArticle)
				protectedArticles.POST("/:id/publish", articleHandler.PublishArticle)
				protectedArticles.POST("/:id/bookmark", articleHandler.AddBookmark)
			}
		}
	}
	return r
}

// TestArticleCRUD یک تست یکپارچه برای جریان کامل CRUD مقالات است.
func TestArticleCRUD(t *testing.T) {
	// --- Setup ---
	// تمیز کردن کالکشن‌ها قبل از هر تست
	err := mongoDB.Collection("articles").Drop(context.Background())
	assert.NoError(t, err)
	_, err = pgDB.Exec("TRUNCATE TABLE article_bookmarks")
	assert.NoError(t, err)

	router := setupArticleRouter()
	
	// --- 1. Create Article ---
	createPayload := gin.H{
		"locale":  "fa",
		"title":   "یک مقاله تستی برای جریان CRUD",
		"content": "این متن محتوای مقاله است.",
		"tags":    []string{"تست", "یکپارچه‌سازی"},
	}
	body, _ := json.Marshal(createPayload)
	req, _ := http.NewRequest("POST", "/api/v1/articles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var createArticleResponse struct {
		ID   string `json:"id"`
		Slug string `json:"slug"`
	}
	json.Unmarshal(w.Body.Bytes(), &createArticleResponse)
	assert.NotEmpty(t, createArticleResponse.ID)
	assert.NotEmpty(t, createArticleResponse.Slug)
	
	articleID := createArticleResponse.ID
	articleSlug := createArticleResponse.Slug

	// --- 2. Get Article (Should Fail - Not Published Yet) ---
	req, _ = http.NewRequest("GET", "/api/v1/articles/"+articleSlug, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// --- 3. Publish Article ---
	req, _ = http.NewRequest("POST", "/api/v1/articles/"+articleID+"/publish", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	
	// --- 4. Get Article (Should Succeed Now) ---
	req, _ = http.NewRequest("GET", "/api/v1/articles/"+articleSlug, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	
	var getArticleResponse struct {
		Title string `json:"title"`
	}
	json.Unmarshal(w.Body.Bytes(), &getArticleResponse)
	assert.Equal(t, "یک مقاله تستی برای جریان CRUD", getArticleResponse.Title)

	// --- 5. Update Article ---
	updatePayload := gin.H{
		"title": "عنوان آپدیت شده مقاله",
	}
	body, _ = json.Marshal(updatePayload)
	req, _ = http.NewRequest("PUT", "/api/v1/articles/"+articleID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	
	// --- 6. Add Bookmark ---
	req, _ = http.NewRequest("POST", "/api/v1/articles/"+articleID+"/bookmark", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	
	// --- 7. Delete Article (Soft Delete) ---
	req, _ = http.NewRequest("DELETE", "/api/v1/articles/"+articleID, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
	
	// --- 8. Get Article (Should Fail - Deleted) ---
	req, _ = http.NewRequest("GET", "/api/v1/articles/"+articleSlug, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}