// apps/backend/internal/interfaces/http/handlers/community_handler.go

package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	app_community "narrative-architecture/apps/backend/internal/application/community"
	"narrative-architecture/apps/backend/internal/interfaces/http/dto"
)

// CommunityHandler کنترلر HTTP برای تمام عملیات مربوط به انجمن گفتگو است.
type CommunityHandler struct {
	createTopicUC *app_community.CreateTopicUseCase
	listTopicsUC  *app_community.ListTopicsUseCase
	getTopicUC    *app_community.GetTopicUseCase
	createPostUC  *app_community.CreatePostUseCase
	listPostsUC   *app_community.ListPostsUseCase
	addVoteUC     *app_community.AddVoteUseCase
}

// NewCommunityHandler یک نمونه جدید از CommunityHandler ایجاد می‌کند.
func NewCommunityHandler(
	createTopicUC *app_community.CreateTopicUseCase,
	listTopicsUC *app_community.ListTopicsUseCase,
	getTopicUC *app_community.GetTopicUseCase,
	createPostUC *app_community.CreatePostUseCase,
	listPostsUC *app_community.ListPostsUseCase,
	addVoteUC *app_community.AddVoteUseCase,
) *CommunityHandler {
	return &CommunityHandler{
		createTopicUC: createTopicUC,
		listTopicsUC:  listTopicsUC,
		getTopicUC:    getTopicUC,
		createPostUC:  createPostUC,
		listPostsUC:   listPostsUC,
		addVoteUC:     addVoteUC,
	}
}

// ListTopics لیستی از تاپیک‌ها را با صفحه‌بندی برمی‌گرداند.
// GET /api/v1/community/topics?locale=fa&page=1&pageSize=10
func (h *CommunityHandler) ListTopics(c *gin.Context) {
	locale := c.DefaultQuery("locale", "fa")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	res, err := h.listTopicsUC.Execute(c.Request.Context(), locale, page, pageSize)
	if err != nil {
		HandleError(c, err)
		return
	}

	var topicDTOs []*dto.TopicResponse
	for _, topic := range res.Topics {
		topicDTOs = append(topicDTOs, dto.ToTopicResponse(topic, false))
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       topicDTOs,
		"total":      res.Total,
		"page":       res.Page,
		"pageSize":   res.PageSize,
		"totalPages": res.TotalPages,
	})
}

// GetTopic جزئیات یک تاپیک خاص را برمی‌گرداند.
// GET /api/v1/community/topics/:id
func (h *CommunityHandler) GetTopic(c *gin.Context) {
	topicID := c.Param("id")

	topic, err := h.getTopicUC.Execute(c.Request.Context(), topicID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToTopicResponse(topic, true))
}

// CreateTopic یک تاپیک جدید در انجمن ایجاد می‌کند.
// POST /api/v1/community/topics
func (h *CommunityHandler) CreateTopic(c *gin.Context) {
	var req struct {
		Locale string   `json:"locale" binding:"required,oneof=fa en"`
		Title  string   `json:"title" binding:"required,min=5"`
		Body   string   `json:"body" binding:"required,min=10"`
		Tags   []string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	useCaseReq := app_community.CreateTopicRequest{
		AuthorID:   c.GetString("userID"),
		AuthorName: c.GetString("username"), // فرض بر وجود این در Middleware
		Locale:     req.Locale,
		Title:      req.Title,
		Body:       req.Body,
		Tags:       req.Tags,
	}

	topic, err := h.createTopicUC.Execute(c.Request.Context(), useCaseReq)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.ToTopicResponse(topic, true))
}

// ListPosts لیستی از پست‌های یک تاپیک را با صفحه‌بندی برمی‌گرداند.
// GET /api/v1/community/topics/:id/posts?page=1&pageSize=20
func (h *CommunityHandler) ListPosts(c *gin.Context) {
	topicID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	res, err := h.listPostsUC.Execute(c.Request.Context(), topicID, page, pageSize)
	if err != nil {
		HandleError(c, err)
		return
	}
	
	var postDTOs []*dto.PostResponse
	for _, post := range res.Posts {
		postDTOs = append(postDTOs, dto.ToPostResponse(post))
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       postDTOs,
		"total":      res.Total,
		"page":       res.Page,
		"pageSize":   res.PageSize,
		"totalPages": res.TotalPages,
	})
}

// CreatePost یک پست جدید در یک تاپیک ایجاد می‌کند.
// POST /api/v1/community/topics/:id/posts
func (h *CommunityHandler) CreatePost(c *gin.Context) {
	topicID := c.Param("id")
	var req struct {
		Body     string  `json:"body" binding:"required,min=1"`
		ParentID *string `json:"parentId"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	useCaseReq := app_community.CreatePostRequest{
		TopicID:  topicID,
		ParentID: req.ParentID,
		UserID:   c.GetString("userID"),
		Username: c.GetString("username"),
		Avatar:   c.GetString("avatar"), // فرض بر وجود این در Middleware
		Body:     req.Body,
	}

	post, err := h.createPostUC.Execute(c.Request.Context(), useCaseReq)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.ToPostResponse(post))
}

// AddVote یک رأی (لایک/دیسلایک) برای یک پست یا کامنت ثبت می‌کند.
// POST /api/v1/community/votes
func (h *CommunityHandler) AddVote(c *gin.Context) {
	var req struct {
		TargetID   string `json:"targetId" binding:"required"`
		TargetType string `json:"targetType" binding:"required,oneof=forum_post article_comment"`
		Value      int    `json:"value" binding:"required,oneof=1 -1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userID := c.GetString("userID")

	useCaseReq := app_community.AddVoteRequest{
		UserID:     userID,
		TargetID:   req.TargetID,
		TargetType: req.TargetType,
		Value:      req.Value,
	}

	if err := h.addVoteUC.Execute(c.Request.Context(), useCaseReq); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vote registered successfully"})
}