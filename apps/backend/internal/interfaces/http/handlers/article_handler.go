// apps/backend/internal/interfaces/http/handlers/article_handler.go

package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	app_article "narrative-architecture/apps/backend/internal/application/article"
	"narrative-architecture/apps/backend/internal/domain/article"
	"narrative-architecture/apps/backend/internal/interfaces/http/dto"
)

// ArticleHandler کنترلر HTTP برای تمام عملیات مربوط به مقالات و بوکمارک‌ها است.
type ArticleHandler struct {
	getArticleUC     *app_article.GetArticleUseCase
	createArticleUC  *app_article.CreateArticleUseCase
	updateArticleUC  *app_article.UpdateArticleUseCase
	deleteArticleUC  *app_article.DeleteArticleUseCase
	listArticlesUC   *app_article.ListArticlesUseCase
	publishArticleUC *app_article.PublishArticleUseCase
	addBookmarkUC    *app_article.AddBookmarkUseCase
	removeBookmarkUC *app_article.RemoveBookmarkUseCase
	listBookmarksUC  *app_article.ListUserBookmarksUseCase
}

// NewArticleHandler یک نمونه جدید از ArticleHandler با تمام وابستگی‌های لازم ایجاد می‌کند.
func NewArticleHandler(
	getUC *app_article.GetArticleUseCase, createUC *app_article.CreateArticleUseCase,
	updateUC *app_article.UpdateArticleUseCase, deleteUC *app_article.DeleteArticleUseCase,
	listUC *app_article.ListArticlesUseCase, publishUC *app_article.PublishArticleUseCase,
	addBookmarkUC *app_article.AddBookmarkUseCase, removeBookmarkUC *app_article.RemoveBookmarkUseCase,
	listBookmarksUC *app_article.ListUserBookmarksUseCase,
) *ArticleHandler {
	return &ArticleHandler{
		getArticleUC:     getUC,
		createArticleUC:  createUC,
		updateArticleUC:  updateUC,
		deleteArticleUC:  deleteUC,
		listArticlesUC:   listUC,
		publishArticleUC: publishUC,
		addBookmarkUC:    addBookmarkUC,
		removeBookmarkUC: removeBookmarkUC,
		listBookmarksUC:  listBookmarksUC,
	}
}

// ListArticles لیستی از مقالات منتشر شده را با فیلتر و صفحه‌بندی برمی‌گرداند.
func (h *ArticleHandler) ListArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	
	filter := article.Filter{
		Page:     page,
		PageSize: pageSize,
		Tags:     c.QueryArray("tags"),
		Category: c.Query("category"),
		Status:   string(article.StatusPublished), // فقط مقالات منتشر شده برای عموم
		Locale:   c.DefaultQuery("locale", "fa"),
	}

	res, err := h.listArticlesUC.Execute(c.Request.Context(), filter)
	if err != nil {
		HandleError(c, err)
		return
	}

	var articleDTOs []*dto.ArticleResponse
	for _, art := range res.Articles {
		articleDTOs = append(articleDTOs, dto.ToArticleResponse(art))
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data":       articleDTOs,
		"total":      res.Total,
		"page":       res.Page,
		"pageSize":   res.PageSize,
		"totalPages": res.TotalPages,
	})
}

// GetArticle یک مقاله را بر اساس اسلاگ آن برمی‌گرداند.
func (h *ArticleHandler) GetArticle(c *gin.Context) {
	slug := c.Param("slug")

	article, err := h.getArticleUC.Execute(c.Request.Context(), slug)
	if err != nil {
		HandleError(c, err)
		return
	}
	
	// بررسی اینکه آیا مقاله منتشر شده یا کاربر دسترسی ادمین دارد
	// if !article.IsPublished() && !userHasRole(c, "admin") {
	//	HandleError(c, ErrPermissionDenied)
	//	return
	// }

	c.JSON(http.StatusOK, dto.ToArticleResponse(article))
}

// CreateArticle یک مقاله جدید در حالت پیش‌نویس ایجاد می‌کند.
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var req app_article.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, err)
		return
	}

	req.AuthorID = c.GetString("userID")

	article, err := h.createArticleUC.Execute(c.Request.Context(), req)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.ToArticleResponse(article))
}

// UpdateArticle یک مقاله موجود را آپدیت می‌کند.
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	articleID := c.Param("id")
	var req app_article.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, err)
		return
	}
	req.ArticleID = articleID
	req.UserID = c.GetString("userID")

	updatedArticle, err := h.updateArticleUC.Execute(c.Request.Context(), req)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.ToArticleResponse(updatedArticle))
}

// DeleteArticle یک مقاله را به‌صورت نرم حذف می‌کند.
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	articleID := c.Param("id")
	userID := c.GetString("userID")

	if err := h.deleteArticleUC.Execute(c.Request.Context(), articleID, userID); err != nil {
		HandleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// PublishArticle وضعیت یک مقاله را به "منتشر شده" تغییر می‌دهد.
func (h *ArticleHandler) PublishArticle(c *gin.Context) {
	articleID := c.Param("id")
	userID := c.GetString("userID")
	
	if err := h.publishArticleUC.Execute(c.Request.Context(), articleID, userID); err != nil {
		HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Article published successfully"})
}

// AddBookmark یک مقاله را به لیست بوکمارک‌های کاربر اضافه می‌کند.
func (h *ArticleHandler) AddBookmark(c *gin.Context) {
	articleID := c.Param("id")
	userID := c.GetString("userID")

	err := h.addBookmarkUC.Execute(c.Request.Context(), app_article.AddBookmarkRequest{
		UserID:    userID,
		ArticleID: articleID,
	})
	
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Bookmark added successfully"})
}

// RemoveBookmark یک بوکمارک را از لیست کاربر حذف می‌کند.
func (h *ArticleHandler) RemoveBookmark(c *gin.Context) {
	articleID := c.Param("id")
	userID := c.GetString("userID")

	err := h.removeBookmarkUC.Execute(c.Request.Context(), userID, articleID)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// ListUserBookmarks لیست مقالات بوکمارک شده توسط کاربر را برمی‌گرداند.
func (h *ArticleHandler) ListUserBookmarks(c *gin.Context) {
	userID := c.GetString("userID")

	articles, err := h.listBookmarksUC.Execute(c.Request.Context(), userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	var response []*dto.ArticleResponse
	for _, art := range articles {
		response = append(response, dto.ToArticleResponse(art))
	}
	
	c.JSON(http.StatusOK, response)
}