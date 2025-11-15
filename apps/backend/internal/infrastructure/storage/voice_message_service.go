// apps/backend/internal/infrastructure/storage/voice_message_service.go

package storage

import (
	"context"
	"fmt"
	"mime/multipart"

	"narrative-architecture/apps/backend/internal/domain/community" // برای تعریف نوع داده
)

// VoiceMessageService برای مدیریت پیام‌های صوتی است.
type VoiceMessageService struct {
	mediaService *MediaService
}

// NewVoiceMessageService یک نمونه جدید ایجاد می‌کند.
func NewVoiceMessageService(mediaService *MediaService) *VoiceMessageService {
	return &VoiceMessageService{mediaService: mediaService}
}

// UploadVoiceMessage یک پیام صوتی را آپلود می‌کند.
func (s *VoiceMessageService) UploadVoiceMessage(ctx context.Context, file *multipart.FileHeader) (string, error) {
	// پیام‌های صوتی در پوشه "voice" ذخیره می‌شوند.
	return s.mediaService.UploadFile(ctx, file, "voice")
}