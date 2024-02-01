package usecase

import (
	"context"
	"testing"

	"github.com/bgg/go-flow-gateway/internal/entity"
	"github.com/bgg/go-flow-gateway/pkg/logger"
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

func (m *MockUserUploadedFileRepo) Create(ctx context.Context, u entity.UserUploadedFile) (int, error) {
	args := m.Called(ctx, u)
	return args.Int(0), args.Error(1)
}

func (m *MockUserUploadedFilePublisher) Publish(ctx context.Context, file entity.UserUploadedFile) error {
	args := m.Called(ctx, file)
	return args.Error(0)
}

func (m *MockUserUploadedFileEmailSender) Send(ctx context.Context, file entity.UserUploadedFile) error {
	args := m.Called(ctx, file)
	return args.Error(0)
}

func (m *MockUserUploadedFileRepo) GetPaginatedFiles(ctx context.Context, lastID, userID, limit int) ([]entity.UserUploadedFile, int, error) {
	args := m.Called(ctx, lastID, userID, limit)
	return args.Get(0).([]entity.UserUploadedFile), args.Get(1).(int), args.Error(2)
}

func (m *MockUserUploadedFileRepo) UpdateEmailSent(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupUserUploadedFileUseCase(t *testing.T) (*UserUploadedFileUseCase, *MockUserUploadedFileRepo, *MockUserUploadedFilePublisher, *MockUserUploadedFileEmailSender) {
	t.Helper()

	mockRepo := new(MockUserUploadedFileRepo)
	mockPub := new(MockUserUploadedFilePublisher)
	mockSender := new(MockUserUploadedFileEmailSender)
	uc := NewUserUploadedFileUseCase(mockRepo, mockPub, mockSender, logger.New("debug"))
	return uc, mockRepo, mockPub, mockSender
}

func TestUserUploadedFileUseCase_Create(t *testing.T) {

	const (
		ID      = 1
		name    = "test.txt"
		size    = 100
		content = "test"
		userID  = 123
	)
	t.Run("Create user uploaded file successfully", func(t *testing.T) {
		// Arrange
		uc, mockRepo, mockPub, _ := setupUserUploadedFileUseCase(t)
		ctx := context.Background()

		userUploadedFile := entity.UserUploadedFile{
			ID:      ID,
			Name:    name,
			Size:    size,
			Content: []byte(content),
			UserID:  userID,
		}

		mockRepo.On("Create", ctx, userUploadedFile).Return(ID, nil)
		mockPub.On("Publish", ctx, userUploadedFile).Return(nil)

		// Act
		result, err := uc.Create(ctx, userUploadedFile)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, userUploadedFile, result)
		mockRepo.AssertExpectations(t)
		mockPub.AssertExpectations(t)
	})

	t.Run("Create user uploaded file with empty file", func(t *testing.T) {
		// Arrange
		uc, mockRepo, _, _ := setupUserUploadedFileUseCase(t)
		ctx := context.Background()

		userUploadedFile := entity.UserUploadedFile{}
		mockRepo.On("Create", ctx, userUploadedFile).Return(0, assert.AnError)

		// Act
		result, err := uc.Create(ctx, userUploadedFile)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, entity.UserUploadedFile{}, result)
		mockRepo.AssertExpectations(t)
	})

}

func TestUserUploadedFileUseCase_SendEmail(t *testing.T) {
	const (
		ID      = 1
		name    = "test.txt"
		size    = 100
		content = "test"
		userID  = 123
	)

	t.Run("Send email successfully", func(t *testing.T) {
		// Arrange
		uc, mockRepo, _, mockSender := setupUserUploadedFileUseCase(t)
		ctx := context.Background()

		userUploadedFile := entity.UserUploadedFile{
			ID:      ID,
			Name:    name,
			Size:    size,
			Content: []byte(content),
			UserID:  userID,
		}

		mockRepo.On("UpdateEmailSent", ctx, userUploadedFile.ID).Return(nil)
		mockSender.On("Send", ctx, userUploadedFile).Return(nil)

		// Act
		err := uc.SendEmail(ctx, userUploadedFile)

		// Assert
		assert.NoError(t, err)
		mockSender.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Send email with empty email recipient", func(t *testing.T) {
		// Arrange
		uc, _, _, mockSender := setupUserUploadedFileUseCase(t)
		ctx := context.Background()

		userUploadedFile := entity.UserUploadedFile{
			ID:             ID,
			Name:           name,
			Size:           size,
			Content:        []byte(content),
			UserID:         userID,
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

func TestUserUploadedFileUseCase_GetPaginatedFiles(t *testing.T) {

	const (
		userID  = 123
		lastID  = 0
		limit   = 10
		ID      = 1
		name    = "test.txt"
		size    = 100
		content = "test"
	)
	t.Run("Get paginated files successfully", func(t *testing.T) {
		// Arrange
		uc, mockRepo, _, _ := setupUserUploadedFileUseCase(t)
		ctx := context.Background()

		userUploadedFiles := []entity.UserUploadedFile{
			{
				ID:      ID,
				Name:    name,
				Size:    size,
				Content: []byte(content),
				UserID:  userID,
			},
		}

		mockRepo.On("GetPaginatedFiles", ctx, lastID, userID, limit).Return(userUploadedFiles, len(userUploadedFiles), nil)

		// Act
		result, totalRecords, err := uc.GetPaginatedFiles(ctx, lastID, userID, limit)

		// Assert
		assert.Equal(t, len(userUploadedFiles), totalRecords)
		assert.NoError(t, err)
		assert.Equal(t, userUploadedFiles, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Get paginated files with invalid user ID", func(t *testing.T) {
		// Arrange
		uc, mockRepo, _, _ := setupUserUploadedFileUseCase(t)
		ctx := context.Background()

		mockRepo.On("GetPaginatedFiles", ctx, lastID, userID, limit).Return([]entity.UserUploadedFile{}, 0, assert.AnError)

		// Act
		_, _, err := uc.GetPaginatedFiles(ctx, lastID, userID, limit)

		// Assert
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
