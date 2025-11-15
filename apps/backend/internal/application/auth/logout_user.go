// apps/backend/internal/application/auth/logout_user.go

package auth

import (
	"context"
	"narrative-architecture/apps/backend/internal/domain/user"
)

type LogoutUserUseCase struct {
	tokenRepo user.TokenRepository
}

func NewLogoutUserUseCase(tokenRepo user.TokenRepository) *LogoutUserUseCase {
	return &LogoutUserUseCase{tokenRepo: tokenRepo}
}

func (uc *LogoutUserUseCase) Execute(ctx context.Context, refreshToken string) error {
	return uc.tokenRepo.InvalidateRefreshToken(ctx, refreshToken)
}