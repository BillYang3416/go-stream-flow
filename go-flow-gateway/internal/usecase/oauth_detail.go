package usecase

import (
	"context"
	"fmt"

	"github.com/bgg/go-flow-gateway/internal/entity"
)

type OAuthDetailUseCase struct {
	repo OAuthDetailRepo
}

func NewOAuthDetailUseCase(r OAuthDetailRepo) *OAuthDetailUseCase {
	return &OAuthDetailUseCase{repo: r}
}

func (uc *OAuthDetailUseCase) Create(ctx context.Context, oauthDetail entity.OAuthDetail) error {
	err := uc.repo.Create(ctx, oauthDetail)
	if err != nil {
		return fmt.Errorf("OAuthDetailUseCase - Create - s.repo.Create: %w", err)
	}
	return nil
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
