package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bgg/go-flow-gateway/internal/infra/repo"
	"github.com/bgg/go-flow-gateway/internal/usecase"
	"github.com/bgg/go-flow-gateway/pkg/postgres"
	"github.com/gin-gonic/gin"
)

func setupUserProfilesTable(t *testing.T) (*postgres.Postgres, func()) {
	t.Helper()

	pg, dbTeardown := setupDatabase(t)

	createTableSQL := `CREATE TABLE user_profiles (
		user_id SERIAL PRIMARY KEY,
		display_name VARCHAR(255) NOT NULL,
		picture_url VARCHAR(255)
	);`

	if _, err := pg.Pool.Exec(context.Background(), createTableSQL); err != nil {
		t.Fatalf("could not create user_profiles table: %s", err)
	}

	return pg, dbTeardown
}

func setupUserProfileRoutes(t *testing.T) (*gin.Engine, *http.Cookie, *postgres.Postgres, func()) {
	t.Helper()

	l := setupLogger(t)

	pg, dbTeardown := setupUserProfilesTable(t)

	userProfileUseCase := usecase.NewUserProfileUseCase(repo.NewUserProfileRepo(pg, l), l)

	router, redisTeardown := setupRouter(t)

	NewUserProfileRoutes(router.Group("/api/v1"), userProfileUseCase, l)

	sessionCookie := setupSessions(t, router)
	return router, sessionCookie, pg, func() {
		dbTeardown()
		redisTeardown()
	}

}

func TestUserProfileRoute_Create(t *testing.T) {

	router, sessionCookie, _, teardown := setupUserProfileRoutes(t)
	defer teardown()

	const (
		url         = "/api/v1/user-profiles/"
		httpMethod  = "POST"
		contentType = "application/json"
		displayName = "Test User"
		pictureURL  = "https://test.com/test.jpg"
	)

	t.Run("create user profile successfully", func(t *testing.T) {

		payload := map[string]interface{}{
			"DisplayName": displayName,
			"PictureURL":  pictureURL,
		}

		// create a actual request with session cookie
		jsonPayload, _ := json.Marshal(payload)
		req, _ := http.NewRequest(httpMethod, url, bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", contentType)
		req.AddCookie(sessionCookie)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("create user profile with invalid payload", func(t *testing.T) {

		payload := map[string]interface{}{
			"DisplayName": "",
			"PictureURL":  "",
		}

		// create a actual request with session cookie
		jsonPayload, _ := json.Marshal(payload)
		req, _ := http.NewRequest(httpMethod, url, bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", contentType)
		req.AddCookie(sessionCookie)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

}

func TestUserProfileRoute_Get(t *testing.T) {
	const (
		userId      = 1
		displayName = "Test User"
		pictureURL  = "https://test.com/test.jpg"
		httpMethod  = "GET"
		contentType = "application/json"
	)

	router, sessionCookie, pg, teardown := setupUserProfileRoutes(t)
	defer teardown()

	query := `INSERT INTO user_profiles (user_id, display_name, picture_url) VALUES ($1, $2,$3);`
	_, err := pg.Pool.Exec(context.Background(), query, userId, displayName, pictureURL)
	if err != nil {
		t.Fatalf("could not insert test data: %s", err)
	}

	t.Run("get user profile successfully", func(t *testing.T) {

		url := fmt.Sprintf("/api/v1/user-profiles/%d", userId)
		// create a actual request with session cookie
		req, _ := http.NewRequest(httpMethod, url, nil)
		req.Header.Set("Content-Type", contentType)
		req.AddCookie(sessionCookie)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("get user profile with invalid user id", func(t *testing.T) {

		url := fmt.Sprintf("/api/v1/user-profiles/%d", 2)
		// create a actual request with session cookie
		req, _ := http.NewRequest(httpMethod, url, nil)
		req.Header.Set("Content-Type", contentType)
		req.AddCookie(sessionCookie)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})

}
