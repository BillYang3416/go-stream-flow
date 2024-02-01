package usecase

import (
	"context"
	"fmt"

	"github.com/bgg/go-flow-gateway/internal/entity"
	"github.com/bgg/go-flow-gateway/internal/usecase/apperrors"
	"github.com/bgg/go-flow-gateway/pkg/logger"
)

type UserCredentialUseCase struct {
	repo        UserCredentialRepo
	hasher      PasswordHasher
	userProfile UserProfile
	logger      logger.Logger
}

func NewUserCredentialUseCase(repo UserCredentialRepo, hasher PasswordHasher, userProfile UserProfile, logger logger.Logger) *UserCredentialUseCase {
	return &UserCredentialUseCase{
		repo:        repo,
		hasher:      hasher,
		userProfile: userProfile,
		logger:      logger,
	}
}

func (uc *UserCredentialUseCase) Register(ctx context.Context, displayName, username, password string) (int, error) {
	_, err := uc.repo.GetByUsername(ctx, username)
	if err == nil {
		uc.logger.Warn("UserCredentialUseCase - Register: duplicate username found", "username", username)
		return 0, fmt.Errorf("UserCredentialUseCase - Register - GetByUsername: has duplicate username")
	} else if err != nil && !apperrors.IsNoRowsAffectedError(err) {
		uc.logger.Error("UserCredentialUseCase - Register - GetByUsername", "error", err)
		return 0, fmt.Errorf("UserCredentialUseCase - Register - GetByUsername: %w", err)
	}

	up, err := uc.userProfile.Create(ctx, entity.UserProfile{
		DisplayName: displayName,
	})
	if err != nil {
		uc.logger.Error("UserCredentialUseCase - Register - userProfile.Create : error creating user profile", "error", err)
		return 0, fmt.Errorf("UserCredentialUseCase - Register - userProfile.Create: %w", err)
	}
	uc.logger.Info("UserCredentialUseCase - Register - userProfile.Create: user profile created", "userID", up.UserID)

	hashedPassword, err := uc.hasher.GenerateHash(ctx, password)
	if err != nil {
		uc.logger.Error("UserCredentialUseCase - Register - hasher.GenerateHash: error generating hash", "error", err)
		return 0, fmt.Errorf("UserCredentialUseCase - Register - hasher.GenerateHash: %w", err)
	}

	u := entity.UserCredential{
		UserID:       up.UserID,
		Username:     username,
		PasswordHash: hashedPassword,
	}
	err = uc.repo.Create(ctx, u)
	if err != nil {
		uc.logger.Error("UserCredentialUseCase - Register - repo.Create: error creating user credential", "error", err)
		return 0, fmt.Errorf("UserCredentialUseCase - Register - repo.Create: %w", err)
	}

	uc.logger.Info("UserCredentialUseCase - Register: user registered", "userID", up.UserID)
	return up.UserID, nil
}

func (uc *UserCredentialUseCase) GetByUsername(ctx context.Context, username string) (entity.UserCredential, error) {
	u, err := uc.repo.GetByUsername(ctx, username)
	if err != nil {
		uc.logger.Error("UserCredentialUseCase - GetByUsername - repo.GetByUsername: failed to get user credential", "error", err)
		return entity.UserCredential{}, fmt.Errorf("UserCredentialUseCase - GetByUsername: %w", err)
	}

	uc.logger.Info("UserCredentialUseCase - GetByUsername: user credential retrieved", "userID", u.UserID)
	return u, nil
}

func (uc *UserCredentialUseCase) Login(ctx context.Context, username, password string) (entity.UserCredential, error) {
	u, err := uc.GetByUsername(ctx, username)
	if err != nil {
		return entity.UserCredential{}, fmt.Errorf("UserCredentialUseCase - Login - GetByUsername: %w", err)
	}

	err = uc.hasher.CompareHash(ctx, password, u.PasswordHash)
	if err != nil {
		uc.logger.Error("UserCredentialUseCase - Login - CompareHash: failed to compare hash", "error", err)
		return entity.UserCredential{}, fmt.Errorf("UserCredentialUseCase - Login - CompareHash: %w", err)
	}

	uc.logger.Info("UserCredentialUseCase - Login: user logged in", "userID", u.UserID)
	return u, nil
}
