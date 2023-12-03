package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bgg/go-file-gate/internal/usecase"
	"github.com/bgg/go-file-gate/internal/usecase/repo"
)

func TestUserProfileRoute_Create(t *testing.T) {

	pg, dbTeardown := setupDatabase(t)
	defer dbTeardown()

	createTableSQL := `CREATE TABLE user_profiles (
		user_id VARCHAR(255) PRIMARY KEY,
		display_name VARCHAR(255) NOT NULL,
		picture_url VARCHAR(255),
		access_token VARCHAR(255) NOT NULL,
		refresh_token VARCHAR(255) NOT NULL
	);`

	if _, err := pg.Pool.Exec(context.Background(), createTableSQL); err != nil {
		t.Fatalf("could not create user_profiles table: %s", err)
	}

	userProfileUseCase := usecase.NewUserProfileUseCase(repo.NewUserProfileRepo(pg))

	router, redisTeardown := setupRouter(t)
	defer redisTeardown()
	l := setupLogger(t)

	NewUserProfileRoutes(router.Group("/api/v1"), userProfileUseCase, l)

	sessionCookie := setupSessions(t, router)

	payload := map[string]interface{}{
		"UserID":      "testuser",
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
}
