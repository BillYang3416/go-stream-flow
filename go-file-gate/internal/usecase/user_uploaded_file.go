package usecase

import (
	"context"
	"fmt"

	"github.com/bgg/go-file-gate/internal/entity"
)

type UserUploadedFileUseCase struct {
	repo UserUploadedFileRepo
}

func NewUserUploadedFileUseCase(r UserUploadedFileRepo) *UserUploadedFileUseCase {
	return &UserUploadedFileUseCase{repo: r}
}

func (uc *UserUploadedFileUseCase) Create(ctx context.Context, userUploadedFile entity.UserUploadedFile) (entity.UserUploadedFile, error) {
	err := uc.repo.Create(ctx, userUploadedFile)
	if err != nil {
		return entity.UserUploadedFile{}, fmt.Errorf("UserUploadedFileUseCase - Create - s.repo.Create: %w", err)
	}
	return userUploadedFile, nil
}
