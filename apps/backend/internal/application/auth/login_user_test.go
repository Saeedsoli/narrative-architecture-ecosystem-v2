// apps/backend/internal/application/auth/login_user_test.go

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

func TestLoginUserUseCase_Execute(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockTokenRepo := new(mocks.MockTokenRepository)

	useCase := NewLoginUserUseCase(mockUserRepo, mockTokenRepo, "secret", "refresh_secret", 15*time.Minute, 30*24*time.Hour)

	// یک کاربر نمونه با رمز عبور هش شده
	testUser := user.NewUser("01H...", "test@example.com", "testuser", "Test User")
	_ = testUser.SetPassword("Password123!")

	tests := []struct {
		name          string
		request       LoginUserRequest
		setupMocks    func()
		expectError   bool
		expectedError string
	}{
		{
			name: "Successful login",
			request: LoginUserRequest{
				Email:    "test@example.com",
				Password: "Password123!",
			},
			setupMocks: func() {
				mockUserRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(testUser, nil).Once()
				mockTokenRepo.On("SaveRefreshToken", mock.Anything, testUser.ID, mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(nil).Once()
			},
			expectError: false,
		},
		{
			name: "User not found",
			request: LoginUserRequest{
				Email:    "notfound@example.com",
				Password: "Password123!",
			},
			setupMocks: func() {
				mockUserRepo.On("FindByEmail", mock.Anything, "notfound@example.com").Return(nil, errors.New("not found")).Once()
			},
			expectError:   true,
			expectedError: "invalid credentials",
		},
		{
			name: "Incorrect password",
			request: LoginUserRequest{
				Email:    "test@example.com",
				Password: "WrongPassword",
			},
			setupMocks: func() {
				mockUserRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(testUser, nil).Once()
			},
			expectError:   true,
			expectedError: "invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo.Mock = mock.Mock{}
			mockTokenRepo.Mock = mock.Mock{}
			tt.setupMocks()

			res, err := useCase.Execute(context.Background(), tt.request)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, testUser.ID, res.User.ID)
			}
			
			mockUserRepo.AssertExpectations(t)
			mockTokenRepo.AssertExpectations(t)
		})
	}
}