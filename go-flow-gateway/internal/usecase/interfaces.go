package usecase

import (
	"context"

	"github.com/bgg/go-flow-gateway/internal/entity"
	"github.com/bgg/go-flow-gateway/internal/usecase/dto"
)

type UserProfile interface {
	Create(context.Context, entity.UserProfile) (entity.UserProfile, error)
	GetByID(context.Context, int) (entity.UserProfile, error)
}

type UserProfileRepo interface {
	Create(context.Context, entity.UserProfile) (entity.UserProfile, error)
	GetByID(context.Context, int) (entity.UserProfile, error)
}

type UserUploadedFile interface {
	Create(ctx context.Context, userUploadedFile entity.UserUploadedFile) (entity.UserUploadedFile, error)
	SendEmail(ctx context.Context, userUploadedFile entity.UserUploadedFile) error
	GetPaginatedFiles(ctx context.Context, lastID, userID, limit int) ([]entity.UserUploadedFile, int, error)
}

type UserUploadedFileRepo interface {
	Create(ctx context.Context, userUploadedFile entity.UserUploadedFile) (int, error)
	GetPaginatedFiles(ctx context.Context, lastID, userID, limit int) ([]entity.UserUploadedFile, int, error)
	UpdateEmailSent(ctx context.Context, userUploadedFileID int) error
}

type UserUploadedFilePublisher interface {
	Publish(ctx context.Context, userUploadedFile entity.UserUploadedFile) error
}

type UserUploadedFileEmailSender interface {
	Send(ctx context.Context, userUploadedFile entity.UserUploadedFile) error
}

type OAuthDetail interface {
	HandleOAuthCallback(ctx context.Context, code, domainUrl, provider, clientID string) (entity.OAuthDetail, error)
	UpdateRefreshToken(ctx context.Context, userID, refreshToken string) error
	GetByOAuthID(ctx context.Context, oAuthID string) (entity.OAuthDetail, error)
}

type OAuthDetailRepo interface {
	Create(ctx context.Context, oAuthDetail entity.OAuthDetail) error
	UpdateRefreshToken(ctx context.Context, oAuthID, refreshToken string) error
	GetByOAuthID(ctx context.Context, oAuthID string) (entity.OAuthDetail, error)
}

type TokenService interface {
	ExchangeCodeForTokens(code, domainUrl string) (*dto.TokenResponse, error)
	VerifyIDToken(idToken, clietID string) (*dto.LineUserProfile, error)
}

type UserCredential interface {
	Register(ctx context.Context, displayName, username, password string) (int, error)
	GetByUsername(ctx context.Context, username string) (entity.UserCredential, error)
	Login(ctx context.Context, username, password string) (entity.UserCredential, error)
}

type UserCredentialRepo interface {
	Create(ctx context.Context, userCredential entity.UserCredential) error
	GetByUsername(ctx context.Context, username string) (entity.UserCredential, error)
}

type PasswordHasher interface {
	GenerateHash(ctx context.Context, password string) (string, error)
	CompareHash(ctx context.Context, password, hashedPassword string) error
}
