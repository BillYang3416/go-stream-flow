package usecase

import (
	"context"
	"testing"

	"github.com/bgg/go-flow-gateway/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOAuthDetailRepo struct {
	mock.Mock
}

func (m *MockOAuthDetailRepo) Create(ctx context.Context, u entity.OAuthDetail) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockOAuthDetailRepo) UpdateRefreshToken(ctx context.Context, userId string, refreshToken string) error {
	args := m.Called(ctx, userId, refreshToken)
	return args.Error(0)
}

func (m *MockOAuthDetailRepo) GetByOAuthID(ctx context.Context, oauthId string) (entity.OAuthDetail, error) {
	args := m.Called(ctx, oauthId)
	return args.Get(0).(entity.OAuthDetail), args.Error(1)
}

func TestOAuthDetailUsecase_Create(t *testing.T) {

	t.Run("Create oauth detail successfully", func(t *testing.T) {

		mockRepo := new(MockOAuthDetailRepo)
		uc := NewOAuthDetailUseCase(mockRepo)
		ctx := context.Background()

		oauthDetail := entity.OAuthDetail{
			UserID:       123,
			Provider:     "line",
			AccessToken:  "123",
			RefreshToken: "123",
		}

		mockRepo.On("Create", ctx, oauthDetail).Return(nil)

		err := uc.Create(ctx, oauthDetail)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Create oauth detail with invalid input", func(t *testing.T) {

		mockRepo := new(MockOAuthDetailRepo)
		uc := NewOAuthDetailUseCase(mockRepo)
		ctx := context.Background()

		oauthDetail := entity.OAuthDetail{}

		mockRepo.On("Create", ctx, oauthDetail).Return(assert.AnError)

		err := uc.Create(ctx, oauthDetail)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestOAuthDetailUsecase_UpdateRefreshToken(t *testing.T) {

	t.Run("Update refresh token successfully", func(t *testing.T) {

		mockRepo := new(MockOAuthDetailRepo)
		uc := NewOAuthDetailUseCase(mockRepo)
		ctx := context.Background()

		userId := "123"
		refreshToken := "123"

		mockRepo.On("UpdateRefreshToken", ctx, userId, refreshToken).Return(nil)

		err := uc.UpdateRefreshToken(ctx, userId, refreshToken)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Update refresh token with invalid input", func(t *testing.T) {

		mockRepo := new(MockOAuthDetailRepo)
		uc := NewOAuthDetailUseCase(mockRepo)
		ctx := context.Background()

		userId := ""
		refreshToken := ""

		mockRepo.On("UpdateRefreshToken", ctx, userId, refreshToken).Return(assert.AnError)

		err := uc.UpdateRefreshToken(ctx, userId, refreshToken)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestOAuthDetailUsecase_GetByOAuthID(t *testing.T) {
	t.Run("Get oauth detail by id successfully", func(t *testing.T) {

		mockRepo := new(MockOAuthDetailRepo)
		uc := NewOAuthDetailUseCase(mockRepo)
		ctx := context.Background()

		oauthDetail := entity.OAuthDetail{
			OAuthID:      "123",
			UserID:       123,
			Provider:     "line",
			AccessToken:  "123",
			RefreshToken: "123",
		}

		mockRepo.On("GetByOAuthID", ctx, oauthDetail.OAuthID).Return(oauthDetail, nil)

		result, err := uc.GetByOAuthID(ctx, oauthDetail.OAuthID)

		assert.NoError(t, err)
		assert.Equal(t, oauthDetail, result)
		mockRepo.AssertExpectations(t)
	})
	t.Run("Get oauth detail by id with invalid oauth ID", func(t *testing.T) {

		mockRepo := new(MockOAuthDetailRepo)
		uc := NewOAuthDetailUseCase(mockRepo)
		ctx := context.Background()

		mockRepo.On("GetByOAuthID", ctx, "123").Return(entity.OAuthDetail{}, assert.AnError)

		result, err := uc.GetByOAuthID(ctx, "123")

		assert.Error(t, err)
		assert.Equal(t, entity.OAuthDetail{}, result)
		mockRepo.AssertExpectations(t)
	})

}
