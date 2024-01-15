package usecase

import (
	"context"
	"testing"

	"github.com/bgg/go-flow-gateway/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserCredentialRepo struct {
	mock.Mock
}

type MockPasswordHasher struct {
	mock.Mock
}

func (m *MockUserCredentialRepo) Create(ctx context.Context, u entity.UserCredential) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockPasswordHasher) GenerateHash(ctx context.Context, password string) (string, error) {
	args := m.Called(ctx, password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordHasher) CompareHash(ctx context.Context, password, hashedPassword string) error {
	args := m.Called(ctx, password, hashedPassword)
	return args.Error(0)
}

func TestUserCredentialUsecase_Create(t *testing.T) {

	t.Run("Create user credential successfully", func(t *testing.T) {

		mockRepo := new(MockUserCredentialRepo)
		mockHasher := new(MockPasswordHasher)
		uc := NewUserCredentialUseCase(mockRepo, mockHasher)
		ctx := context.Background()

		password := "test"
		hashedPassword := "$2a$10"
		mockHasher.On("GenerateHash", ctx, password).Return(hashedPassword, nil)

		userCredential := entity.UserCredential{
			UserID:       123,
			Username:     "test",
			PasswordHash: hashedPassword,
		}
		mockRepo.On("Create", ctx, userCredential).Return(nil)

		err := uc.Create(ctx, userCredential.UserID, userCredential.Username, password)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Create user credential with invalid input", func(t *testing.T) {

		mockRepo := new(MockUserCredentialRepo)
		mockHasher := new(MockPasswordHasher)
		uc := NewUserCredentialUseCase(mockRepo, mockHasher)
		ctx := context.Background()
		password := "test"
		hashedPassword := "$2a$10"
		mockHasher.On("GenerateHash", ctx, password).Return(hashedPassword, nil)

		userCredential := entity.UserCredential{
			PasswordHash: hashedPassword,
		}
		mockRepo.On("Create", ctx, userCredential).Return(assert.AnError)

		err := uc.Create(ctx, userCredential.UserID, userCredential.Username, password)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
