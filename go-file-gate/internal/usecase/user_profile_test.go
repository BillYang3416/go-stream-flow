package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/bgg/go-file-gate/internal/entity"
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

func TestCreateUserProfile(t *testing.T) {

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
}

func TestCreateUserProfileWithError(t *testing.T) {
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

	expectedError := errors.New("failed to create user profile")
	mockRepo.On("Create", ctx, userProfile).Return(expectedError)

	// Act
	result, err := uc.Create(ctx, userProfile)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, entity.UserProfile{}, result)
	mockRepo.AssertExpectations(t)
}

func TestGetUserProfileByID(t *testing.T) {
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
}

func TestGetUserProfileByIDWithError(t *testing.T) {
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

	expectedError := errors.New("failed to get user profile")
	mockRepo.On("GetByID", ctx, userProfile.UserID).Return(entity.UserProfile{}, expectedError)

	// Act
	result, err := uc.GetByID(ctx, userProfile.UserID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, entity.UserProfile{}, result)
	mockRepo.AssertExpectations(t)
}

func TestUpdateRefreshTokenOfUserProfile(t *testing.T) {
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
}

func TestUpdateRefreshTokenOfUserProfileWithError(t *testing.T) {
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

	expectedError := errors.New("failed to update refresh token")
	mockRepo.On("UpdateRefreshToken", ctx, userProfile.UserID, userProfile.RefreshToken).Return(expectedError)

	// Act
	err := uc.UpdateRefreshToken(ctx, userProfile.UserID, userProfile.RefreshToken)

	// Assert
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
