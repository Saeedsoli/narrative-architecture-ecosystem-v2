// apps/backend/internal/application/auth/register_user_test.go

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
)

func TestRegisterUserUseCase_Execute(t *testing.T) {
	// 1. Arrange
	mockUserRepo := new(mocks.MockUserRepository)
	mockTokenRepo := new(mocks.MockTokenRepository)
	
	useCase := NewRegisterUserUseCase(mockUserRepo, mockTokenRepo, "secret", "refresh_secret", 15*time.Minute, 30*24*time.Hour)

	tests := []struct {
		name          string
		request       RegisterUserRequest
		setupMocks    func()
		expectError   bool
		expectedError string
	}{
		{
			name: "Successful registration",
			request: RegisterUserRequest{
				Email:    "test@example.com",
				Password: "Password123!",
				Username: "testuser",
				FullName: "Test User",
			},
			setupMocks: func() {
				mockUserRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(nil, errors.New("not found")).Once()
				mockUserRepo.On("Save", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil).Once()
				mockTokenRepo.On("SaveRefreshToken", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(nil).Once()
			},
			expectError: false,
		},
		{
			name: "Email already exists",
			request: RegisterUserRequest{
				Email:    "existing@example.com",
				Password: "Password123!",
			},
			setupMocks: func() {
				mockUserRepo.On("FindByEmail", mock.Anything, "existing@example.com").Return(&user.User{}, nil).Once()
			},
			expectError:   true,
			expectedError: "email already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks for each test
			mockUserRepo.Mock = mock.Mock{}
			mockTokenRepo.Mock = mock.Mock{}
			tt.setupMocks()

			// 2. Act
			res, err := useCase.Execute(context.Background(), tt.request)

			// 3. Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, tt.request.Email, res.User.Email)
				assert.NotEmpty(t, res.AccessToken)
				assert.NotEmpty(t, res.RefreshToken)
			}

			mockUserRepo.AssertExpectations(t)
			mockTokenRepo.AssertExpectations(t)
		})
	}
}