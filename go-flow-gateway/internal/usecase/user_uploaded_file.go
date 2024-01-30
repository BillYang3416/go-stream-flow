package usecase

import (
	"context"
	"fmt"

	"github.com/bgg/go-flow-gateway/internal/entity"
	"github.com/bgg/go-flow-gateway/pkg/logger"
)

type UserUploadedFileUseCase struct {
	repo   UserUploadedFileRepo
	pub    UserUploadedFilePublisher
	sender UserUploadedFileEmailSender
	logger logger.Logger
}

func NewUserUploadedFileUseCase(r UserUploadedFileRepo, p UserUploadedFilePublisher, s UserUploadedFileEmailSender, l logger.Logger) *UserUploadedFileUseCase {
	return &UserUploadedFileUseCase{repo: r, pub: p, sender: s, logger: l}
}

func (uc *UserUploadedFileUseCase) Create(ctx context.Context, userUploadedFile entity.UserUploadedFile) (entity.UserUploadedFile, error) {
	returnedID, err := uc.repo.Create(ctx, userUploadedFile)
	if err != nil {
		uc.logger.Error("UserUploadedFileUseCase - Create - repo.Create : error creating user uploaded file", "error", err)
		return entity.UserUploadedFile{}, fmt.Errorf("UserUploadedFileUseCase - Create - s.repo.Create: %w", err)
	}
	uc.logger.Info("UserUploadedFileUseCase - Create : user uploaded file created", "userUploadedFileID", returnedID)

	userUploadedFile.ID = returnedID
	err = uc.pub.Publish(ctx, userUploadedFile)
	if err != nil {
		uc.logger.Error("UserUploadedFileUseCase - Create - pub.Publish : error publishing user uploaded file", "error", err)
		return entity.UserUploadedFile{}, fmt.Errorf("UserUploadedFileUseCase - Create - s.pub.Publish: %w", err)
	}

	uc.logger.Info("UserUploadedFileUseCase - Create : user uploaded file published", "userUploadedFileID", userUploadedFile.ID)
	return userUploadedFile, nil
}

func (uc *UserUploadedFileUseCase) SendEmail(ctx context.Context, userUploadedFile entity.UserUploadedFile) error {
	err := uc.sender.Send(ctx, userUploadedFile)
	if err != nil {
		uc.logger.Error("UserUploadedFileUseCase - SendEmail - sender.Send : error sending email", "error", err)
		return fmt.Errorf("UserUploadedFileUseCase - SendEmail - s.sender.Send: %w", err)
	}
	uc.logger.Info("UserUploadedFileUseCase - SendEmail : email sent", "userUploadedFileID", userUploadedFile.ID)

	err = uc.repo.UpdateEmailSent(ctx, userUploadedFile.ID)
	if err != nil {
		uc.logger.Error("UserUploadedFileUseCase - SendEmail - repo.UpdateEmailSent : error updating email sent", "error", err)
		return fmt.Errorf("UserUploadedFileUseCase - SendEmail - s.repo.UpdateEmailSent: %w", err)
	}

	uc.logger.Info("UserUploadedFileUseCase - SendEmail : email sent updated", "userUploadedFileID", userUploadedFile.ID)
	return nil
}

func (uc *UserUploadedFileUseCase) GetPaginatedFiles(ctx context.Context, lastID, userID, limit int) ([]entity.UserUploadedFile, int, error) {
	files, totalRecords, err := uc.repo.GetPaginatedFiles(ctx, lastID, userID, limit)
	if err != nil {
		uc.logger.Error("UserUploadedFileUseCase - GetPaginatedFiles - repo.GetPaginatedFiles : error getting paginated files", "error", err)
		return nil, 0, fmt.Errorf("UserUploadedFileUseCase - GetPaginatedFiles - s.repo.GetPaginatedFiles: %w", err)
	}

	uc.logger.Info("UserUploadedFileUseCase - GetPaginatedFiles : paginated files retrieved", "totalRecords", totalRecords)
	return files, totalRecords, nil
}
