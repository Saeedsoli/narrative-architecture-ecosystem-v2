// apps/backend/internal/interfaces/router.go

package interfaces

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	
	"narrative-architecture/apps/backend/internal/interfaces/http/handlers"
	"narrative-architecture/apps/backend/internal/interfaces/http/middleware"
	"narrative-architecture/apps/backend/pkg/config"
)

// SetupRouter روتر اصلی Gin را با تمام Middlewareها و روت‌ها پیکربندی می‌کند.
func SetupRouter(
	cfg *config.Config,
	redisClient *redis.Client,
	authHandler *handlers.AuthHandler,
	articleHandler *handlers.ArticleHandler,
	submissionHandler *handlers.SubmissionHandler,
	communityHandler *handlers.CommunityHandler,
	adminHandler *handlers.AdminHandler,
	contentHandler *handlers.ContentHandler,
	healthHandler *handlers.HealthHandler,
) *gin.Engine {
	r := gin.New()

	// --- Global Middlewares ---
	r.Use(gin.Recovery()) // بازیابی از panicها
	r.Use(middleware.LoggingMiddleware()) // لاگ‌گیری ساختاریافته
	r.Use(middleware.CorsMiddleware([]string{"http://localhost:3000"})) // تنظیم CORS
	r.Use(middleware.RateLimiterMiddleware(redisClient, 100, 1*time.Minute)) // محدودسازی درخواست‌ها

	// --- Routes ---
	r.GET("/health", healthHandler.CheckHealth)
	
	v1 := r.Group("/api/v1")
	{
		// --- Public Routes ---
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/login", authHandler.Login)
			authRoutes.POST("/refresh", authHandler.Refresh)
		}

		publicArticleRoutes := v1.Group("/articles")
		{
			publicArticleRoutes.GET("", articleHandler.ListArticles)
			publicArticleRoutes.GET("/:slug", articleHandler.GetArticle)
		}
		
		publicCommunityRoutes := v1.Group("/community")
		{
			publicCommunityRoutes.GET("/topics", communityHandler.ListTopics)
			publicCommunityRoutes.GET("/topics/:id", communityHandler.GetTopic)
			publicCommunityRoutes.GET("/topics/:id/posts", communityHandler.ListPosts)
		}

		// --- Protected Routes ---
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWT.AccessSecret))
		{
			protected.GET("/auth/me", authHandler.Me)
			protected.POST("/auth/logout", authHandler.Logout)

			protectedArticleRoutes := protected.Group("/articles")
			{
				protectedArticleRoutes.POST("", articleHandler.CreateArticle)
				protectedArticleRoutes.PUT("/:id", articleHandler.UpdateArticle)
				protectedArticleRoutes.DELETE("/:id", articleHandler.DeleteArticle)
				protectedArticleRoutes.POST("/:id/publish", articleHandler.PublishArticle)
				protectedArticleRoutes.POST("/:id/bookmark", articleHandler.AddBookmark)
				protectedArticleRoutes.DELETE("/:id/bookmark", articleHandler.RemoveBookmark)
			}
			
			protected.GET("/bookmarks", articleHandler.ListUserBookmarks)

			protectedCommunityRoutes := protected.Group("/community")
			{
				protectedCommunityRoutes.POST("/topics", communityHandler.CreateTopic)
				protectedCommunityRoutes.POST("/topics/:id/posts", communityHandler.CreatePost)
				protectedCommunityRoutes.POST("/votes", communityHandler.AddVote)
			}

			submissions := protected.Group("/submissions")
			{
				submissions.POST("", submissionHandler.SubmitExercise)
				submissions.GET("", submissionHandler.GetUserSubmissions)
				submissions.POST("/:id/analyze", submissionHandler.AnalyzeSubmission)
			}
			
			protected.GET("/content/chapter-16", contentHandler.GetChapter16)
			
			// --- Admin Routes ---
			adminRoutes := protected.Group("/admin")
			adminRoutes.Use(middleware.RoleMiddleware("admin", "moderator"))
			{
				adminRoutes.GET("/users", adminHandler.ListUsers)
				adminRoutes.PUT("/users/:id/status", adminHandler.UpdateUserStatus)
				adminRoutes.GET("/moderation", adminHandler.ListModerationQueue)
				adminRoutes.POST("/moderation/:id", adminHandler.ModerateContent)
			}
		}
	}

	return r
}