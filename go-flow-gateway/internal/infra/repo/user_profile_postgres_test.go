package repo

import (
	"context"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/bgg/go-flow-gateway/internal/entity"
	"github.com/bgg/go-flow-gateway/pkg/postgres"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

func setupUserProfileRepoTest(t *testing.T) (context.Context, pgxmock.PgxPoolIface, *UserProfileRepo) {
	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err, "Error should not have occurred when opening a stub database connection")
	defer mock.Close()

	pg := &postgres.Postgres{Pool: mock, Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)}
	repo := NewUserProfileRepo(pg)
	return ctx, mock, repo
}

func TestUserProfileRepo_Create(t *testing.T) {

	t.Run("should create a user profile", func(t *testing.T) {
		// Arrange
		ctx, mock, repo := setupUserProfileRepoTest(t)

		userProfile := entity.UserProfile{
			UserID:       "123",
			DisplayName:  "John Doe",
			PictureURL:   "https://example.com/picture.jpg",
			AccessToken:  "123",
			RefreshToken: "123",
		}

		mock.ExpectExec("INSERT INTO user_profiles").
			WithArgs(userProfile.UserID, userProfile.DisplayName, userProfile.PictureURL, userProfile.AccessToken, userProfile.RefreshToken).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		// Act
		err := repo.Create(ctx, userProfile)

		// Assert
		assert.NoError(t, err, "Error should not have occurred when creating a user profile")
		mock.ExpectationsWereMet()
	})

	t.Run("should return an error when creating a user profile", func(t *testing.T) {
		// Arrange
		ctx, mock, repo := setupUserProfileRepoTest(t)

		userProfile := entity.UserProfile{}

		mock.ExpectExec("INSERT INTO user_profiles").
			WithArgs(userProfile.UserID, userProfile.DisplayName, userProfile.PictureURL, userProfile.AccessToken, userProfile.RefreshToken).
			WillReturnError(assert.AnError)

		// Act
		err := repo.Create(ctx, userProfile)

		// Assert
		assert.Error(t, err, "Error should have occurred when creating a user profile")
		mock.ExpectationsWereMet()
	})
}

func TestUserProfileRepo_GetByID(t *testing.T) {

	t.Run("should return a user profile by id", func(t *testing.T) {
		// Arrange
		ctx, mock, repo := setupUserProfileRepoTest(t)

		userProfile := entity.UserProfile{
			UserID:      "123",
			DisplayName: "John Doe",
			PictureURL:  "https://example.com/picture.jpg",
		}

		mock.ExpectQuery("SELECT").
			WithArgs(userProfile.UserID).
			WillReturnRows(pgxmock.NewRows([]string{"user_id", "display_name", "picture_url"}).
				AddRow(userProfile.UserID, userProfile.DisplayName, userProfile.PictureURL))

		// Act
		result, err := repo.GetByID(ctx, userProfile.UserID)

		// Assert
		assert.NoError(t, err, "Error should not have occurred when getting a user profile")
		assert.Equal(t, userProfile, result, "User profile should be equal to expected")
		mock.ExpectationsWereMet()
	})

	t.Run("should return an error when getting a user profile by id", func(t *testing.T) {
		// Arrange
		ctx, mock, repo := setupUserProfileRepoTest(t)

		userProfile := entity.UserProfile{}

		mock.ExpectQuery("SELECT").
			WithArgs(userProfile.UserID).
			WillReturnError(assert.AnError)

		// Act
		result, err := repo.GetByID(ctx, userProfile.UserID)

		// Assert
		assert.Error(t, err, "Error should have occurred when getting a user profile")
		assert.Equal(t, entity.UserProfile{}, result, "User profile should be equal to expected")
		mock.ExpectationsWereMet()
	})
}

func TestUserProfileRepo_UpdateRefreshToken(t *testing.T) {

	t.Run("should update a refresh token", func(t *testing.T) {
		// Arrange
		ctx, mock, repo := setupUserProfileRepoTest(t)

		userID := "123"
		refreshToken := "test"

		mock.ExpectExec("UPDATE user_profiles").
			WithArgs(refreshToken, userID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		// Act
		err := repo.UpdateRefreshToken(ctx, userID, refreshToken)

		// Assert
		assert.NoError(t, err, "Error should not have occurred when updating a refresh token")
		mock.ExpectationsWereMet()
	})

	t.Run("should return an error when updating a refresh token", func(t *testing.T) {
		// Arrange
		ctx, mock, repo := setupUserProfileRepoTest(t)

		userID := "123"
		refreshToken := "test"

		mock.ExpectExec("UPDATE user_profiles").
			WithArgs(refreshToken, userID).
			WillReturnError(assert.AnError)

		// Act
		err := repo.UpdateRefreshToken(ctx, userID, refreshToken)

		// Assert
		assert.Error(t, err, "Error should have occurred when updating a refresh token")
		mock.ExpectationsWereMet()
	})
}
