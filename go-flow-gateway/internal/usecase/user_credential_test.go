package usecase

import (
	"context"
	"testing"

	"github.com/bgg/go-flow-gateway/internal/entity"
	"github.com/bgg/go-flow-gateway/internal/usecase/apperrors"
	"github.com/bgg/go-flow-gateway/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserCredentialRepo struct {
	mock.Mock
}

type MockPasswordHasher struct {
	mock.Mock
}

func (m *MockUserCredentialRepo) Create(ctx context.Context, userCredential entity.UserCredential) error {
	args := m.Called(ctx, userCredential)
	return args.Error(0)
}

func (m *MockUserCredentialRepo) GetByUsername(ctx context.Context, username string) (entity.UserCredential, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(entity.UserCredential), args.Error(1)
}

func (m *MockPasswordHasher) GenerateHash(ctx context.Context, password string) (string, error) {
	args := m.Called(ctx, password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordHasher) CompareHash(ctx context.Context, password, hashedPassword string) error {
	args := m.Called(ctx, password, hashedPassword)
	return args.Error(0)
}

func setupUserCredentialUsecase(t *testing.T) (UserCredential, *MockUserCredentialRepo, *MockPasswordHasher, *MockUserProfileUseCase) {
	t.Helper()
	mockRepo := new(MockUserCredentialRepo)
	mockHasher := new(MockPasswordHasher)
	mockUserProfileUseCase := new(MockUserProfileUseCase)
	uc := NewUserCredentialUseCase(mockRepo, mockHasher, mockUserProfileUseCase, logger.New("debug"))
	return uc, mockRepo, mockHasher, mockUserProfileUseCase
}

func TestUserCredentialUsecase_Register(t *testing.T) {

	const (
		displayName    = "hank"
		username       = "test"
		password       = "test"
		hashedPassword = "$2a$10"
		userID         = 123
	)

	t.Run("user registered successfully", func(t *testing.T) {

		uc, mockRepo, mockHasher, mockUserProfileUseCase := setupUserCredentialUsecase(t)

		ctx := context.Background()

		userCredential := entity.UserCredential{
			UserID:       userID,
			Username:     username,
			PasswordHash: hashedPassword,
		}

		mockRepo.On("GetByUsername", ctx, username).Return(entity.UserCredential{}, apperrors.NewNoRowsAffectedError("test", "test"))
		mockUserProfileUseCase.On("Create", ctx, entity.UserProfile{
			DisplayName: displayName,
		}).Return(entity.UserProfile{
			UserID: userID,
		}, nil)
		mockHasher.On("GenerateHash", ctx, password).Return(hashedPassword, nil)
		mockRepo.On("Create", ctx, userCredential).Return(nil)

		gotUserID, err := uc.Register(ctx, displayName, username, password)

		assert.NoError(t, err)
		assert.Equal(t, userID, gotUserID)
		mockUserProfileUseCase.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
	})

	t.Run("Create user credential with invalid input", func(t *testing.T) {

		uc, mockRepo, mockHasher, mockUserProfileUseCase := setupUserCredentialUsecase(t)

		ctx := context.Background()

		userCredential := entity.UserCredential{
			UserID:       userID,
			Username:     username,
			PasswordHash: hashedPassword,
		}

		mockRepo.On("GetByUsername", ctx, username).Return(entity.UserCredential{}, apperrors.NewNoRowsAffectedError("test", "test"))
		mockUserProfileUseCase.On("Create", ctx, entity.UserProfile{
			DisplayName: displayName,
		}).Return(entity.UserProfile{
			UserID: userID,
		}, nil)
		mockHasher.On("GenerateHash", ctx, password).Return(hashedPassword, nil)
		mockRepo.On("Create", ctx, userCredential).Return(assert.AnError)

		_, err := uc.Register(ctx, displayName, username, password)

		assert.Error(t, err)
		mockUserProfileUseCase.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
	})
}

func TestUserCredentialUsecase_GetByUsername(t *testing.T) {

	const (
		userID         = 123
		userName       = "test"
		hashedPassword = "$2a$10"
	)

	t.Run("Get user credential by username successfully", func(t *testing.T) {

		uc, mockRepo, _, _ := setupUserCredentialUsecase(t)

		ctx := context.Background()

		userCredential := entity.UserCredential{
			UserID:       userID,
			Username:     userName,
			PasswordHash: hashedPassword,
		}

		mockRepo.On("GetByUsername", ctx, userCredential.Username).Return(userCredential, nil)

		u, err := uc.GetByUsername(ctx, userCredential.Username)

		assert.NoError(t, err)
		assert.Equal(t, userCredential, u)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Get user credential by username with invalid input", func(t *testing.T) {

		uc, mockRepo, _, _ := setupUserCredentialUsecase(t)

		ctx := context.Background()

		userCredential := entity.UserCredential{
			UserID:       userID,
			Username:     userName,
			PasswordHash: hashedPassword,
		}

		mockRepo.
			On("GetByUsername", ctx, userCredential.Username).
			Return(entity.UserCredential{}, assert.AnError)

		u, err := uc.GetByUsername(ctx, userCredential.Username)

		assert.Error(t, err)
		assert.Equal(t, entity.UserCredential{}, u)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserCredentialUsecase_Login(t *testing.T) {

	const (
		userID         = 123
		userName       = "test"
		password       = "test"
		hashedPassword = "$2a$10"
	)

	t.Run("Login successfully", func(t *testing.T) {

		uc, mockRepo, mockHasher, _ := setupUserCredentialUsecase(t)

		ctx := context.Background()

		userCredential := entity.UserCredential{
			UserID:       userID,
			Username:     userName,
			PasswordHash: hashedPassword,
		}

		mockRepo.On("GetByUsername", ctx, userCredential.Username).Return(userCredential, nil)
		mockHasher.On("CompareHash", ctx, password, userCredential.PasswordHash).Return(nil)

		u, err := uc.Login(ctx, userCredential.Username, password)

		assert.NoError(t, err)
		assert.Equal(t, userCredential, u)
		mockRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
	})

	t.Run("Login with invalid input", func(t *testing.T) {

		uc, mockRepo, mockHasher, _ := setupUserCredentialUsecase(t)

		ctx := context.Background()

		userCredential := entity.UserCredential{
			UserID:       userID,
			Username:     userName,
			PasswordHash: hashedPassword,
		}

		mockRepo.On("GetByUsername", ctx, userCredential.Username).Return(userCredential, nil)
		mockHasher.On("CompareHash", ctx, password, userCredential.PasswordHash).Return(assert.AnError)

		u, err := uc.Login(ctx, userCredential.Username, password)

		assert.Error(t, err)
		assert.Equal(t, entity.UserCredential{}, u)
		mockRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
	})
}
