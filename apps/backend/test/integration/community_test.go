// apps/backend/test/integration/community_test.go

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
	
	app_community "narrative-architecture/apps/backend/internal/application/community"
	infra_mongo "narrative-architecture/apps/backend/internal/infrastructure/database/mongodb"
	"narrative-architecture/apps/backend/internal/interfaces/http/handlers"
)

func setupCommunityRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// --- Dependency Injection ---
	communityRepo := infra_mongo.NewCommunityRepository(mongoDB)
	createTopicUC := app_community.NewCreateTopicUseCase(communityRepo)
	listTopicsUC := app_community.NewListTopicsUseCase(communityRepo)
	getTopicUC := app_community.NewGetTopicUseCase(communityRepo)
	createPostUC := app_community.NewCreatePostUseCase(communityRepo)
	listPostsUC := app_community.NewListPostsUseCase(communityRepo)
	
	communityHandler := handlers.NewCommunityHandler(createTopicUC, listTopicsUC, getTopicUC, createPostUC, listPostsUC)
	
	// --- Router Setup ---
	v1 := r.Group("/api/v1")
	{
		community := v1.Group("/community")
		{
			community.GET("/topics", communityHandler.ListTopics)
			
			protected := community.Group("")
			protected.Use(func(c *gin.Context) {
				c.Set("userID", "01HUSERIDTEST")
				c.Set("username", "testuser")
				c.Next()
			})
			{
				protected.POST("/topics", communityHandler.CreateTopic)
			}
		}
	}
	return r
}

func TestCommunityIntegration(t *testing.T) {
	// --- Setup ---
	mongoDB.Collection("forum_topics").Drop(context.Background())
	mongoDB.Collection("forum_posts").Drop(context.Background())

	router := setupCommunityRouter()

	// --- 1. تست ایجاد تاپیک جدید ---
	topicPayload := gin.H{
		"locale": "fa",
		"title":  "تاپیک تستی",
		"body":   "این اولین پست تاپیک است.",
		"tags":   []string{"تست"},
	}
	body, _ := json.Marshal(topicPayload)
	req, _ := http.NewRequest("POST", "/api/v1/community/topics", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var topicResponse struct {
		ID string `json:"ID"`
	}
	json.Unmarshal(w.Body.Bytes(), &topicResponse)
	assert.NotEmpty(t, topicResponse.ID)
	
	// --- 2. تست لیست کردن تاپیک‌ها ---
	req, _ = http.NewRequest("GET", "/api/v1/community/topics?locale=fa", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	
	var listResponse struct {
		Data []interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &listResponse)
	assert.Len(t, listResponse.Data, 1)
}