package usecase

import (
	"context"

	"github.com/bgg/go-flow-gateway/internal/entity"
	"github.com/stretchr/testify/mock"
)

type MockUserProfileUseCase struct {
	mock.Mock
}

func (m *MockUserProfileUseCase) Create(ctx context.Context, userProfile entity.UserProfile) (entity.UserProfile, error) {
	args := m.Called(ctx, userProfile)
	return args.Get(0).(entity.UserProfile), args.Error(1)
}

func (m *MockUserProfileUseCase) GetByID(ctx context.Context, userID int) (entity.UserProfile, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(entity.UserProfile), args.Error(1)
}
