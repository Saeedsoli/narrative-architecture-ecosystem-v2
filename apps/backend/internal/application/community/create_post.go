// apps/backend/internal/application/community/create_post.go

package community

import (
	"context"
	"errors"
	"time"
	"github.com/oklog/ulid/v2"
	"narrative-architecture/apps/backend/internal/domain/community"
)

type CreatePostRequest struct {
	TopicID  string
	ParentID *string
	UserID   string
	Username string
	Avatar   string
	Body     string
}

type CreatePostUseCase struct {
	repo community.Repository
}

func NewCreatePostUseCase(repo community.Repository) *CreatePostUseCase {
	return &CreatePostUseCase{repo: repo}
}

func (uc *CreatePostUseCase) Execute(ctx context.Context, req CreatePostRequest) (*community.Post, error) {
	// بررسی وجود تاپیک
	topic, err := uc.repo.FindTopicByID(ctx, req.TopicID)
	if err != nil {
		return nil, errors.New("topic not found")
	}
	if topic.Status == "locked" {
		return nil, errors.New("topic is locked")
	}

	newPost := &community.Post{
		ID:       ulid.New().String(),
		TopicID:  req.TopicID,
		ParentID: req.ParentID,
		User: struct {
			ID       string
			Username string
			Avatar   string
		}{ID: req.UserID, Username: req.Username, Avatar: req.Avatar},
		Body:      req.Body,
		CreatedAt: time.Now(),
	}

	if err := uc.repo.CreatePost(ctx, newPost); err != nil {
		return nil, err
	}
	
	return newPost, nil
}