package usecase

import (
	"context"
	"fmt"

	"github.com/bgg/go-flow-gateway/internal/entity"
)

type UserProfileUseCase struct {
	repo UserProfileRepo
}

func NewUserProfileUseCase(r UserProfileRepo) *UserProfileUseCase {
	return &UserProfileUseCase{repo: r}
}

func (uc *UserProfileUseCase) Create(ctx context.Context, userProfile entity.UserProfile) (entity.UserProfile, error) {
	userProfile, err := uc.repo.Create(ctx, userProfile)
	if err != nil {
		return entity.UserProfile{}, fmt.Errorf("UserProfileUseCase - Create - s.repo.Create: %w", err)
	}
	return userProfile, nil
}

func (uc *UserProfileUseCase) GetByID(ctx context.Context, userId int) (entity.UserProfile, error) {
	userProfile, err := uc.repo.GetByID(ctx, userId)
	if err != nil {
		return entity.UserProfile{}, fmt.Errorf("UserProfileUseCase - Get - s.repo.Find: %w", err)
	}
	return userProfile, nil
}
