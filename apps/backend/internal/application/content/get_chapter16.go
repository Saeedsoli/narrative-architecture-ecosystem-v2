package content

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"errors"
	"fmt"

	"narrative-architecture/apps/backend/internal/domain/entitlement"
)

// ContentStorage رابطی برای خواندن محتوای خام است.
type ContentStorage interface {
	GetEncryptedChapter16(ctx context.Context) ([]byte, error)
}

type GetChapter16Response struct {
	Content   json.RawMessage `json:"content"` // محتوای رمزگشایی شده
	Watermark string          `json:"watermark"`
}

type GetChapter16UseCase struct {
	entitlementRepo entitlement.Repository
	contentStorage  ContentStorage
	decryptionKey   []byte // این کلید باید از یک جای امن خوانده شود
}

func NewGetChapter16UseCase(
	entitlementRepo entitlement.Repository,
	storage ContentStorage,
	key string,
) *GetChapter16UseCase {
	return &GetChapter16UseCase{
		entitlementRepo: entitlementRepo,
		contentStorage:  storage,
		decryptionKey:   []byte(key),
	}
}

func (uc *GetChapter16UseCase) Execute(ctx context.Context, userID, userIP string) (*GetChapter16Response, error) {
	// 1. بررسی حق دسترسی کاربر
	hasAccess, err := uc.entitlementRepo.HasAccess(ctx, userID, "chapter", "16")
	if err != nil {
		return nil, err
	}
	if !hasAccess {
		return nil, errors.New("user does not have access to this chapter")
	}

	// 2. خواندن محتوای رمزنگاری شده از S3
	encryptedData, err := uc.contentStorage.GetEncryptedChapter16(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get content: %w", err)
	}

	// 3. رمزگشایی محتوا (AES-GCM)
	block, err := aes.NewCipher(uc.decryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create gcm: %w", err)
	}
	
	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	
	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]
	decryptedData, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt content: %w", err)
	}

	// 4. ایجاد Watermark
	watermarkText := fmt.Sprintf("%s|%s|%d", userID, userIP, time.Now().Unix())

	return &GetChapter16Response{
		Content:   json.RawMessage(decryptedData),
		Watermark: watermarkText,
	}, nil
}