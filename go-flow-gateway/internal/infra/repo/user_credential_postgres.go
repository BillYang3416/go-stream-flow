package repo

import (
	"context"
	"fmt"

	"github.com/bgg/go-flow-gateway/internal/entity"
	"github.com/bgg/go-flow-gateway/pkg/postgres"
)

type UserCredentialRepo struct {
	*postgres.Postgres
}

func NewUserCredentialRepo(pg *postgres.Postgres) *UserCredentialRepo {
	return &UserCredentialRepo{Postgres: pg}
}

func (r *UserCredentialRepo) Create(ctx context.Context, u entity.UserCredential) error {

	sql, args, err := r.Builder.
		Insert("user_credentials").
		Columns("user_id", "username", "password_hash").
		Values(u.UserID, u.Username, u.PasswordHash).
		ToSql()

	if err != nil {
		return fmt.Errorf("UserCredentialRepo - Create - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		pgErrorChecker := postgres.NewPGErrorChecker()
		if pgErrorChecker.IsUniqueViolation(err) {
			return NewUniqueConstraintError("duplicate key", fmt.Sprintf("UserCredentialRepo - Create - r.Pool.Exec: %s", err.Error()))
		}
		return fmt.Errorf("UserCredentialRepo - Create - r.Pool.Exec: %w", err)
	}

	return nil
}