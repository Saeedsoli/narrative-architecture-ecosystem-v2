// apps/backend/internal/application/admin/moderate_content.go

package admin

import (
	"context"
	"errors"
	"narrative-architecture/apps/backend/internal/domain/moderation"
)

type ModerateContentRequest struct {
	ItemID      string
	ModeratorID string
	Action      string // "approve" or "reject"
	Reason      string
}

type ModerateContentUseCase struct {
	moderationRepo moderation.Repository
	contentRepo    ContentRepository // یک رابط عمومی برای آپدیت محتوا در MongoDB
}

// ... NewModerateContentUseCase ...

func (uc *ModerateContentUseCase) Execute(ctx context.Context, req ModerateContentRequest) error {
	// 1. تغییر وضعیت آیتم در صف بررسی
	status := moderation.StatusApproved
	if req.Action == "reject" {
		status = moderation.StatusRejected
	}
	
	err := uc.moderationRepo.UpdateQueueItemStatus(ctx, req.ItemID, status, req.ModeratorID, req.Reason)
	if err != nil {
		return err
	}

	// 2. اگر رد شده، محتوای اصلی را نیز آپدیت کن (مثلاً soft delete)
	if status == moderation.StatusRejected {
		item, err := uc.moderationRepo.FindQueueItemByID(ctx, req.ItemID)
		if err != nil { return err }

		// این بخش به یک رابط عمومی برای کار با محتواهای مختلف نیاز دارد
		return uc.contentRepo.Moderate(ctx, item.TargetType, item.TargetID, "rejected")
	}

	return nil
}