package usecase

import (
	"context"
	"testing"

	"github.com/bgg/go-flow-gateway/internal/entity"
	"github.com/bgg/go-flow-gateway/internal/usecase/apperrors"
	"github.com/bgg/go-flow-gateway/internal/usecase/dto"
	"github.com/bgg/go-flow-gateway/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOAuthDetailRepo struct {
	mock.Mock
}

type MockTokenService struct {
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

func (m *MockTokenService) ExchangeCodeForTokens(code, domainUrl string) (*dto.TokenResponse, error) {
	args := m.Called(code, domainUrl)
	return args.Get(0).(*dto.TokenResponse), args.Error(1)
}

func (m *MockTokenService) VerifyIDToken(idToken, clientID string) (*dto.LineUserProfile, error) {
	args := m.Called(idToken, clientID)
	return args.Get(0).(*dto.LineUserProfile), args.Error(1)
}

func setupOAuthDetailUsecase(t *testing.T) (*OAuthDetailUseCase, *MockOAuthDetailRepo, *MockUserProfileUseCase, *MockTokenService) {
	t.Helper()
	mockRepo := new(MockOAuthDetailRepo)
	mockUserProfileUseCase := new(MockUserProfileUseCase)
	mockTokenService := new(MockTokenService)

	uc := NewOAuthDetailUseCase(mockRepo, mockUserProfileUseCase, logger.New("debug"), mockTokenService)
	return uc, mockRepo, mockUserProfileUseCase, mockTokenService
}

func TestOAuthDetailUsecase_HandleOAuthCallback(t *testing.T) {

	const (
		code         = "testcode"
		domainURL    = "https://example.com/callback"
		clientID     = "testclientid"
		userID       = 1
		provider     = "line"
		idToken      = "testidtoken"
		expiresIn    = 3600
		accessToken  = "testaccesstoken"
		refreshToken = "testrefreshtoken"
		sub          = "testsub"
		displayName  = "Test User"
		pictureURL   = "https://example.com/profile.jpg"
	)

	t.Run("Create oauth detail successfully", func(t *testing.T) {

		uc, mockRepo, mockUserProfileUseCase, mockTokenService := setupOAuthDetailUsecase(t)

		ctx := context.Background()

		oAuthDetail := entity.OAuthDetail{
			OAuthID:      sub,
			UserID:       userID,
			Provider:     provider,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		mockTokenService.On("ExchangeCodeForTokens", code, domainURL).Return(&dto.TokenResponse{
			IDToken:      idToken,
			ExpiresIn:    expiresIn,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}, nil)
		mockTokenService.On("VerifyIDToken", idToken, clientID).Return(&dto.LineUserProfile{
			Sub:     sub,
			Name:    displayName,
			Picture: pictureURL,
		}, nil)
		mockRepo.On("GetByOAuthID", ctx, sub).Return(entity.OAuthDetail{}, apperrors.NewNoRowsAffectedError("test", "test"))

		mockUserProfileUseCase.On("Create", ctx, entity.UserProfile{
			DisplayName: displayName,
			PictureURL:  pictureURL,
		}).Return(entity.UserProfile{
			UserID: userID,
		}, nil)

		mockRepo.On("Create", ctx, oAuthDetail).Return(nil)

		gotOAuthDetail, err := uc.HandleOAuthCallback(ctx, code, domainURL, provider, clientID)

		assert.NoError(t, err)
		assert.Equal(t, oAuthDetail, gotOAuthDetail)
		mockRepo.AssertExpectations(t)
		mockUserProfileUseCase.AssertExpectations(t)
		mockTokenService.AssertExpectations(t)
	})

	t.Run("Create oauth detail with invalid input", func(t *testing.T) {

		uc, mockRepo, mockUserProfileUseCase, mockTokenService := setupOAuthDetailUsecase(t)
		ctx := context.Background()

		oAuthDetail := entity.OAuthDetail{
			OAuthID:      sub,
			UserID:       userID,
			Provider:     provider,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		mockTokenService.On("ExchangeCodeForTokens", code, domainURL).Return(&dto.TokenResponse{
			IDToken:      idToken,
			ExpiresIn:    expiresIn,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}, nil)
		mockTokenService.On("VerifyIDToken", idToken, clientID).Return(&dto.LineUserProfile{
			Sub:     sub,
			Name:    displayName,
			Picture: pictureURL,
		}, nil)
		mockRepo.On("GetByOAuthID", ctx, sub).Return(entity.OAuthDetail{}, apperrors.NewNoRowsAffectedError("test", "test"))

		mockUserProfileUseCase.On("Create", ctx, entity.UserProfile{
			DisplayName: displayName,
			PictureURL:  pictureURL,
		}).Return(entity.UserProfile{
			UserID: userID,
		}, nil)

		mockRepo.On("Create", ctx, oAuthDetail).Return(assert.AnError)

		_, err := uc.HandleOAuthCallback(ctx, code, domainURL, provider, clientID)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
		mockUserProfileUseCase.AssertExpectations(t)
		mockTokenService.AssertExpectations(t)
	})
}

func TestOAuthDetailUsecase_UpdateRefreshToken(t *testing.T) {

	t.Run("Update refresh token successfully", func(t *testing.T) {

		uc, mockRepo, _, _ := setupOAuthDetailUsecase(t)
		ctx := context.Background()

		userId := "123"
		refreshToken := "123"

		mockRepo.On("UpdateRefreshToken", ctx, userId, refreshToken).Return(nil)

		err := uc.UpdateRefreshToken(ctx, userId, refreshToken)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Update refresh token with invalid input", func(t *testing.T) {

		uc, mockRepo, _, _ := setupOAuthDetailUsecase(t)
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

		uc, mockRepo, _, _ := setupOAuthDetailUsecase(t)
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

		uc, mockRepo, _, _ := setupOAuthDetailUsecase(t)
		ctx := context.Background()

		mockRepo.On("GetByOAuthID", ctx, "123").Return(entity.OAuthDetail{}, assert.AnError)

		result, err := uc.GetByOAuthID(ctx, "123")

		assert.Error(t, err)
		assert.Equal(t, entity.OAuthDetail{}, result)
		mockRepo.AssertExpectations(t)
	})

}
