package article

import (
	"context"
	
	"narrative-architecture/apps/backend/internal/domain/article"
)

type PublishArticleUseCase struct {
	repo         article.Repository
	domainService *article.Service // <-- استفاده از سرویس دامنه
}

func NewPublishArticleUseCase(repo article.Repository, domainService *article.Service) *PublishArticleUseCase {
	return &PublishArticleUseCase{
		repo:         repo,
		domainService: domainService,
	}
}

func (uc *PublishArticleUseCase) Execute(ctx context.Context, articleID string) error {
	// 1. دریافت مقاله از دیتابیس
	art, err := uc.repo.FindByID(ctx, articleID)
	if err != nil {
		return err
	}

	// 2. اجرای منطق تجاری با استفاده از سرویس دامنه
	if err := uc.domainService.Publish(art); err != nil {
		return err // اگر مقاله شرایط انتشار را نداشته باشد، خطا برمی‌گردد
	}
	
	// (اختیاری) محاسبه مجدد زمان مطالعه
	art.Metadata.ReadTime = uc.domainService.CalculateReadingTime(art.Content)

	// 3. ذخیره تغییرات در دیتابیس
	if err := uc.repo.Save(ctx, art); err != nil {
		return err
	}
	
	// 4. ارسال رویداد برای سایر سرویس‌ها
	// eventBus.Publish("article.published", art.ID)

	return nil
}