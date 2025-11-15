// apps/backend/internal/application/article/add_bookmark.go

package article

import (
	"context"
	"github.com/oklog/ulid/v2"
	"narrative-architecture/apps/backend/internal/domain/article"
)

type AddBookmarkRequest struct {
	UserID    string
	ArticleID string
}

type AddBookmarkUseCase struct {
	repo article.BookmarkRepository
}

func NewAddBookmarkUseCase(repo article.BookmarkRepository) *AddBookmarkUseCase {
	return &AddBookmarkUseCase{repo: repo}
}

func (uc *AddBookmarkUseCase) Execute(ctx context.Context, req AddBookmarkRequest) error {
	newBookmark := article.NewBookmark(
		ulid.New().String(),
		req.UserID,
		req.ArticleID,
	)
	return uc.repo.Create(ctx, newBookmark)
}