// apps/backend/cmd/api/main.go

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	
	// لایه‌های اپلیکیشن
	app_admin "narrative-architecture/apps/backend/internal/application/admin"
	app_auth "narrative-architecture/apps/backend/internal/application/auth"
	app_article "narrative-architecture/apps/backend/internal/application/article"
	app_community "narrative-architecture/apps/backend/internal/application/community"
	app_content "narrative-architecture/apps/backend/internal/application/content"
	app_submission "narrative-architecture/apps/backend/internal/application/submission"
	
	// لایه دامنه
	domain_article "narrative-architecture/apps/backend/internal/domain/article"

	// لایه زیرساخت
	infra_ai "narrative-architecture/apps/backend/internal/infrastructure/ai"
	infra_cache "narrative-architecture/apps/backend/internal/infrastructure/cache"
	infra_mongo "narrative-architecture/apps/backend/internal/infrastructure/database/mongodb"
	infra_postgres "narrative-architecture/apps/backend/internal/infrastructure/database/postgres"
	infra_storage "narrative-architecture/apps/backend/internal/infrastructure/storage"

	// لایه رابط
	"narrative-architecture/apps/backend/internal/interfaces"
	"narrative-architecture/apps/backend/internal/interfaces/http/handlers"
	custom_validator "narrative-architecture/apps/backend/internal/interfaces/http/validator"
	
	// پکیج‌های کمکی
	pkg_config "narrative-architecture/apps/backend/pkg/config"
	pkg_logger "narrative-architecture/apps/backend/pkg/logger"
)

