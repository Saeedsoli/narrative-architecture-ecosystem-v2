package auth

import (
	"context"
	"errors"
	"time"

	"narrative-architecture/apps/backend/pkg/jwt"
)

type RefreshTokenUseCase struct {
	userRepo      UserRepository
	tokenRepo     TokenRepository
	accessSecret  string
	refreshSecret string
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

// ... NewRefreshTokenUseCase ...

func (uc *RefreshTokenUseCase) Execute(ctx context.Context, providedRefreshToken string) (*LoginUserResponse, error) {
	// 1. اعتبارسنجی Refresh Token
	claims, err := jwt.ValidateRefreshToken(providedRefreshToken, uc.refreshSecret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// 2. بررسی اینکه آیا توکن در دیتابیس وجود دارد و باطل نشده است
	if err := uc.tokenRepo.ValidateRefreshToken(ctx, claims.Subject, providedRefreshToken); err != nil {
		return nil, err
	}

	// 3. پیدا کردن کاربر
	u, err := uc.userRepo.FindByID(ctx, claims.Subject)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 4. تولید یک جفت توکن جدید (با چرخش Refresh Token)
	newAccessToken, newRefreshToken, err := jwt.GenerateTokens(u.ID, u.Email, u.Roles, uc.accessSecret, uc.refreshSecret, uc.accessTTL, uc.refreshTTL)
	if err != nil {
		return nil, err
	}

	// 5. باطل کردن توکن قدیمی و ذخیره توکن جدید
	if err := uc.tokenRepo.RotateRefreshToken(ctx, claims.Subject, providedRefreshToken, newRefreshToken, uc.refreshTTL); err != nil {
		return nil, err
	}

	return &LoginUserResponse{
		User:         u,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}