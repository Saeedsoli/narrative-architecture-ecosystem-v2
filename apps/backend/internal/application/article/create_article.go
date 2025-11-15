// apps/backend/internal/application/article/create_article.go

package article

import (
	"context"
	"time"

	"narrative-architecture/apps/backend/internal/domain/article"
	"narrative-architecture/apps/backend/lib/utils"
	"github.com/oklog/ulid/v2"
)

type CreateArticleRequest struct {
	AuthorID       string
	Locale         string
	Title          string
	Content        string
	Excerpt        string
	Tags           []string
	Category       string
}

type CreateArticleUseCase struct {
	repo article.Repository
}

func NewCreateArticleUseCase(repo article.Repository) *CreateArticleUseCase {
	return &CreateArticleUseCase{repo: repo}
}

func (uc *CreateArticleUseCase) Execute(ctx context.Context, req CreateArticleRequest) (*article.Article, error) {
	now := time.Now()
	
	newArticle := &article.Article{
		ID:             ulid.New().String(),
		ContentGroupID: ulid.New().String(), // یک گروه جدید برای این مقاله ایجاد می‌شود
		Locale:         req.Locale,
		Title:          req.Title,
		Content:        req.Content,
		Excerpt:        req.Excerpt,
		Author: struct {
			ID     string
			Name   string // باید از userRepo گرفته شود
			Avatar string
		}{ID: req.AuthorID},
		Metadata: struct {
			Tags       []string
			Category   string
			ReadTime   int
			Difficulty string
		}{
			Tags:     req.Tags,
			Category: req.Category,
			ReadTime: utils.CalculateReadingTime(req.Content),
		},
		Status:    article.StatusDraft,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// تولید خودکار اسلاگ در لایه Repository انجام می‌شود
	if err := uc.repo.Save(ctx, newArticle); err != nil {
		return nil, err
	}
	
	return newArticle, nil
}