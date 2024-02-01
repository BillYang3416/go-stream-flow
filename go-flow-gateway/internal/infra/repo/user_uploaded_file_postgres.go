package repo

import (
	"context"
	"fmt"

	"github.com/bgg/go-flow-gateway/internal/entity"
	"github.com/bgg/go-flow-gateway/pkg/logger"
	"github.com/bgg/go-flow-gateway/pkg/postgres"
)

type UserUploadedFileRepo struct {
	*postgres.Postgres
	logger logger.Logger
}

func NewUserUploadedFileRepo(pg *postgres.Postgres, l logger.Logger) *UserUploadedFileRepo {
	return &UserUploadedFileRepo{Postgres: pg, logger: l}
}

func (r *UserUploadedFileRepo) Create(ctx context.Context, u entity.UserUploadedFile) (int, error) {
	// Build the SQL query using squirrel
	sql, args, err := r.Builder.
		Insert("user_uploaded_files").
		Columns("name", "size", "content", "user_id", "email_sent", "email_recipient").
		Values(u.Name, u.Size, u.Content, u.UserID, u.EmailSent, u.EmailRecipient).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		r.logger.Error("UserUploadedFileRepo - Create - r.Builder: failed to build query", "error", err)
		return 0, fmt.Errorf("UserUploadedFileRepo - Create - r.Builder: %w", err)
	}

	// Execute the query using pgx
	var userUploadedFileID int
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&userUploadedFileID)
	if err != nil {
		r.logger.Error("UserUploadedFileRepo - Create - r.Pool.QueryRow: failed to execute query", "error", err)
		return 0, fmt.Errorf("UserUploadedFileRepo - Save - r.Pool.QueryRow: %w", err)
	}

	r.logger.Info("UserUploadedFileRepo - Create: successfully created user uploaded file", "userUploadedFileID", userUploadedFileID)
	return userUploadedFileID, nil
}

func (r *UserUploadedFileRepo) GetPaginatedFiles(ctx context.Context, lastID, userID, limit int) ([]entity.UserUploadedFile, int, error) {

	// Query to get the total number of records first
	var totalRecords int
	countSql, countArgs, err := r.Builder.
		Select("COUNT(id)").
		From("user_uploaded_files").
		Where("user_id = ?", userID).
		ToSql()

	if err != nil {
		r.logger.Error("UserUploadedFileRepo - GetPaginatedFiles - r.Builder: failed to build total records query", "error", err)
		return nil, 0, fmt.Errorf("UserUploadedFileRepo - GetPaginatedFiles - r.Builder: %w", err)
	}

	err = r.Pool.QueryRow(ctx, countSql, countArgs...).Scan(&totalRecords)
	if err != nil {
		r.logger.Error("UserUploadedFileRepo - GetPaginatedFiles - r.Pool.QueryRow: failed to execute total records query", "error", err)
		return nil, 0, fmt.Errorf("UserUploadedFileRepo - GetPaginatedFiles - r.Pool.QueryRow: %w", err)
	}

	// Build the SQL query using squirrel
	sql, args, err := r.Builder.
		Select("id", "name", "size", "content", "user_id", "created_at", "email_sent", "email_sent_at", "email_recipient", "error_message").
		From("user_uploaded_files").
		Where("user_id = ?", userID).
		Where("id > ?", lastID).
		Limit(uint64(limit)).
		OrderBy("id ASC").
		ToSql()

	if err != nil {
		r.logger.Error("UserUploadedFileRepo - GetPaginatedFiles - r.Builder: failed to build user uploaded files query", "error", err)
		return nil, totalRecords, fmt.Errorf("UserUploadedFileRepo - GetPaginatedFiles - r.Builder: %w", err)
	}

	// Execute the query using pgx
	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.logger.Error("UserUploadedFileRepo - GetPaginatedFiles - r.Pool.Query: failed to execute user uploaded files query", "error", err)
		return nil, totalRecords, fmt.Errorf("UserUploadedFileRepo - GetPaginatedFiles - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var files []entity.UserUploadedFile
	for rows.Next() {
		var file entity.UserUploadedFile
		err := rows.Scan(&file.ID, &file.Name, &file.Size, &file.Content, &file.UserID, &file.CreatedAt, &file.EmailSent, &file.EmailSentAt, &file.EmailRecipient, &file.ErrorMessage)
		if err != nil {
			r.logger.Error("UserUploadedFileRepo - GetPaginatedFiles - rows.Scan: failed to scan user uploaded files query", "error", err)
			return nil, totalRecords, fmt.Errorf("UserUploadedFileRepo - GetPaginatedFiles - rows.Scan: %w", err)
		}
		files = append(files, file)
	}

	r.logger.Info("UserUploadedFileRepo - GetPaginatedFiles: successfully retrieved user uploaded files", "totalRecords", totalRecords)
	return files, totalRecords, nil
}

func (r *UserUploadedFileRepo) UpdateEmailSent(ctx context.Context, ID int) error {
	// Build the SQL query using squirrel
	sql, args, err := r.Builder.
		Update("user_uploaded_files").
		Set("email_sent", true).
		Set("email_sent_at", "NOW()").
		Where("id = ?", ID).
		ToSql()

	if err != nil {
		r.logger.Error("UserUploadedFileRepo - UpdateEmailSent - r.Builder: failed to build query", "error", err)
		return fmt.Errorf("UserUploadedFileRepo - UpdateEmailSent - r.Builder: %w", err)
	}

	// Execute the query using pgx
	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		r.logger.Error("UserUploadedFileRepo - UpdateEmailSent - r.Pool.Exec: failed to execute query", "error", err)
		return fmt.Errorf("UserUploadedFileRepo - UpdateEmailSent - r.Pool.Exec: %w", err)
	}

	r.logger.Info("UserUploadedFileRepo - UpdateEmailSent: successfully updated user uploaded file", "userUploadedFileID", ID)
	return nil
}
