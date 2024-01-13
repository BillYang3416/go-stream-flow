package usecase

import (
	"context"
	"testing"

	"github.com/bgg/go-flow-gateway/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserProfileRepo struct {
	mock.Mock
}

func (m *MockUserProfileRepo) Create(ctx context.Context, u entity.UserProfile) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserProfileRepo) GetByID(ctx context.Context, userId string) (entity.UserProfile, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(entity.UserProfile), args.Error(1)
}

func (m *MockUserProfileRepo) UpdateRefreshToken(ctx context.Context, userId string, refreshToken string) error {
	args := m.Called(ctx, userId, refreshToken)
	return args.Error(0)
}

func TestUserProfileUsecase_Create(t *testing.T) {

	t.Run("Create user profile successfully", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserProfileRepo)
		uc := NewUserProfileUseCase(mockRepo)
		ctx := context.Background()

		userProfile := entity.UserProfile{
			UserID:       "U1234567890",
			DisplayName:  "test",
			PictureURL:   "https://example.com",
			AccessToken:  "test",
			RefreshToken: "test",
		}

		mockRepo.On("Create", ctx, userProfile).Return(nil)

		// Act
		result, err := uc.Create(ctx, userProfile)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, userProfile, result)
		mockRepo.AssertExpectations(t)
	})
	t.Run("Create user profile with invalid input", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserProfileRepo)
		uc := NewUserProfileUseCase(mockRepo)
		ctx := context.Background()

		userProfile := entity.UserProfile{}

		mockRepo.On("Create", ctx, userProfile).Return(assert.AnError)

		// Act
		result, err := uc.Create(ctx, userProfile)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, entity.UserProfile{}, result)
		mockRepo.AssertExpectations(t)
	})

}

func TestUserProfileUsecase_GetByID(t *testing.T) {
	t.Run("Get user profile by id successfully", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserProfileRepo)
		uc := NewUserProfileUseCase(mockRepo)
		ctx := context.Background()

		userProfile := entity.UserProfile{
			UserID:       "U1234567890",
			DisplayName:  "test",
			PictureURL:   "https://example.com",
			AccessToken:  "test",
			RefreshToken: "test",
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
		mockRepo := new(MockUserProfileRepo)
		uc := NewUserProfileUseCase(mockRepo)
		ctx := context.Background()

		mockRepo.On("GetByID", ctx, "123").Return(entity.UserProfile{}, assert.AnError)

		// Act
		result, err := uc.GetByID(ctx, "123")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, entity.UserProfile{}, result)
		mockRepo.AssertExpectations(t)
	})

}

func TestUserProfileUsecase_UpdateRefreshToken(t *testing.T) {
	t.Run("Update refresh token of user profile successfully", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserProfileRepo)
		uc := NewUserProfileUseCase(mockRepo)
		ctx := context.Background()

		userProfile := entity.UserProfile{
			UserID:       "U1234567890",
			DisplayName:  "test",
			PictureURL:   "https://example.com",
			AccessToken:  "test",
			RefreshToken: "test",
		}

		mockRepo.On("UpdateRefreshToken", ctx, userProfile.UserID, userProfile.RefreshToken).Return(nil)

		// Act
		err := uc.UpdateRefreshToken(ctx, userProfile.UserID, userProfile.RefreshToken)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
	t.Run("Update refresh token of user profile which does not existed", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserProfileRepo)
		uc := NewUserProfileUseCase(mockRepo)
		ctx := context.Background()

		userID := "U1234567890"
		refreshToken := "test"

		mockRepo.On("UpdateRefreshToken", ctx, userID, refreshToken).Return(assert.AnError)

		// Act
		err := uc.UpdateRefreshToken(ctx, userID, refreshToken)

		// Assert
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

}
