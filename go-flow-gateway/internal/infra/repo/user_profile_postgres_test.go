package repo

import (
	"context"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/bgg/go-flow-gateway/internal/entity"
	"github.com/bgg/go-flow-gateway/pkg/logger"
	"github.com/bgg/go-flow-gateway/pkg/postgres"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

func setupUserProfileRepoTest(t *testing.T) (context.Context, pgxmock.PgxPoolIface, *UserProfileRepo) {
	t.Helper()

	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err, "Error should not have occurred when opening a stub database connection")
	defer mock.Close()

	pg := &postgres.Postgres{Pool: mock, Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)}
	repo := NewUserProfileRepo(pg, logger.New("debug"))
	return ctx, mock, repo
}

func TestUserProfileRepo_Create(t *testing.T) {

	const (
		userID      = 1234567890
		displayName = "hank"
		pictureURL  = "https://example.com"
	)

	t.Run("should create a user profile", func(t *testing.T) {
		// Arrange
		ctx, mock, repo := setupUserProfileRepoTest(t)

		userProfile := entity.UserProfile{
			UserID:      userID,
			DisplayName: displayName,
			PictureURL:  pictureURL,
		}

		rows := pgxmock.NewRows([]string{"user_id"}).AddRow(userID)

		mock.ExpectQuery("INSERT INTO user_profiles").
			WithArgs(userProfile.DisplayName, userProfile.PictureURL).
			WillReturnRows(rows)

		// Act
		_, err := repo.Create(ctx, userProfile)

		// Assert
		assert.NoError(t, err, "Error should not have occurred when creating a user profile")
		mock.ExpectationsWereMet()
	})

	t.Run("should return an error when creating a user profile", func(t *testing.T) {
		// Arrange
		ctx, mock, repo := setupUserProfileRepoTest(t)

		userProfile := entity.UserProfile{}

		mock.ExpectQuery("INSERT INTO user_profiles").
			WithArgs(userProfile.UserID, userProfile.DisplayName, userProfile.PictureURL).
			WillReturnError(assert.AnError)

		// Act
		_, err := repo.Create(ctx, userProfile)

		// Assert
		assert.Error(t, err, "Error should have occurred when creating a user profile")
		mock.ExpectationsWereMet()
	})
}

func TestUserProfileRepo_GetByID(t *testing.T) {

	const (
		userID      = 1234567890
		displayName = "hank"
		pictureURL  = "https://example.com"
	)

	t.Run("should return a user profile by id", func(t *testing.T) {
		// Arrange
		ctx, mock, repo := setupUserProfileRepoTest(t)

		userProfile := entity.UserProfile{
			UserID:      userID,
			DisplayName: displayName,
			PictureURL:  pictureURL,
		}

		mock.ExpectQuery("SELECT").
			WithArgs(userID).
			WillReturnRows(pgxmock.NewRows([]string{"user_id", "display_name", "picture_url"}).
				AddRow(userID, displayName, pictureURL))

		// Act
		result, err := repo.GetByID(ctx, userID)

		// Assert
		assert.NoError(t, err, "Error should not have occurred when getting a user profile")
		assert.Equal(t, userProfile, result, "User profile should be equal to expected")
		mock.ExpectationsWereMet()
	})

	t.Run("should return an error when getting a user profile by id", func(t *testing.T) {
		// Arrange
		ctx, mock, repo := setupUserProfileRepoTest(t)

		mock.ExpectQuery("SELECT").
			WithArgs(userID).
			WillReturnError(assert.AnError)

		// Act
		result, err := repo.GetByID(ctx, userID)

		// Assert
		assert.Error(t, err, "Error should have occurred when getting a user profile")
		assert.Equal(t, entity.UserProfile{}, result, "User profile should be equal to expected")
		mock.ExpectationsWereMet()
	})
}
