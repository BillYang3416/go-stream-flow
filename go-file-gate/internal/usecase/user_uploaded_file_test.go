package usecase

import (
	"context"
	"testing"

	"github.com/bgg/go-file-gate/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserUploadedFilePublisher struct {
	mock.Mock
}

type MockUserUploadedFileRepo struct {
	mock.Mock
}

func (m *MockUserUploadedFileRepo) Create(ctx context.Context, u entity.UserUploadedFile) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserUploadedFilePublisher) Publish(ctx context.Context, file entity.UserUploadedFile) error {
	args := m.Called(ctx, file)
	return args.Error(0)
}

func TestUserUploadedFileUseCase_Create(t *testing.T) {
	t.Run("Create user uploaded file successfully", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserUploadedFileRepo)
		mockPub := new(MockUserUploadedFilePublisher)
		uc := NewUserUploadedFileUseCase(mockRepo, mockPub)
		ctx := context.Background()
		userUploadedFile := entity.UserUploadedFile{
			Name:    "test.txt",
			Size:    100,
			Content: []byte("test"),
			UserID:  "123",
		}
		mockRepo.On("Create", ctx, userUploadedFile).Return(nil)
		mockPub.On("Publish", ctx, userUploadedFile).Return(nil)

		// Act
		result, err := uc.Create(ctx, userUploadedFile)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, userUploadedFile, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Create user uploaded file with empty file", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserUploadedFileRepo)
		mockPub := new(MockUserUploadedFilePublisher)
		uc := NewUserUploadedFileUseCase(mockRepo, mockPub)
		ctx := context.Background()
		userUploadedFile := entity.UserUploadedFile{
			UserID: "123",
		}
		mockRepo.On("Create", ctx, userUploadedFile).Return(assert.AnError)

		// Act
		result, err := uc.Create(ctx, userUploadedFile)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, entity.UserUploadedFile{}, result)
		mockRepo.AssertExpectations(t)
	})

}
