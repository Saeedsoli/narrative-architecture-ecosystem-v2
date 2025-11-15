package community

import (
	"context"
	"time"
	"github.com/oklog/ulid/v2"
	"narrative-architecture/apps/backend/internal/domain/community"
)

type CreateTopicRequest struct {
	AuthorID   string
	AuthorName string
	Locale     string
	Title      string
	Body       string
	Tags       []string
}

type CreateTopicUseCase struct {
	repo community.Repository
}

func NewCreateTopicUseCase(repo community.Repository) *CreateTopicUseCase {
	return &CreateTopicUseCase{repo: repo}
}

func (uc *CreateTopicUseCase) Execute(ctx context.Context, req CreateTopicRequest) (*community.Topic, error) {
	newTopic := &community.Topic{
		ID:     ulid.New().String(),
		Locale: req.Locale,
		Title:  req.Title,
		Body:   req.Body,
		Tags:   req.Tags,
		Author: struct {
			ID       string
			Username string
		}{ID: req.AuthorID, Username: req.AuthorName},
		Status:    "open",
		CreatedAt: time.Now(),
	}

	if err := uc.repo.CreateTopic(ctx, newTopic); err != nil {
		return nil, err
	}
	
	return newTopic, nil
}