package repo

import (
	"context"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/bgg/go-file-gate/internal/entity"
	"github.com/bgg/go-file-gate/pkg/postgres"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

func TestUserProfileRepo_Create(t *testing.T) {

	// Arrange
	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err, "Error should not have occurred when opening a stub database connection")
	defer mock.Close()

	pg := &postgres.Postgres{Pool: mock, Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)}
	repo := NewUserProfileRepo(pg)

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
	err = repo.Create(ctx, userProfile)

	// Assert
	assert.NoError(t, err, "Error should not have occurred when creating a user profile")
	mock.ExpectationsWereMet()
}
