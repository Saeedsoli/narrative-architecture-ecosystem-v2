// apps/backend/internal/application/auth/login_user.go

package auth

import (
	"context"
	"errors"
	"time"

	"narrative-architecture/apps/backend/internal/domain/user"
	"narrative-architecture/apps/backend/pkg/jwt"
)

type LoginUserRequest struct {
	Email    string
	Password string
}

type LoginUserResponse struct {
	User         *user.User
	AccessToken  string
	RefreshToken string
}

type LoginUserUseCase struct {
	userRepo      user.UserRepository
	tokenRepo     user.TokenRepository
	accessSecret  string
	refreshSecret string
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewLoginUserUseCase(userRepo user.UserRepository, tokenRepo user.TokenRepository, accessSecret, refreshSecret string, accessTTL, refreshTTL time.Duration) *LoginUserUseCase {
	return &LoginUserUseCase{
		userRepo:      userRepo,
		tokenRepo:     tokenRepo,
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

func (uc *LoginUserUseCase) Execute(ctx context.Context, req LoginUserRequest) (*LoginUserResponse, error) {
	u, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !u.CheckPassword(req.Password) {
		return nil, errors.New("invalid email or password")
	}

	accessToken, refreshToken, err := jwt.GenerateTokens(u.ID, u.Email, u.Roles, uc.accessSecret, uc.refreshSecret, uc.accessTTL, uc.refreshTTL)
	if err != nil {
		return nil, err
	}

	if err := uc.tokenRepo.SaveRefreshToken(ctx, u.ID, refreshToken, uc.refreshTTL); err != nil {
		return nil, err
	}

	return &LoginUserResponse{
		User:         u,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}