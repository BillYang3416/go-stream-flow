package usecase

import (
	"context"

	"github.com/bgg/go-flow-gateway/internal/entity"
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
	Create(context.Context, entity.UserUploadedFile) (entity.UserUploadedFile, error)
	SendEmail(context.Context, entity.UserUploadedFile) error
}

type UserUploadedFileRepo interface {
	Create(context.Context, entity.UserUploadedFile) error
}

type UserUploadedFilePublisher interface {
	Publish(context.Context, entity.UserUploadedFile) error
}

type UserUploadedFileEmailSender interface {
	Send(context.Context, entity.UserUploadedFile) error
}

type OAuthDetail interface {
	Create(context.Context, entity.OAuthDetail) error
	UpdateRefreshToken(context.Context, string, string) error
	GetByOAuthID(context.Context, string) (entity.OAuthDetail, error)
}

type OAuthDetailRepo interface {
	Create(context.Context, entity.OAuthDetail) error
	UpdateRefreshToken(context.Context, string, string) error
	GetByOAuthID(context.Context, string) (entity.OAuthDetail, error)
}

type UserCredential interface {
	Create(context.Context, int, string, string) error
	GetByUsername(context.Context, string) (entity.UserCredential, error)
}

type UserCredentialRepo interface {
	Create(context.Context, entity.UserCredential) error
	GetByUsername(context.Context, string) (entity.UserCredential, error)
}

type PasswordHasher interface {
	GenerateHash(context.Context, string) (string, error)
	CompareHash(context.Context, string, string) error
}
