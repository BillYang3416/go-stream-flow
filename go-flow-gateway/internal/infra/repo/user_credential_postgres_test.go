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

func setupUserCredentialRepoTest(t *testing.T) (context.Context, pgxmock.PgxPoolIface, *UserCredentialRepo) {
	t.Helper()

	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err, "Error should not have occurred when opening a stub database connection")
	defer mock.Close()

	pg := &postgres.Postgres{Pool: mock, Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)}
	repo := NewUserCredentialRepo(pg, logger.New("debug"))

	return ctx, mock, repo
}

func TestUserCredentialRepo_Create(t *testing.T) {

	const (
		userID       = 123
		userName     = "123"
		passwordHash = "123"
	)

	t.Run("should create a user credential", func(t *testing.T) {

		ctx, mock, repo := setupUserCredentialRepoTest(t)

		userCredential := entity.UserCredential{
			UserID:       userID,
			Username:     userName,
			PasswordHash: passwordHash,
		}

		mock.ExpectExec("INSERT INTO user_credentials").
			WithArgs(userCredential.UserID, userCredential.Username, userCredential.PasswordHash).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err := repo.Create(ctx, userCredential)
		assert.NoError(t, err, "Error should not have occurred when creating a user credential")
		mock.ExpectationsWereMet()
	})

	t.Run("should return an error when creating a user credential", func(t *testing.T) {

		ctx, mock, repo := setupUserCredentialRepoTest(t)

		userCredential := entity.UserCredential{}

		mock.ExpectExec("INSERT INTO user_credentials").
			WithArgs(userCredential.UserID, userCredential.Username, userCredential.PasswordHash).
			WillReturnError(assert.AnError)

		err := repo.Create(ctx, userCredential)
		assert.Error(t, err, "Error should have occurred when creating a user credential")
		mock.ExpectationsWereMet()
	})
}

func TestUserCredentialRepo_GetByUsername(t *testing.T) {

	const (
		userID       = 123
		userName     = "123"
		passwordHash = "123"
	)

	t.Run("should get a user credential by username", func(t *testing.T) {

		ctx, mock, repo := setupUserCredentialRepoTest(t)

		userCredential := entity.UserCredential{
			UserID:       userID,
			Username:     userName,
			PasswordHash: passwordHash,
		}

		mock.ExpectQuery("SELECT user_id, username, password_hash FROM user_credentials").
			WithArgs(userCredential.Username).
			WillReturnRows(
				mock.NewRows([]string{"user_id", "username", "password_hash"}).
					AddRow(userCredential.UserID, userCredential.Username, userCredential.PasswordHash),
			)

		u, err := repo.GetByUsername(ctx, userCredential.Username)
		assert.NoError(t, err, "Error should not have occurred when getting a user credential by username")
		assert.Equal(t, userCredential, u, "UserCredential should have been returned")
		mock.ExpectationsWereMet()
	})

	t.Run("should return an error when getting a user credential by username", func(t *testing.T) {

		ctx, mock, repo := setupUserCredentialRepoTest(t)

		userCredential := entity.UserCredential{}

		mock.ExpectQuery("SELECT user_id, username, password_hash FROM user_credentials").
			WithArgs(userCredential.Username).
			WillReturnError(assert.AnError)

		_, err := repo.GetByUsername(ctx, userCredential.Username)
		assert.Error(t, err, "Error should have occurred when getting a user credential by username")
		mock.ExpectationsWereMet()
	})
}
