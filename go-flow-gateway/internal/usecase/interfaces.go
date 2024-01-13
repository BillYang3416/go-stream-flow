package usecase

import (
	"context"

	"github.com/bgg/go-flow-gateway/internal/entity"
)

type UserProfile interface {
	Create(context.Context, entity.UserProfile) (entity.UserProfile, error)
	GetByID(context.Context, string) (entity.UserProfile, error)
	UpdateRefreshToken(context.Context, string, string) error
}

type UserProfileRepo interface {
	Create(context.Context, entity.UserProfile) error
	GetByID(context.Context, string) (entity.UserProfile, error)
	UpdateRefreshToken(context.Context, string, string) error
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
