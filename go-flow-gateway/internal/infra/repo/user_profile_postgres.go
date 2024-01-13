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

func (r *UserProfileRepo) Create(ctx context.Context, u entity.UserProfile) (entity.UserProfile, error) {
	// Build the SQL query using squirrel
	sql, args, err := r.Builder.
		Insert("user_profiles").                // Assuming 'user_profiles' is the table name
		Columns("display_name", "picture_url"). // Columns in the table
		Values(u.DisplayName, u.PictureURL).
		Suffix("RETURNING user_id"). // Corresponding values from the UserProfile entity
		ToSql()

	if err != nil {
		return entity.UserProfile{}, fmt.Errorf("UserProfileRepo - Create - r.Builder: %w", err)
	}

	var userID int
	// Use QueryRow to execute the query and scan the user_id directly into the userID variable
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&userID)
	if err != nil {
		pgErrorChecker := postgres.NewPGErrorChecker()
		if pgErrorChecker.IsUniqueViolation(err) {
			return entity.UserProfile{}, NewUniqueConstraintError("duplicate key", fmt.Sprintf("UserProfileRepo - Create - r.Pool.Exec: %s", err.Error()))
		}
		return entity.UserProfile{}, fmt.Errorf("UserProfileRepo - Create - r.Pool.Exec: %w", err)
	}

	// Set the userID in the UserProfile entity before returning
	u.UserID = userID
	return u, nil
}

func (r *UserProfileRepo) GetByID(ctx context.Context, userId int) (entity.UserProfile, error) {
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
