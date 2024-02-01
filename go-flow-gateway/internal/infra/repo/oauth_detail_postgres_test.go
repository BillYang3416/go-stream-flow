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

func setupOauthDetailRepoTest(t *testing.T) (context.Context, pgxmock.PgxPoolIface, *OAuthDetailRepo) {
	t.Helper()

	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err, "Error should not have occurred when opening a stub database connection")
	defer mock.Close()

	pg := &postgres.Postgres{Pool: mock, Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)}
	repo := NewOAuthDetailRepo(pg, logger.New("debug"))
	return ctx, mock, repo
}

func TestOAuthDetailRepo_Create(t *testing.T) {

	t.Run("should create an oauth detail", func(t *testing.T) {

		ctx, mock, repo := setupOauthDetailRepoTest(t)

		oauthDetail := entity.OAuthDetail{
			OAuthID:      "123",
			UserID:       123,
			Provider:     "line",
			AccessToken:  "123",
			RefreshToken: "123",
		}

		mock.ExpectExec("INSERT INTO oauth_details").
			WithArgs(oauthDetail.OAuthID, oauthDetail.UserID, oauthDetail.Provider, oauthDetail.AccessToken, oauthDetail.RefreshToken).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err := repo.Create(ctx, oauthDetail)

		assert.NoError(t, err, "Error should not have occurred when creating an oauth detail")
		mock.ExpectationsWereMet()
	})

	t.Run("should return an error when creating an oauth detail", func(t *testing.T) {

		ctx, mock, repo := setupOauthDetailRepoTest(t)

		oauthDetail := entity.OAuthDetail{}

		mock.ExpectExec("INSERT INTO oauth_details").
			WithArgs(oauthDetail.UserID, oauthDetail.Provider, oauthDetail.AccessToken, oauthDetail.RefreshToken).
			WillReturnError(assert.AnError)

		err := repo.Create(ctx, oauthDetail)

		assert.Error(t, err, "Error should have occurred when creating an oauth detail")
		mock.ExpectationsWereMet()
	})

}

func TestOAuthDetailRepo_UpdateRefreshToken(t *testing.T) {

	t.Run("should update refresh token", func(t *testing.T) {

		ctx, mock, repo := setupOauthDetailRepoTest(t)

		mock.ExpectExec("UPDATE oauth_details").
			WithArgs("123", "123").
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := repo.UpdateRefreshToken(ctx, "123", "123")

		assert.NoError(t, err, "Error should not have occurred when updating refresh token")
		mock.ExpectationsWereMet()
	})

	t.Run("should return an error when updating refresh token", func(t *testing.T) {

		ctx, mock, repo := setupOauthDetailRepoTest(t)

		mock.ExpectExec("UPDATE oauth_details").
			WithArgs("123", "123").
			WillReturnError(assert.AnError)

		err := repo.UpdateRefreshToken(ctx, "123", "123")

		assert.Error(t, err, "Error should have occurred when updating refresh token")
		mock.ExpectationsWereMet()
	})
}

func TestOAuthDetailRepo_GetByOAuthID(t *testing.T) {

	t.Run("should return an oauth detail by oauth id", func(t *testing.T) {

		ctx, mock, repo := setupOauthDetailRepoTest(t)

		oauthDetail := entity.OAuthDetail{
			OAuthID:      "123",
			UserID:       123,
			Provider:     "line",
			AccessToken:  "123",
			RefreshToken: "123",
		}

		mock.ExpectQuery("SELECT").
			WithArgs(oauthDetail.OAuthID).
			WillReturnRows(mock.NewRows([]string{"oauth_id", "user_id", "provider", "access_token", "refresh_token"}).
				AddRow(oauthDetail.OAuthID, oauthDetail.UserID, oauthDetail.Provider, oauthDetail.AccessToken, oauthDetail.RefreshToken))

		result, err := repo.GetByOAuthID(ctx, oauthDetail.OAuthID)

		assert.NoError(t, err, "Error should not have occurred when getting an oauth detail by oauth id")
		assert.Equal(t, oauthDetail, result)
		mock.ExpectationsWereMet()
	})

	t.Run("should return an error when getting an oauth detail by oauth id", func(t *testing.T) {

		ctx, mock, repo := setupOauthDetailRepoTest(t)

		mock.ExpectQuery("SELECT").
			WithArgs("123").
			WillReturnError(assert.AnError)

		_, err := repo.GetByOAuthID(ctx, "123")

		assert.Error(t, err, "Error should have occurred when getting an oauth detail by oauth id")
		mock.ExpectationsWereMet()
	})
}
