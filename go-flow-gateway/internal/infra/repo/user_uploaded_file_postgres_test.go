package repo

import (
	"context"
	"testing"
	"time"

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
			UserID:    123,
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

func TestUserUploadedFile_GetPaginatedFiles(t *testing.T) {

	t.Run("should return a list of user uploaded files", func(t *testing.T) {

		// Arrange
		ctx, mock, repo := setupUserUploadedFileRepoTest(t)

		lastID := 0
		userID := 123
		limit := 10

		userUploadedFiles := []entity.UserUploadedFile{
			{
				ID:             3,
				Name:           "test.txt",
				Size:           123,
				Content:        []byte("test"),
				UserID:         123,
				CreatedAt:      time.Now(),
				EmailSent:      false,
				EmailSentAt:    time.Now(),
				EmailRecipient: "",
				ErrorMessage:   "",
			},
		}

		mock.ExpectQuery("SELECT (.+) FROM user_uploaded_files").
			WithArgs(userID, lastID).
			WillReturnRows(mock.NewRows([]string{"id", "name", "size", "content", "user_id", "created_at", "email_sent", "email_sent_at", "email_recipient", "error_message"}).
				AddRow(userUploadedFiles[0].ID, userUploadedFiles[0].Name, userUploadedFiles[0].Size, userUploadedFiles[0].Content, userUploadedFiles[0].UserID, userUploadedFiles[0].CreatedAt, userUploadedFiles[0].EmailSent, userUploadedFiles[0].EmailSentAt, userUploadedFiles[0].EmailRecipient, userUploadedFiles[0].ErrorMessage))

		// Act
		files, err := repo.GetPaginatedFiles(ctx, lastID, userID, limit)

		// Assert
		assert.NoError(t, err, "Error should not have occurred when getting a list of user uploaded files")
		assert.Equal(t, userUploadedFiles, files, "The returned list of user uploaded files should match the expected list")
		mock.ExpectationsWereMet()
	})
}
