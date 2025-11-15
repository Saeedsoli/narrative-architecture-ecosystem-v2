// apps/backend/internal/mocks/token_repository_mock.go

package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockTokenRepository struct {
	mock.Mock
}

func (m *MockTokenRepository) SaveRefreshToken(ctx context.Context, userID, token string, ttl time.Duration) error {
	args := m.Called(ctx, userID, token, ttl)
	return args.Error(0)
}

func (m *MockTokenRepository) ValidateRefreshToken(ctx context.Context, userID, token string) error {
	args := m.Called(ctx, userID, token)
	return args.Error(0)
}

func (m *MockTokenRepository) RotateRefreshToken(ctx context.Context, userID, oldToken, newToken string, ttl time.Duration) error {
	args := m.Called(ctx, userID, oldToken, newToken, ttl)
	return args.Error(0)
}

func (m *MockTokenRepository) InvalidateRefreshToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}