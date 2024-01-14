package usecase

import (
	"context"
	"testing"

	"github.com/bgg/go-flow-gateway/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserUploadedFilePublisher struct {
	mock.Mock
}

type MockUserUploadedFileRepo struct {
	mock.Mock
}

type MockUserUploadedFileEmailSender struct {
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

func (m *MockUserUploadedFileEmailSender) Send(ctx context.Context, file entity.UserUploadedFile) error {
	args := m.Called(ctx, file)
	return args.Error(0)
}

func TestUserUploadedFileUseCase_Create(t *testing.T) {
	t.Run("Create user uploaded file successfully", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserUploadedFileRepo)
		mockPub := new(MockUserUploadedFilePublisher)
		mockSender := new(MockUserUploadedFileEmailSender)

		uc := NewUserUploadedFileUseCase(mockRepo, mockPub, mockSender)
		ctx := context.Background()
		userUploadedFile := entity.UserUploadedFile{
			Name:    "test.txt",
			Size:    100,
			Content: []byte("test"),
			UserID:  123,
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
		mockSender := new(MockUserUploadedFileEmailSender)

		uc := NewUserUploadedFileUseCase(mockRepo, mockPub, mockSender)
		ctx := context.Background()
		userUploadedFile := entity.UserUploadedFile{
			UserID: 123,
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

func TestUserUploadedFileUseCase_SendEmail(t *testing.T) {
	t.Run("Send email successfully", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserUploadedFileRepo)
		mockPub := new(MockUserUploadedFilePublisher)
		mockSender := new(MockUserUploadedFileEmailSender)
		uc := NewUserUploadedFileUseCase(mockRepo, mockPub, mockSender)
		ctx := context.Background()
		userUploadedFile := entity.UserUploadedFile{
			Name:           "test.txt",
			Size:           100,
			Content:        []byte("test"),
			UserID:         123,
			EmailRecipient: "",
		}
		mockSender.On("Send", ctx, userUploadedFile).Return(nil)

		// Act
		err := uc.SendEmail(ctx, userUploadedFile)

		// Assert
		assert.NoError(t, err)
		mockSender.AssertExpectations(t)
	})

	t.Run("Send email with empty email recipient", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserUploadedFileRepo)
		mockPub := new(MockUserUploadedFilePublisher)
		mockSender := new(MockUserUploadedFileEmailSender)
		uc := NewUserUploadedFileUseCase(mockRepo, mockPub, mockSender)
		ctx := context.Background()
		userUploadedFile := entity.UserUploadedFile{
			Name:           "test.txt",
			Size:           100,
			Content:        []byte("test"),
			UserID:         123,
			EmailRecipient: "",
		}
		mockSender.On("Send", ctx, userUploadedFile).Return(assert.AnError)

		// Act
		err := uc.SendEmail(ctx, userUploadedFile)

		// Assert
		assert.Error(t, err)
		mockSender.AssertExpectations(t)
	})
}
