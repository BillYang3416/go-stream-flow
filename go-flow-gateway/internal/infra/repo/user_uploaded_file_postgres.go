package repo

import (
	"context"
	"fmt"

	"github.com/bgg/go-flow-gateway/internal/entity"
	"github.com/bgg/go-flow-gateway/pkg/postgres"
)

type UserUploadedFileRepo struct {
	*postgres.Postgres
}

func NewUserUploadedFileRepo(pg *postgres.Postgres) *UserUploadedFileRepo {
	return &UserUploadedFileRepo{Postgres: pg}
}

func (r *UserUploadedFileRepo) Create(ctx context.Context, u entity.UserUploadedFile) error {
	// Build the SQL query using squirrel
	sql, args, err := r.Builder.
		Insert("user_uploaded_files").
		Columns("name", "size", "content", "user_id", "email_sent").
		Values(u.Name, u.Size, u.Content, u.UserID, u.EmailSent).
		ToSql()

	if err != nil {
		return fmt.Errorf("UserUploadedFileRepo - Save - r.Builder: %w", err)
	}

	// Execute the query using pgx
	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		pgErrorChecker := postgres.NewPGErrorChecker()
		if pgErrorChecker.IsUniqueViolation(err) {
			return NewUniqueConstraintError("duplicate key", fmt.Sprintf("UserUploadedFileRepo - Save - r.Pool.Exec: %s", err.Error()))
		}
		return fmt.Errorf("UserUploadedFileRepo - Save - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *UserUploadedFileRepo) GetPaginatedFiles(ctx context.Context, lastID, userID, limit int) ([]entity.UserUploadedFile, error) {
	// Build the SQL query using squirrel
	sql, args, err := r.Builder.
		Select("id", "name", "size", "content", "user_id", "created_at", "email_sent", "email_sent_at", "email_recipient", "error_message").
		From("user_uploaded_files").
		Where("user_id = ?", userID).
		Where("id > ?", lastID).
		Limit(uint64(limit)).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("UserUploadedFileRepo - GetPaginatedFiles - r.Builder: %w", err)
	}

	// Execute the query using pgx
	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("UserUploadedFileRepo - GetPaginatedFiles - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var files []entity.UserUploadedFile
	for rows.Next() {
		var file entity.UserUploadedFile
		err := rows.Scan(&file.ID, &file.Name, &file.Size, &file.Content, &file.UserID, &file.CreatedAt, &file.EmailSent, &file.EmailSentAt, &file.EmailRecipient, &file.ErrorMessage)
		if err != nil {
			return nil, fmt.Errorf("UserUploadedFileRepo - GetPaginatedFiles - rows.Scan: %w", err)
		}
		files = append(files, file)
	}

	return files, nil
}
