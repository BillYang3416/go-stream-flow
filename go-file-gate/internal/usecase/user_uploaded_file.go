package usecase

import (
	"context"
	"fmt"

	"github.com/bgg/go-file-gate/internal/entity"
)

type UserUploadedFileUseCase struct {
	repo UserUploadedFileRepo
	pub  UserUploadedFilePublisher
}

func NewUserUploadedFileUseCase(r UserUploadedFileRepo, p UserUploadedFilePublisher) *UserUploadedFileUseCase {
	return &UserUploadedFileUseCase{repo: r, pub: p}
}

func (uc *UserUploadedFileUseCase) Create(ctx context.Context, userUploadedFile entity.UserUploadedFile) (entity.UserUploadedFile, error) {
	err := uc.repo.Create(ctx, userUploadedFile)
	if err != nil {
		return entity.UserUploadedFile{}, fmt.Errorf("UserUploadedFileUseCase - Create - s.repo.Create: %w", err)
	}

	err = uc.pub.Publish(ctx, userUploadedFile)
	if err != nil {
		return entity.UserUploadedFile{}, fmt.Errorf("UserUploadedFileUseCase - Create - s.pub.Publish: %w", err)
	}

	return userUploadedFile, nil
}
