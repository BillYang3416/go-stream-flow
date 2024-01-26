package repo

import (
	"context"
	"fmt"

	"github.com/bgg/go-flow-gateway/internal/entity"
	"github.com/bgg/go-flow-gateway/internal/usecase/apperrors"
	"github.com/bgg/go-flow-gateway/pkg/postgres"
)

type OAuthDetailRepo struct {
	*postgres.Postgres
}

func NewOAuthDetailRepo(pg *postgres.Postgres) *OAuthDetailRepo {
	return &OAuthDetailRepo{Postgres: pg}
}

func (r *OAuthDetailRepo) Create(ctx context.Context, u entity.OAuthDetail) error {

	sql, args, err := r.Builder.
		Insert("oauth_details").
		Columns("oauth_id", "user_id", "provider", "access_token", "refresh_token").
		Values(u.OAuthID, u.UserID, u.Provider, u.AccessToken, u.RefreshToken).
		ToSql()

	if err != nil {
		return fmt.Errorf("OAuthDetailRepo - Create - r.Builder: %w", err)
	}
	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		pgErrorChecker := postgres.NewPGErrorChecker()
		if pgErrorChecker.IsUniqueViolation(err) {
			return apperrors.NewUniqueConstraintError("duplicate key", fmt.Sprintf("OAuthDetailRepo - Create - r.Pool.Exec: %s", err.Error()))
		}
		return fmt.Errorf("OAuthDetailRepo - Create - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *OAuthDetailRepo) UpdateRefreshToken(ctx context.Context, oauthID string, refreshToken string) error {

	sql, args, err := r.Builder.
		Update("oauth_details").
		Set("refresh_token", refreshToken).
		Where("oauth_id = ?", oauthID).
		ToSql()

	if err != nil {
		return fmt.Errorf("OAuthDetailRepo - UpdateRefreshToken - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("OAuthDetailRepo - UpdateRefreshToken - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *OAuthDetailRepo) GetByOAuthID(ctx context.Context, oauthId string) (entity.OAuthDetail, error) {

	sql, args, err := r.Builder.
		Select("oauth_id", "user_id", "provider", "access_token", "refresh_token").
		From("oauth_details").
		Where("oauth_id = ?", oauthId).
		ToSql()

	if err != nil {
		return entity.OAuthDetail{}, fmt.Errorf("OAuthDetailRepo - GetByOAuthID - r.Builder: %w", err)
	}

	var u entity.OAuthDetail
	row := r.Pool.QueryRow(ctx, sql, args...)
	err = row.Scan(&u.OAuthID, &u.UserID, &u.Provider, &u.AccessToken, &u.RefreshToken)
	if err != nil {
		pgErrorChecker := postgres.NewPGErrorChecker()
		if pgErrorChecker.IsNoRows(err) {
			return entity.OAuthDetail{}, apperrors.NewNoRowsAffectedError("oauth detail not found", fmt.Sprintf("OAuthDetailRepo - GetByOAuthID - row.Scan: %s", err.Error()))
		}
		return entity.OAuthDetail{}, fmt.Errorf("OAuthDetailRepo - GetByOAuthID - row.Scan: %w", err)
	}

	return u, nil
}
