package v1

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bgg/go-flow-gateway/internal/infra/email"
	"github.com/bgg/go-flow-gateway/internal/infra/messaging/rabbitmq"
	"github.com/bgg/go-flow-gateway/internal/infra/repo"
	"github.com/bgg/go-flow-gateway/internal/usecase"
	"github.com/bgg/go-flow-gateway/pkg/postgres"
	"github.com/gin-gonic/gin"
)

func setupUserUploadedFilesTable(t *testing.T) (*postgres.Postgres, func()) {
	t.Helper()

	pg, dbTeardown := setupUserProfilesTable(t)

	// insert test data
	_, err := pg.Pool.Exec(context.Background(), `INSERT INTO user_profiles (display_name, picture_url) VALUES ('Test User', 'https://test.com/test.jpg');`)
	if err != nil {
		t.Fatalf("could not insert test data: %s", err)
	}

	createTableSQL := `CREATE TABLE user_uploaded_files (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		size BIGINT NOT NULL,
		content BYTEA NOT NULL,
		user_id INT NOT NULL REFERENCES user_profiles(user_id),
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

func setupUserUploadedFileRoute(t *testing.T) (*gin.Engine, *http.Cookie, *postgres.Postgres, func()) {
	pg, dbTeardown := setupUserUploadedFilesTable(t)

	ch, rabbitMQTeardown := setupRabbitMQ(t)

	smtpClient, smtpTeardown := setupMailhog(t)

	l := setupLogger(t)

	userUploadedFileUseCase := usecase.NewUserUploadedFileUseCase(repo.NewUserUploadedFileRepo(pg, l), rabbitmq.NewUserUploadedFilePublisher(l, ch), email.NewUserUploadedFileEmailSender(smtpClient, l), l)

	router, redisTeardown := setupRouter(t)

	NewUserUploadedFileRoutes(router.Group("/api/v1"), userUploadedFileUseCase, l)

	sessionCookie := setupSessions(t, router)
	return router, sessionCookie, pg, func() {
		dbTeardown()
		redisTeardown()
		rabbitMQTeardown()
		smtpTeardown()
	}
}

func TestUserUploadedFileRoute_Create(t *testing.T) {

	router, sessionCookie, _, teardown := setupUserUploadedFileRoute(t)
	defer teardown()

	const (
		url            = "/api/v1/user-uploaded-files/"
		httpMethod     = "POST"
		emailRecipient = "johndoe@email.com"
		fileName       = "test.txt"
		fileContent    = "dummy file content"
	)

	t.Run("create user uploaded file successfully", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		_ = writer.WriteField("emailRecipient", emailRecipient)

		fileContent := bytes.NewBufferString(fileContent)
		part, err := writer.CreateFormFile("file", fileName)
		if err != nil {
			t.Fatalf("could not create form file: %s", err)
		}
		_, err = io.Copy(part, fileContent)
		if err != nil {
			t.Fatalf("could not copy file content: %s", err)
		}

		writer.Close()

		req, err := http.NewRequest(httpMethod, url, body)
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

		fileContent := bytes.NewBufferString(fileContent)
		part, err := writer.CreateFormFile("file", fileName)
		if err != nil {
			t.Fatalf("could not create form file: %s", err)
		}
		_, err = io.Copy(part, fileContent)
		if err != nil {
			t.Fatalf("could not copy file content: %s", err)
		}

		writer.Close()

		req, err := http.NewRequest(httpMethod, url, body)
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
