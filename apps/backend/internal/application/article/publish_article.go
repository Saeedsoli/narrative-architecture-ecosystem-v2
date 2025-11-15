// apps/backend/internal/application/article/publish_article.go

package article

import (
	"context"
	"errors"

	"narrative-architecture/apps/backend/internal/domain/article"
)

type PublishArticleUseCase struct {
	repo          article.Repository
	domainService *article.Service
	// eventBus      EventBus // برای ارسال رویداد
}

func NewPublishArticleUseCase(repo article.Repository, ds *article.Service) *PublishArticleUseCase {
	return &PublishArticleUseCase{repo: repo, domainService: ds}
}

func (uc *PublishArticleUseCase) Execute(ctx context.Context, articleID, userID string) error {
	// 1. دریافت مقاله از دیتابیس
	art, err := uc.repo.FindByID(ctx, articleID) // فرض بر وجود این متد در Repository
	if err != nil {
		return err
	}

	// 2. بررسی دسترسی (کاربر فقط می‌تواند مقاله خود را منتشر کند، مگر اینکه ادمین باشد)
	if art.Author.ID != userID /* && !user.IsAdmin() */ {
		return errors.New("user does not have permission to publish this article")
	}

	// 3. اجرای منطق تجاری انتشار با استفاده از سرویس دامنه
	if err := uc.domainService.Publish(art); err != nil {
		return err
	}

	// 4. ذخیره تغییرات در دیتابیس
	if err := uc.repo.Save(ctx, art); err != nil {
		return err
	}
	
	// 5. ارسال رویداد برای همگام‌سازی با Elasticsearch و ارسال نوتیفیکیشن
	// uc.eventBus.Publish("article.published", art)

	return nil
}