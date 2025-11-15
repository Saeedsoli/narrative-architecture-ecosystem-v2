package article

import (
	"context"
	"errors"
	"narrative-architecture/apps/backend/internal/domain/article"
)

type UpdateArticleRequest struct {
	ArticleID string
	UserID    string
	Title     *string
	Content   *string
	Excerpt   *string
	Tags      []string
}

type UpdateArticleUseCase struct {
	repo  article.Repository
	cache ArticleCache
}

func NewUpdateArticleUseCase(repo article.Repository, cache ArticleCache) *UpdateArticleUseCase {
	return &UpdateArticleUseCase{repo: repo, cache: cache}
}

func (uc *UpdateArticleUseCase) Execute(ctx context.Context, req UpdateArticleRequest) (*article.Article, error) {
	art, err := uc.repo.FindByID(ctx, req.ArticleID)
	if err != nil {
		return nil, err
	}

	if art.Author.ID != req.UserID {
		return nil, errors.New("user does not have permission to update this article")
	}

	if req.Title != nil { art.Title = *req.Title }
	if req.Content != nil { art.Content = *req.Content }
	if req.Excerpt != nil { art.Excerpt = *req.Excerpt }
	if req.Tags != nil { art.Metadata.Tags = req.Tags }

	if err := uc.repo.Save(ctx, art); err != nil {
		return nil, err
	}
	
	go uc.cache.Delete(context.Background(), art)

	return art, nil
}