package usecase

import (
	"context"
	"fmt"

	"github.com/bgg/go-flow-gateway/internal/entity"
)

type UserUploadedFileUseCase struct {
	repo   UserUploadedFileRepo
	pub    UserUploadedFilePublisher
	sender UserUploadedFileEmailSender
}

func NewUserUploadedFileUseCase(r UserUploadedFileRepo, p UserUploadedFilePublisher, s UserUploadedFileEmailSender) *UserUploadedFileUseCase {
	return &UserUploadedFileUseCase{repo: r, pub: p, sender: s}
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

func (uc *UserUploadedFileUseCase) SendEmail(ctx context.Context, userUploadedFile entity.UserUploadedFile) error {
	err := uc.sender.Send(ctx, userUploadedFile)
	if err != nil {
		return fmt.Errorf("UserUploadedFileUseCase - SendEmail - s.sender.Send: %w", err)
	}

	return nil
}

func (uc *UserUploadedFileUseCase) GetPaginatedFiles(ctx context.Context, lastID, userID, limit int) ([]entity.UserUploadedFile, int, error) {
	files, totalRecords, err := uc.repo.GetPaginatedFiles(ctx, lastID, userID, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("UserUploadedFileUseCase - GetPaginatedFiles - s.repo.GetPaginatedFiles: %w", err)
	}

	return files, totalRecords, nil
}
