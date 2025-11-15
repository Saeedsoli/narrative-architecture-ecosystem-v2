// apps/backend/internal/application/community/add_comment.go

package community

import (
	"context"
	"errors"
	"time"

	"github.com/oklog/ulid/v2"
	"narrative-architecture/apps/backend/internal/domain/community"
)

// AddCommentRequest ساختار درخواست برای افزودن یک نظر است.
type AddCommentRequest struct {
	ArticleID string
	UserID    string
	Username  string
	Avatar    string
	ParentID  *string
	Body      string
}

// AddCommentUseCase منطق تجاری برای افزودن نظر را کپسوله می‌کند.
type AddCommentUseCase struct {
	repo community.Repository
}

// NewAddCommentUseCase یک نمونه جدید از AddCommentUseCase ایجاد می‌کند.
func NewAddCommentUseCase(repo community.Repository) *AddCommentUseCase {
	return &AddCommentUseCase{repo: repo}
}

// Execute متد اصلی برای اجرای Use Case است.
func (uc *AddCommentUseCase) Execute(ctx context.Context, req AddCommentRequest) (*community.Comment, error) {
	if req.Body == "" {
		return nil, errors.New("comment body cannot be empty")
	}

	newComment := &community.Comment{
		ID:        ulid.New().String(),
		ArticleID: req.ArticleID,
		UserID:    req.UserID,
		User: struct {
			ID       string
			Username string
			Avatar   string
		}{ID: req.UserID, Username: req.Username, Avatar: req.Avatar},
		ParentID:  req.ParentID,
		Body:      req.Body,
		CreatedAt: time.Now(),
	}

	// منطق تعیین thread_id
	if req.ParentID != nil {
		parentComment, err := uc.repo.FindCommentByID(ctx, *req.ParentID)
		if err != nil {
			return nil, errors.New("parent comment not found")
		}
		newComment.ThreadID = parentComment.ThreadID
	} else {
		newComment.ThreadID = newComment.ID
	}

	if err := uc.repo.CreateComment(ctx, newComment); err != nil {
		return nil, err
	}

	return newComment, nil
}