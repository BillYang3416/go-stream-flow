package usecase

import (
	"context"
	"testing"

	"github.com/bgg/go-flow-gateway/internal/entity"
	"github.com/bgg/go-flow-gateway/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserProfileRepo struct {
	mock.Mock
}

func (m *MockUserProfileRepo) Create(ctx context.Context, u entity.UserProfile) (entity.UserProfile, error) {
	args := m.Called(ctx, u)
	return args.Get(0).(entity.UserProfile), args.Error(1)
}

func (m *MockUserProfileRepo) GetByID(ctx context.Context, userId int) (entity.UserProfile, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(entity.UserProfile), args.Error(1)
}

func setupUserProfileUseCase(t *testing.T) (*UserProfileUseCase, *MockUserProfileRepo) {
	t.Helper()

	mockRepo := new(MockUserProfileRepo)
	uc := NewUserProfileUseCase(mockRepo, logger.New("debug"))
	return uc, mockRepo
}

func TestUserProfileUsecase_Create(t *testing.T) {

	const (
		userID      = 1234567890
		displayName = "hank"
		pictureURL  = "https://example.com"
	)

	t.Run("Create user profile successfully", func(t *testing.T) {
		// Arrange
		uc, mockRepo := setupUserProfileUseCase(t)
		ctx := context.Background()

		userProfile := entity.UserProfile{
			UserID:      userID,
			DisplayName: displayName,
			PictureURL:  pictureURL,
		}

		mockRepo.On("Create", ctx, userProfile).Return(userProfile, nil)

		// Act
		result, err := uc.Create(ctx, userProfile)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, userProfile, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Create user profile with invalid input", func(t *testing.T) {
		// Arrange
		uc, mockRepo := setupUserProfileUseCase(t)
		ctx := context.Background()

		userProfile := entity.UserProfile{}

		mockRepo.On("Create", ctx, userProfile).Return(userProfile, assert.AnError)

		// Act
		result, err := uc.Create(ctx, userProfile)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, entity.UserProfile{}, result)
		mockRepo.AssertExpectations(t)
	})

}

func TestUserProfileUsecase_GetByID(t *testing.T) {

	const (
		userID      = 1234567890
		displayName = "hank"
		pictureURL  = "https://example.com"
	)

	t.Run("Get user profile by id successfully", func(t *testing.T) {
		// Arrange
		uc, mockRepo := setupUserProfileUseCase(t)
		ctx := context.Background()

		userProfile := entity.UserProfile{
			UserID:      userID,
			DisplayName: displayName,
			PictureURL:  pictureURL,
		}

		mockRepo.On("GetByID", ctx, userProfile.UserID).Return(userProfile, nil)

		// Act
		result, err := uc.GetByID(ctx, userProfile.UserID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, userProfile, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Get user profile by id with invalid user ID", func(t *testing.T) {
		// Arrange
		uc, mockRepo := setupUserProfileUseCase(t)
		ctx := context.Background()

		mockRepo.On("GetByID", ctx, userID).Return(entity.UserProfile{}, assert.AnError)

		// Act
		result, err := uc.GetByID(ctx, userID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, entity.UserProfile{}, result)
		mockRepo.AssertExpectations(t)
	})

}
