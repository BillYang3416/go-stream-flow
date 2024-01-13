package repo

import (
	"context"
	"fmt"

	"github.com/bgg/go-flow-gateway/internal/entity"
	"github.com/bgg/go-flow-gateway/pkg/postgres"
)

type UserProfileRepo struct {
	*postgres.Postgres
}

func NewUserProfileRepo(pg *postgres.Postgres) *UserProfileRepo {
	return &UserProfileRepo{Postgres: pg}
}

func (r *UserProfileRepo) Create(ctx context.Context, u entity.UserProfile) error {
	// Build the SQL query using squirrel
	sql, args, err := r.Builder.
		Insert("user_profiles").                                                            // Assuming 'user_profiles' is the table name
		Columns("user_id", "display_name", "picture_url", "access_token", "refresh_token"). // Columns in the table
		Values(u.UserID, u.DisplayName, u.PictureURL, u.AccessToken, u.RefreshToken).       // Corresponding values from the UserProfile entity
		ToSql()

	if err != nil {
		return fmt.Errorf("UserProfileRepo - Save - r.Builder: %w", err)
	}

	// Execute the query using pgx
	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		pgErrorChecker := postgres.NewPGErrorChecker()
		if pgErrorChecker.IsUniqueViolation(err) {
			return NewUniqueConstraintError("duplicate key", fmt.Sprintf("UserProfileRepo - Save - r.Pool.Exec: %s", err.Error()))
		}
		return fmt.Errorf("UserProfileRepo - Save - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *UserProfileRepo) GetByID(ctx context.Context, userId string) (entity.UserProfile, error) {
	sql, args, err := r.Builder.Select("user_id", "display_name", "picture_url").From("user_profiles").Where("user_id = ?", userId).ToSql()

	if err != nil {
		return entity.UserProfile{}, fmt.Errorf("UserProfileRepo - GetByID - r.Builder: %w", err)
	}

	var u entity.UserProfile
	row := r.Pool.QueryRow(ctx, sql, args...)
	err = row.Scan(&u.UserID, &u.DisplayName, &u.PictureURL)
	if err != nil {
		pgErrorChecker := postgres.NewPGErrorChecker()
		if pgErrorChecker.IsNoRows(err) {
			return entity.UserProfile{}, NewNoRowsAffectedError("user profile not found", fmt.Sprintf("UserProfileRepo - GetByID - row.Scan: %s", err.Error()))
		}
		return entity.UserProfile{}, fmt.Errorf("UserProfileRepo - GetByID - row.Scan: %w", err)
	}

	return u, nil
}

func (r *UserProfileRepo) UpdateRefreshToken(ctx context.Context, userId string, refreshToken string) error {
	sql, args, err := r.Builder.Update("user_profiles").Set("refresh_token", refreshToken).Where("user_id = ?", userId).ToSql()

	if err != nil {
		return fmt.Errorf("UserProfileRepo - UpdateRefreshToken - r.Builder: %w", err)
	}

	commandTag, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UserProfileRepo - UpdateRefreshToken - r.Pool.Exec: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return NewNoRowsAffectedError("user profile is not found", "UserProfileRepo - UpdateRefreshToken - r.Pool.Exec: no rows affected")
	}

	return nil
}
