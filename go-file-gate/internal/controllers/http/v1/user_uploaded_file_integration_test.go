package v1

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bgg/go-file-gate/internal/infra/messaging/rabbitmq"
	"github.com/bgg/go-file-gate/internal/infra/repo"
	"github.com/bgg/go-file-gate/internal/usecase"
	"github.com/bgg/go-file-gate/pkg/postgres"
)

func setupUserUploadedFilesTable(t *testing.T) (*postgres.Postgres, func()) {
	t.Helper()

	pg, dbTeardown := setupUserProfilesTable(t)

	// insert test data
	_, err := pg.Pool.Exec(context.Background(), `INSERT INTO user_profiles (user_id, display_name, picture_url, access_token, refresh_token) VALUES ('testuser', 'Test User', 'https://test.com/test.jpg', 'test-access-token', 'test-refresh-token');`)
	if err != nil {
		t.Fatalf("could not insert test data: %s", err)
	}

	createTableSQL := `CREATE TABLE user_uploaded_files (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		size BIGINT NOT NULL,
		content BYTEA NOT NULL,
		user_id VARCHAR(255) NOT NULL REFERENCES user_profiles(user_id),
		created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
		email_sent BOOLEAN NOT NULL,
		email_sent_at TIMESTAMPTZ,
		email_recipient VARCHAR(255),
		error_message TEXT
	);`

	if _, err := pg.Pool.Exec(context.Background(), createTableSQL); err != nil {
		t.Fatalf("could not create user_uploaded_files table: %s", err)
	}

	return pg, dbTeardown
}

func TestUserUploadedFileRoute_Create(t *testing.T) {

	pg, dbTeardown := setupUserUploadedFilesTable(t)
	defer dbTeardown()

	ch, rabbitMQTeardown := setupRabbitMQ(t)
	defer rabbitMQTeardown()

	userUploadedFileUseCase := usecase.NewUserUploadedFileUseCase(repo.NewUserUploadedFileRepo(pg), rabbitmq.NewUserUploadedFilePublisher(ch))

	router, redisTeardown := setupRouter(t)
	defer redisTeardown()
	l := setupLogger(t)

	NewUserUploadedFileRoutes(router.Group("/api/v1"), userUploadedFileUseCase, l)

	sessionCookie := setupSessions(t, router)

	t.Run("create user uploaded file successfully", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		_ = writer.WriteField("emailRecipient", "johndoe@email.com")

		fileContent := bytes.NewBufferString("dummy file content")
		part, err := writer.CreateFormFile("file", "test.txt")
		if err != nil {
			t.Fatalf("could not create form file: %s", err)
		}
		_, err = io.Copy(part, fileContent)
		if err != nil {
			t.Fatalf("could not copy file content: %s", err)
		}

		writer.Close()

		req, err := http.NewRequest("POST", "/api/v1/user-uploaded-files/", body)
		if err != nil {
			t.Fatalf("could not create request: %s", err)
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.AddCookie(sessionCookie)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected status code %d, got %d", http.StatusNoContent, w.Code)
		}

	})

	t.Run("create user uploaded file with invalid request body", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		_ = writer.WriteField("emailRecipient", "invalid email")

		fileContent := bytes.NewBufferString("dummy file content")
		part, err := writer.CreateFormFile("file", "test.txt")
		if err != nil {
			t.Fatalf("could not create form file: %s", err)
		}
		_, err = io.Copy(part, fileContent)
		if err != nil {
			t.Fatalf("could not copy file content: %s", err)
		}

		writer.Close()

		req, err := http.NewRequest("POST", "/api/v1/user-uploaded-files/", body)
		if err != nil {
			t.Fatalf("could not create request: %s", err)
		}
		req.AddCookie(sessionCookie)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}
