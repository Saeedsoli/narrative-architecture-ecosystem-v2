// apps/backend/internal/mocks/user_repository_mock.go

package mocks

import (
	"context"
	"errors"

	"github.com/stretchr/testify/mock"
	"narrative-architecture/apps/backend/internal/domain/user"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) Save(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}