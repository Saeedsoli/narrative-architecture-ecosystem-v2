// apps/backend/internal/domain/article/service.go

package article

import (
	"errors"
	"time"
)

// Service کپسوله کننده منطق تجاری خالص مربوط به مقالات است.
// این سرویس به هیچ رابط خارجی (مانند دیتابیس) دسترسی ندارد.
type Service struct{}

// NewService یک نمونه جدید از سرویس دامنه مقاله ایجاد می‌کند.
func NewService() *Service {
	return &Service{}
}

// CanPublish بررسی می‌کند که آیا یک مقاله شرایط لازم برای انتشار را دارد یا خیر.
// این یک مثال از منطق تجاری خالص است.
func (s *Service) CanPublish(art *Article) (bool, error) {
	if art.Title == "" {
		return false, errors.New("article title cannot be empty for publishing")
	}
	if len(art.Content) < 100 { // مثال: حداقل 100 کاراکتر محتوا
		return false, errors.New("article content is too short for publishing")
	}
	if len(art.Metadata.Tags) == 0 {
		return false, errors.New("article must have at least one tag for publishing")
	}
	return true, nil
}

// Publish مقاله را برای انتشار آماده می‌کند.
// این متد وضعیت مقاله را تغییر داده و تاریخ انتشار را تنظیم می‌کند.
func (s *Service) Publish(art *Article) error {
	canPublish, err := s.CanPublish(art)
	if !canPublish {
		return err
	}

	if art.Status == StatusPublished {
		return errors.New("article is already published")
	}

	art.Status = StatusPublished
	now := time.Now().UTC()
	art.PublishedAt = &now
	art.UpdatedAt = now

	return nil
}

// Unpublish یک مقاله را از حالت انتشار خارج می‌کند.
func (s *Service) Unpublish(art *Article) {
	art.Status = StatusDraft
	art.PublishedAt = nil
	art.UpdatedAt = time.Now().UTC()
}

// CalculateReadingTime زمان مطالعه تخمینی یک مقاله را محاسبه می‌کند.
// فرض: 200 کلمه در دقیقه.
func (s *Service) CalculateReadingTime(content string) int {
	wordCount := len(strings.Fields(content))
	readingTime := wordCount / 200
	if readingTime == 0 {
		return 1 // حداقل 1 دقیقه
	}
	return readingTime
}