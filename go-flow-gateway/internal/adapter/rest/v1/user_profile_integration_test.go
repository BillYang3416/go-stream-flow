package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bgg/go-flow-gateway/internal/infra/repo"
	"github.com/bgg/go-flow-gateway/internal/usecase"
	"github.com/bgg/go-flow-gateway/pkg/postgres"
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

func TestUserProfileRoute_Create(t *testing.T) {
	l := setupLogger(t)

	pg, dbTeardown := setupUserProfilesTable(t)
	defer dbTeardown()

	userProfileUseCase := usecase.NewUserProfileUseCase(repo.NewUserProfileRepo(pg, l), l)

	router, redisTeardown := setupRouter(t)
	defer redisTeardown()

	NewUserProfileRoutes(router.Group("/api/v1"), userProfileUseCase, l)

	sessionCookie := setupSessions(t, router)

	t.Run("create user profile successfully", func(t *testing.T) {

		payload := map[string]interface{}{
			"DisplayName": "Test User",
			"PictureURL":  "https://test.com/test.jpg",
		}

		// create a actual request with session cookie
		jsonPayload, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/api/v1/user-profiles/", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(sessionCookie)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("create user profile with invalid payload", func(t *testing.T) {

		payload := map[string]interface{}{
			"DisplayName": "Test User",
			"PictureURL":  "invalid-url",
		}

		// create a actual request with session cookie
		jsonPayload, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/api/v1/user-profiles/", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(sessionCookie)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

}

func TestUserProfileRoute_Get(t *testing.T) {
	l := setupLogger(t)

	pg, dbTeardown := setupUserProfilesTable(t)
	defer dbTeardown()

	// insert test data
	_, err := pg.Pool.Exec(context.Background(), `INSERT INTO user_profiles (user_id, display_name, picture_url, access_token, refresh_token) VALUES ('testuser', 'Test User', 'https://test.com/test.jpg', 'test-access-token', 'test-refresh-token');`)
	if err != nil {
		t.Fatalf("could not insert test data: %s", err)
	}

	userProfileUseCase := usecase.NewUserProfileUseCase(repo.NewUserProfileRepo(pg, l), l)

	router, redisTeardown := setupRouter(t)
	defer redisTeardown()

	NewUserProfileRoutes(router.Group("/api/v1"), userProfileUseCase, l)

	sessionCookie := setupSessions(t, router)

	t.Run("get user profile successfully", func(t *testing.T) {
		// create a actual request with session cookie
		req, _ := http.NewRequest("GET", "/api/v1/user-profiles/testuser", nil)
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(sessionCookie)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("get user profile with invalid user id", func(t *testing.T) {
		// create a actual request with session cookie
		req, _ := http.NewRequest("GET", "/api/v1/user-profiles/invalid-user-id", nil)
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(sessionCookie)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d, got %d", http.StatusNotFound, w.Code)
		}
	})

}
