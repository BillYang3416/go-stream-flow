package usecase

import (
	"context"
	"fmt"

	"github.com/bgg/go-flow-gateway/internal/entity"
)

type UserCredentialUseCase struct {
	repo   UserCredentialRepo
	hasher PasswordHasher
}

func NewUserCredentialUseCase(repo UserCredentialRepo, hasher PasswordHasher) *UserCredentialUseCase {
	return &UserCredentialUseCase{
		repo:   repo,
		hasher: hasher,
	}
}

func (uc *UserCredentialUseCase) Create(ctx context.Context, userID int, username, password string) error {
	hashedPassword, err := uc.hasher.GenerateHash(ctx, password)
	if err != nil {
		return fmt.Errorf("UserCredentialUseCase - Create - hasher.GenerateHash: %w", err)
	}

	u := entity.UserCredential{
		UserID:       userID,
		Username:     username,
		PasswordHash: hashedPassword,
	}

	err = uc.repo.Create(ctx, u)
	if err != nil {
		return fmt.Errorf("UserCredentialUseCase - Create: %w", err)
	}

	return nil
}
