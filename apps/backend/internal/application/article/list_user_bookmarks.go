// apps/backend/internal/application/article/list_user_bookmarks.go

package article

import (
	"context"
	"log"
	"narrative-architecture/apps/backend/internal/domain/article"
)

type ListUserBookmarksUseCase struct {
	bookmarkRepo article.BookmarkRepository
	articleRepo  article.Repository
}

func NewListUserBookmarksUseCase(bookmarkRepo article.BookmarkRepository, articleRepo article.Repository) *ListUserBookmarksUseCase {
	return &ListUserBookmarksUseCase{
		bookmarkRepo: bookmarkRepo,
		articleRepo:  articleRepo,
	}
}

func (uc *ListUserBookmarksUseCase) Execute(ctx context.Context, userID string) ([]*article.Article, error) {
	bookmarks, err := uc.bookmarkRepo.FindByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(bookmarks) == 0 {
		return []*article.Article{}, nil
	}

	var articleIDs []string
	for _, b := range bookmarks {
		articleIDs = append(articleIDs, b.ArticleID)
	}

	// برای این کار، باید یک فیلتر جدید به Repository اضافه کنیم
	filter := article.Filter{IDs: articleIDs}
	articles, _, err := uc.articleRepo.Find(ctx, filter)
	if err != nil {
		log.Printf("Failed to fetch bookmarked articles details: %v", err)
		return []*article.Article{}, nil
	}
	
	return articles, nil
}