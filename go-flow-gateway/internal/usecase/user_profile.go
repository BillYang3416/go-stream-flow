package usecase

import (
	"context"
	"fmt"

	"github.com/bgg/go-flow-gateway/internal/entity"
	"github.com/bgg/go-flow-gateway/pkg/logger"
)

type UserProfileUseCase struct {
	repo   UserProfileRepo
	logger logger.Logger
}

func NewUserProfileUseCase(r UserProfileRepo, l logger.Logger) *UserProfileUseCase {
	return &UserProfileUseCase{repo: r, logger: l}
}

func (uc *UserProfileUseCase) Create(ctx context.Context, userProfile entity.UserProfile) (entity.UserProfile, error) {
	userProfile, err := uc.repo.Create(ctx, userProfile)
	if err != nil {
		uc.logger.Error("UserProfileUseCase - Create - repo.Create: failed to create user profile", err)
		return entity.UserProfile{}, fmt.Errorf("UserProfileUseCase - Create - s.repo.Create: %w", err)
	}

	uc.logger.Info("UserProfileUseCase - Create - repo.Create: successfully created user profile", "userID", userProfile.UserID)
	return userProfile, nil
}

func (uc *UserProfileUseCase) GetByID(ctx context.Context, userID int) (entity.UserProfile, error) {
	userProfile, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		uc.logger.Error("UserProfileUseCase - Get - repo.Find: failed to get user profile", err)
		return entity.UserProfile{}, fmt.Errorf("UserProfileUseCase - Get - s.repo.Find: %w", err)
	}

	uc.logger.Info("UserProfileUseCase - Get - repo.Find: successfully got user profile", "userID", userProfile.UserID)
	return userProfile, nil
}
