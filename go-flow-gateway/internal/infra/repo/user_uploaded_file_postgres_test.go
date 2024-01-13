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

func setupUserUploadedFileRepoTest(t *testing.T) (context.Context, pgxmock.PgxPoolIface, *UserUploadedFileRepo) {

	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err, "Error should not have occurred when opening a stub database connection")
	defer mock.Close()

	pg := &postgres.Postgres{Pool: mock, Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)}
	repo := NewUserUploadedFileRepo(pg)
	return ctx, mock, repo
}

func TestUserUploadedFile_Created(t *testing.T) {

	t.Run("should create a user uploaded file", func(t *testing.T) {

		// Arrange
		ctx, mock, repo := setupUserUploadedFileRepoTest(t)

		userUploadedFile := entity.UserUploadedFile{
			Name:      "test.txt",
			Size:      123,
			Content:   []byte("test"),
			UserID:    "123",
			EmailSent: false,
		}

		mock.ExpectExec("INSERT INTO user_uploaded_files").
			WithArgs(userUploadedFile.Name, userUploadedFile.Size, userUploadedFile.Content, userUploadedFile.UserID, userUploadedFile.EmailSent).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		// Act
		err := repo.Create(ctx, userUploadedFile)

		// Assert
		assert.NoError(t, err, "Error should not have occurred when creating a user uploaded file")
		mock.ExpectationsWereMet()
	})

	t.Run("should return an error when creating a user uploaded file", func(t *testing.T) {

		// Arrange
		ctx, mock, repo := setupUserUploadedFileRepoTest(t)

		userUploadedFile := entity.UserUploadedFile{}

		mock.ExpectExec("INSERT INTO user_uploaded_files").
			WithArgs(userUploadedFile.Name, userUploadedFile.Size, userUploadedFile.Content, userUploadedFile.UserID, userUploadedFile.EmailSent).
			WillReturnError(assert.AnError)

		// Act
		err := repo.Create(ctx, userUploadedFile)

		// Assert
		assert.Error(t, err, "Error should have occurred when creating a user uploaded file")
		mock.ExpectationsWereMet()
	})

}
