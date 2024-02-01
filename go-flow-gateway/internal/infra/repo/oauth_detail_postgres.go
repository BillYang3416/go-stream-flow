package repo

import (
	"context"
	"fmt"

	"github.com/bgg/go-flow-gateway/internal/entity"
	"github.com/bgg/go-flow-gateway/internal/usecase/apperrors"
	"github.com/bgg/go-flow-gateway/pkg/logger"
	"github.com/bgg/go-flow-gateway/pkg/postgres"
)

type OAuthDetailRepo struct {
	*postgres.Postgres
	logger logger.Logger
}

func NewOAuthDetailRepo(pg *postgres.Postgres, l logger.Logger) *OAuthDetailRepo {
	return &OAuthDetailRepo{Postgres: pg, logger: l}
}

func (r *OAuthDetailRepo) Create(ctx context.Context, oAuthDetail entity.OAuthDetail) error {

	sql, args, err := r.Builder.
		Insert("oauth_details").
		Columns("oauth_id", "user_id", "provider", "access_token", "refresh_token").
		Values(oAuthDetail.OAuthID, oAuthDetail.UserID, oAuthDetail.Provider, oAuthDetail.AccessToken, oAuthDetail.RefreshToken).
		ToSql()

	if err != nil {
		r.logger.Error("OAuthDetailRepo - Create - r.Builder: failed to build query", "error", err)
		return fmt.Errorf("OAuthDetailRepo - Create - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		r.logger.Error("OAuthDetailRepo - Create - r.Pool.Exec : failed to execute query", "sql", sql, "error", err)
		pgErrorChecker := postgres.NewPGErrorChecker()
		if pgErrorChecker.IsUniqueViolation(err) {
			return apperrors.NewUniqueConstraintError("duplicate key", fmt.Sprintf("OAuthDetailRepo - Create - r.Pool.Exec: %s", err.Error()))
		}
		return fmt.Errorf("OAuthDetailRepo - Create - r.Pool.Exec: %w", err)
	}

	r.logger.Info("OAuthDetailRepo - Create - oauth detail created successfully", "oauth_id", oAuthDetail.OAuthID)
	return nil
}

func (r *OAuthDetailRepo) UpdateRefreshToken(ctx context.Context, oAuthID string, refreshToken string) error {

	sql, args, err := r.Builder.
		Update("oauth_details").
		Set("refresh_token", refreshToken).
		Where("oauth_id = ?", oAuthID).
		ToSql()

	if err != nil {
		r.logger.Error("OAuthDetailRepo - UpdateRefreshToken - r.Builder: failed to build query", "error", err)
		return fmt.Errorf("OAuthDetailRepo - UpdateRefreshToken - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		r.logger.Error("OAuthDetailRepo - UpdateRefreshToken - r.Pool.Exec : failed to execute query", "sql", sql, "error", err)
		return fmt.Errorf("OAuthDetailRepo - UpdateRefreshToken - r.Pool.Exec: %w", err)
	}

	r.logger.Info("OAuthDetailRepo - UpdateRefreshToken - oauth detail updated successfully", "oauth_id", oAuthID)
	return nil
}

func (r *OAuthDetailRepo) GetByOAuthID(ctx context.Context, oauthId string) (entity.OAuthDetail, error) {

	sql, args, err := r.Builder.
		Select("oauth_id", "user_id", "provider", "access_token", "refresh_token").
		From("oauth_details").
		Where("oauth_id = ?", oauthId).
		ToSql()

	if err != nil {
		r.logger.Error("OAuthDetailRepo - GetByOAuthID - r.Builder: failed to build query", "error", err)
		return entity.OAuthDetail{}, fmt.Errorf("OAuthDetailRepo - GetByOAuthID - r.Builder: %w", err)
	}

	var u entity.OAuthDetail
	row := r.Pool.QueryRow(ctx, sql, args...)
	err = row.Scan(&u.OAuthID, &u.UserID, &u.Provider, &u.AccessToken, &u.RefreshToken)
	if err != nil {
		r.logger.Error("OAuthDetailRepo - GetByOAuthID - row.Scan: failed to scan row", "error", err)
		pgErrorChecker := postgres.NewPGErrorChecker()
		if pgErrorChecker.IsNoRows(err) {
			return entity.OAuthDetail{}, apperrors.NewNoRowsAffectedError("oauth detail not found", fmt.Sprintf("OAuthDetailRepo - GetByOAuthID - row.Scan: %s", err.Error()))
		}
		return entity.OAuthDetail{}, fmt.Errorf("OAuthDetailRepo - GetByOAuthID - row.Scan: %w", err)
	}

	r.logger.Info("OAuthDetailRepo - GetByOAuthID - oauth detail retrieved successfully", "oauth_id", oauthId)
	return u, nil
}
