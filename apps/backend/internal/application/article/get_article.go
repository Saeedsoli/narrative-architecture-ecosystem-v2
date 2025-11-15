// apps/backend/internal/application/article/get_article.go

package article

import (
	"context"
	"log"

	"narrative-architecture/apps/backend/internal/domain/article"
)

// ArticleCache رابطی برای کش کردن مقالات است.
type ArticleCache interface {
	GetBySlug(ctx context.Context, slug string) (*article.Article, error)
	Set(ctx context.Context, art *article.Article) error
	Delete(ctx context.Context, art *article.Article) error
}

// GetArticleUseCase منطق تجاری برای دریافت یک مقاله را کپسوله می‌کند.
type GetArticleUseCase struct {
	repo  article.Repository
	cache ArticleCache
}

// NewGetArticleUseCase یک نمونه جدید از GetArticleUseCase ایجاد می‌کند.
func NewGetArticleUseCase(repo article.Repository, cache ArticleCache) *GetArticleUseCase {
	return &GetArticleUseCase{repo: repo, cache: cache}
}

// Execute متد اصلی برای اجرای Use Case است.
func (uc *GetArticleUseCase) Execute(ctx context.Context, slug string) (*article.Article, error) {
	// 1. ابتدا سعی کن از کش بخوانی
	if cachedArticle, err := uc.cache.GetBySlug(ctx, slug); err == nil {
		log.Println("Cache HIT for article:", slug)
		return cachedArticle, nil
	}
	log.Println("Cache MISS for article:", slug)

	// 2. اگر در کش نبود، از دیتابیس بخوان
	art, err := uc.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	// 3. ترجمه‌های موجود را پیدا کن
	translations, err := uc.repo.FindTranslations(ctx, art.ContentGroupID)
	if err != nil {
		log.Printf("Failed to get translations for article %s: %v", art.ID, err)
	} else {
		art.Translations = translations
	}

	// 4. نتیجه را برای درخواست‌های بعدی در کش ذخیره کن
	go func() {
		if err := uc.cache.Set(context.Background(), art); err != nil {
			log.Printf("Failed to set cache for article %s: %v", art.ID, err)
		}
	}()

	return art, nil
}