// apps/backend/internal/application/article/remove_bookmark.go

package article

import (
	"context"
	"narrative-architecture/apps/backend/internal/domain/article"
)

type RemoveBookmarkUseCase struct {
	repo article.BookmarkRepository
}

func NewRemoveBookmarkUseCase(repo article.BookmarkRepository) *RemoveBookmarkUseCase {
	return &RemoveBookmarkUseCase{repo: repo}
}

func (uc *RemoveBookmarkUseCase) Execute(ctx context.Context, userID, articleID string) error {
	return uc.repo.Delete(ctx, userID, articleID)
}