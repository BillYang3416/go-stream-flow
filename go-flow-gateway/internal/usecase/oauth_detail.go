package usecase

import (
	"context"
	"fmt"

	"github.com/bgg/go-flow-gateway/internal/entity"
	"github.com/bgg/go-flow-gateway/internal/usecase/apperrors"
	"github.com/bgg/go-flow-gateway/pkg/logger"
)

type OAuthDetailUseCase struct {
	repo               OAuthDetailRepo
	userProfileUseCase UserProfile
	logger             logger.Logger
	tokenSvc           TokenService
}

func NewOAuthDetailUseCase(r OAuthDetailRepo, u UserProfile, l logger.Logger, t TokenService) *OAuthDetailUseCase {
	return &OAuthDetailUseCase{repo: r, userProfileUseCase: u, logger: l, tokenSvc: t}
}

func (uc *OAuthDetailUseCase) HandleOAuthCallback(ctx context.Context, code, domainUrl, provider, clientID string) (entity.OAuthDetail, error) {

	// Exchange code for tokens
	tokenResponse, err := uc.tokenSvc.ExchangeCodeForTokens(code, domainUrl)
	if err != nil {
		uc.logger.Error("OAuthDetailUseCase - HandleOAuthCallback - s.tokenSvc.ExchangeCodeForTokens: failed to exchange code for tokens", "error", err)
		return entity.OAuthDetail{}, fmt.Errorf("OAuthDetailUseCase - HandleOAuthCallback - s.tokenSvc.ExchangeCodeForTokens: %w", err)
	}

	// Verify ID Token And Get User Profile of Line
	lineUserProfile, err := uc.tokenSvc.VerifyIDToken(tokenResponse.IDToken, clientID)
	if err != nil {
		uc.logger.Error("OAuthDetailUseCase - HandleOAuthCallback - s.tokenSvc.VerifyIDToken: failed to verify id token", "error", err)
		return entity.OAuthDetail{}, fmt.Errorf("OAuthDetailUseCase - HandleOAuthCallback - s.tokenSvc.VerifyIDToken: %w", err)
	}

	// check the oauth detail is exists
	oauthDetail, err := uc.GetByOAuthID(ctx, lineUserProfile.Sub)
	if err != nil {
		if _, ok := apperrors.AsNoRowsAffectedError(err); ok {

			// create user profile
			userProfile, err := uc.userProfileUseCase.Create(ctx, entity.UserProfile{
				DisplayName: lineUserProfile.Name,
				PictureURL:  lineUserProfile.Picture,
			})
			if err != nil {
				uc.logger.Error("OAuthDetailUseCase - HandleOAuthCallback - s.userProfileUseCase.Create: failed to create user profile", "error", err)
				return entity.OAuthDetail{}, fmt.Errorf("OAuthDetailUseCase - HandleOAuthCallback - s.userProfileUseCase.Create: %w", err)
			}

			// create oauth detail
			oauthDetail = entity.OAuthDetail{
				OAuthID:      lineUserProfile.Sub,
				UserID:       userProfile.UserID,
				Provider:     provider,
				AccessToken:  tokenResponse.AccessToken,
				RefreshToken: tokenResponse.RefreshToken,
			}

			err = uc.repo.Create(ctx, oauthDetail)
			if err != nil {
				uc.logger.Error("OAuthDetailUseCase - Create - s.repo.Create: failed to create oauth detail", "error", err)
				return entity.OAuthDetail{}, fmt.Errorf("OAuthDetailUseCase - Create - s.repo.Create: %w", err)
			}

		} else {
			uc.logger.Error("OAuthDetailUseCase - HandleOAuthCallback - s.GetByOAuthID: failed to get oauth detail by oauth id", "error", err)
			return entity.OAuthDetail{}, fmt.Errorf("OAuthDetailUseCase - HandleOAuthCallback - s.GetByOAuthID: %w", err)
		}
	} else {
		err = uc.UpdateRefreshToken(ctx, fmt.Sprint(oauthDetail.UserID), tokenResponse.RefreshToken)
		if err != nil {
			uc.logger.Error("OAuthDetailUseCase - HandleOAuthCallback - s.UpdateRefreshToken: failed to update refresh token", "error", err)
			return entity.OAuthDetail{}, fmt.Errorf("OAuthDetailUseCase - HandleOAuthCallback - s.UpdateRefreshToken: %w", err)
		}
	}

	uc.logger.Info("OAuthDetailUseCase - HandleOAuthCallback - success", "oauthDetail.UserID", oauthDetail.UserID)
	return oauthDetail, nil
}

func (uc *OAuthDetailUseCase) UpdateRefreshToken(ctx context.Context, userId string, refreshToken string) error {
	err := uc.repo.UpdateRefreshToken(ctx, userId, refreshToken)
	if err != nil {
		return fmt.Errorf("OAuthDetailUseCase - UpdateRefreshToken - s.repo.UpdateRefreshToken: %w", err)
	}
	return nil
}

func (uc *OAuthDetailUseCase) GetByOAuthID(ctx context.Context, oauthId string) (entity.OAuthDetail, error) {
	oauthDetail, err := uc.repo.GetByOAuthID(ctx, oauthId)
	if err != nil {
		return entity.OAuthDetail{}, fmt.Errorf("OAuthDetailUseCase - GetByOAuthID - s.repo.GetByOAuthID: %w", err)
	}
	return oauthDetail, nil
}
