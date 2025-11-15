package article

import (
	"context"
	"errors"
)

type DeleteArticleUseCase struct {
	repo  article.Repository
	cache ArticleCache
}

func NewDeleteArticleUseCase(repo article.Repository, cache ArticleCache) *DeleteArticleUseCase {
	return &DeleteArticleUseCase{repo: repo, cache: cache}
}

func (uc *DeleteArticleUseCase) Execute(ctx context.Context, articleID, userID string) error {
	art, err := uc.repo.FindByID(ctx, articleID)
	if err != nil { return err }

	if art.Author.ID != userID {
		return errors.New("user does not have permission to delete this article")
	}

	if err := uc.repo.Delete(ctx, articleID); err != nil {
		return err
	}
	
	go uc.cache.Delete(context.Background(), art)

	return nil
}