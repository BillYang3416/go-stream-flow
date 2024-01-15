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
	_, err := uc.repo.GetByUsername(ctx, username)
	if err == nil {
		return fmt.Errorf("UserCredentialUseCase - Create - GetByUsername: has duplicate username")
	}

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

func (uc *UserCredentialUseCase) GetByUsername(ctx context.Context, username string) (entity.UserCredential, error) {
	u, err := uc.repo.GetByUsername(ctx, username)
	if err != nil {
		return entity.UserCredential{}, fmt.Errorf("UserCredentialUseCase - GetByUsername: %w", err)
	}

	return u, nil
}

func (uc *UserCredentialUseCase) Login(ctx context.Context, username, password string) (entity.UserCredential, error) {
	u, err := uc.repo.GetByUsername(ctx, username)
	if err != nil {
		return entity.UserCredential{}, fmt.Errorf("UserCredentialUseCase - Login - GetByUsername: %w", err)
	}

	err = uc.hasher.CompareHash(ctx, password, u.PasswordHash)
	if err != nil {
		return entity.UserCredential{}, fmt.Errorf("UserCredentialUseCase - Login - CompareHash: %w", err)
	}

	return u, nil
}
