// apps/backend/internal/application/auth/refresh_token_test.go

package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"narrative-architecture/apps/backend/internal/domain/user"
	"narrative-architecture/apps/backend/internal/mocks"
	"narrative-architecture/apps/backend/pkg/jwt"
)

func TestRefreshTokenUseCase_Execute(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockTokenRepo := new(mocks.MockTokenRepository)
	
	refreshSecret := "very-secret-refresh-key"
	useCase := NewRefreshTokenUseCase(mockUserRepo, mockTokenRepo, "secret", refreshSecret, 15*time.Minute, 30*24*time.Hour)

	testUser := user.NewUser("01HUSERID", "test@example.com", "testuser", "Test User")
	_, validRefreshToken, _ := jwt.GenerateTokens(testUser.ID, testUser.Email, testUser.Roles, "secret", refreshSecret, 15*time.Minute, 30*24*time.Hour)

	tests := []struct {
		name          string
		refreshToken  string
		setupMocks    func()
		expectError   bool
		expectedError string
	}{
		{
			name:         "Successful refresh",
			refreshToken: validRefreshToken,
			setupMocks: func() {
				mockTokenRepo.On("ValidateRefreshToken", mock.Anything, testUser.ID, validRefreshToken).Return(nil).Once()
				mockUserRepo.On("FindByID", mock.Anything, testUser.ID).Return(testUser, nil).Once()
				mockTokenRepo.On("RotateRefreshToken", mock.Anything, testUser.ID, validRefreshToken, mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(nil).Once()
			},
			expectError: false,
		},
		{
			name:         "Invalid token signature",
			refreshToken: "invalid.token.signature",
			setupMocks:   func() {},
			expectError:   true,
			expectedError: "invalid refresh token",
		},
		{
			name:         "Token not found in DB or already used",
			refreshToken: validRefreshToken,
			setupMocks: func() {
				mockTokenRepo.On("ValidateRefreshToken", mock.Anything, testUser.ID, validRefreshToken).Return(errors.New("token invalid")).Once()
			},
			expectError:   true,
			expectedError: "token invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo.Mock = mock.Mock{}
			mockTokenRepo.Mock = mock.Mock{}
			tt.setupMocks()

			res, err := useCase.Execute(context.Background(), tt.refreshToken)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.NotEmpty(t, res.AccessToken)
				assert.NotEmpty(t, res.RefreshToken)
			}
			
			mockUserRepo.AssertExpectations(t)
			mockTokenRepo.AssertExpectations(t)
		})
	}
}