func main() {
	// --- 1. بارگذاری پیکربندی و راه‌اندازی لاگر ---
	cfg := pkg_config.LoadConfig()
	pkg_logger.InitLogger(cfg.Server.GinMode)

	gin.SetMode(cfg.Server.GinMode)
	log.Println("INFO: Starting Narrative Architecture Backend...")

	// --- 2. اتصال به دیتابیس‌ها و سرویس‌های خارجی ---
	log.Println("INFO: Connecting to databases and external services...")

	// PostgreSQL
	pgDB := infra_postgres.ConnectDB(cfg.DB)
	defer pgDB.Close()

	// MongoDB
	mongoClient, mongoDB := infra_mongo.ConnectMongo(cfg.Mongo)
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("ERROR: Could not disconnect from MongoDB: %v", err)
		}
	}()

	// Redis
	redisClient := infra_cache.ConnectRedis(cfg.Redis)
	
	// S3 Client
	s3Client := infra_storage.NewS3Client(infra_storage.S3Config{Region: cfg.Services.S3Region})

	log.Println("INFO: Database connections established successfully.")

	// --- 3. Dependency Injection (ایجاد نمونه از تمام لایه‌ها) ---
	log.Println("INFO: Initializing application layers...")

	// Repositories
	userRepo := infra_postgres.NewUserRepository(pgDB)
	tokenRepo := infra_postgres.NewTokenRepository(pgDB)
	adminRepo := infra_postgres.NewAdminRepository(pgDB)
	submissionRepo := infra_postgres.NewSubmissionRepository(pgDB)
	pgArticleRepo := infra_postgres.NewArticleRepository(pgDB)
	mongoArticleRepo := infra_mongo.NewArticleRepository(mongoDB)
	communityRepo := infra_mongo.NewForumRepository(mongoDB) // استفاده از نام دقیق‌تر
	moderationRepo := infra_postgres.NewModerationRepository(pgDB)
	contentRepo := infra_mongo.NewContentRepository(mongoDB)
	entitlementRepo := infra_postgres.NewEntitlementRepository(pgDB)

	// External Services & Caches
	articleCache := infra_cache.NewArticleCache(redisClient)
	aiClient := infra_ai.NewClient(cfg.Services.AI)
	contentStorage := infra_storage.NewContentStorage(s3Client, cfg.Services.S3Bucket)

	// Domain Services
	articleDomainService := domain_article.NewService()
	
	// Use Cases
	registerUC := app_auth.NewRegisterUserUseCase(userRepo, tokenRepo, cfg.JWT.AccessSecret, cfg.JWT.RefreshSecret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)
	loginUC := app_auth.NewLoginUserUseCase(userRepo, tokenRepo, cfg.JWT.AccessSecret, cfg.JWT.RefreshSecret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)
	logoutUC := app_auth.NewLogoutUserUseCase(tokenRepo)
	refreshUC := app_auth.NewRefreshTokenUseCase(userRepo, tokenRepo, cfg.JWT.AccessSecret, cfg.JWT.RefreshSecret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)

	getArticleUC := app_article.NewGetArticleUseCase(mongoArticleRepo, articleCache)
	listArticlesUC := app_article.NewListArticlesUseCase(mongoArticleRepo)
	createArticleUC := app_article.NewCreateArticleUseCase(mongoArticleRepo, userRepo)
	updateArticleUC := app_article.NewUpdateArticleUseCase(mongoArticleRepo, articleCache)
	deleteArticleUC := app_article.NewDeleteArticleUseCase(mongoArticleRepo, articleCache)
	publishArticleUC := app_article.NewPublishArticleUseCase(mongoArticleRepo, articleDomainService, articleCache)
	addBookmarkUC := app_article.NewAddBookmarkUseCase(pgArticleRepo)
	removeBookmarkUC := app_article.NewRemoveBookmarkUseCase(pgArticleRepo)
	listBookmarksUC := app_article.NewListUserBookmarksUseCase(pgArticleRepo, mongoArticleRepo)

	createTopicUC := app_community.NewCreateTopicUseCase(communityRepo)
	listTopicsUC := app_community.NewListTopicsUseCase(communityRepo)
	getTopicUC := app_community.NewGetTopicUseCase(communityRepo)
	createPostUC := app_community.NewCreatePostUseCase(communityRepo)
	listPostsUC := app_community.NewListPostsUseCase(communityRepo)
	addVoteUC := app_community.NewAddVoteUseCase(communityRepo)

	listUsersUC := app_admin.NewListUsersUseCase(adminRepo)
	updateUserStatusUC := app_admin.NewUpdateUserStatusUseCase(userRepo)
	listModerationUC := app_admin.NewListModerationQueueUseCase(moderationRepo)
	moderateContentUC := app_admin.NewModerateContentUseCase(moderationRepo, contentRepo)

	analyzeSubmissionUC := app_submission.NewAnalyzeSubmissionUseCase(submissionRepo, aiClient)
	submitExerciseUC := app_submission.NewSubmitExerciseUseCase(submissionRepo, userRepo, analyzeSubmissionUC)
	getUserSubmissionsUC := app_submission.NewGetUserSubmissionsUseCase(submissionRepo)
	gradeSubmissionUC := app_submission.NewGradeSubmissionUseCase(submissionRepo, userRepo)
	
	getChapter16UC := app_content.NewGetChapter16UseCase(
		entitlementRepo,
		contentStorage,
		cfg.Services.DecryptionKey,
	)

	// Handlers
	authHandler := handlers.NewAuthHandler(registerUC, loginUC, logoutUC, refreshUC)
	articleHandler := handlers.NewArticleHandler(
		getArticleUC, createArticleUC, updateArticleUC,
		deleteArticleUC, listArticlesUC, publishArticleUC,
		addBookmarkUC, removeBookmarkUC, listBookmarksUC,
	)
	submissionHandler := handlers.NewSubmissionHandler(
		submitExerciseUC, analyzeSubmissionUC, getUserSubmissionsUC, gradeSubmissionUC,
	)
	contentHandler := handlers.NewContentHandler(getChapter16UC)
	communityHandler := handlers.NewCommunityHandler(
		createTopicUC, listTopicsUC, getTopicUC, createPostUC, listPostsUC, addVoteUC,
	)
	adminHandler := handlers.NewAdminHandler(listUsersUC, updateUserStatusUC, listModerationUC, moderateContentUC)
	healthHandler := handlers.NewHealthHandler()

	log.Println("INFO: All services and use cases initialized.")

	// --- 4. راه‌اندازی روتر و ثبت Validator سفارشی ---
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("password", custom_validator.PasswordPolicy)
	}

	router := interfaces.SetupRouter(
		cfg,
		redisClient,
		authHandler,
		articleHandler,
		submissionHandler,
		communityHandler,
		adminHandler,
		contentHandler,
		healthHandler,
	)

	// --- 5. راه‌اندازی و خاموش کردن امن سرور ---
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("FATAL: listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("INFO: Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("FATAL: Server forced to shutdown:", err)
	}

	log.Println("INFO: Server exiting gracefully.")
}