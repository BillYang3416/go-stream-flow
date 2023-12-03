package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bgg/go-file-gate/internal/usecase"
	"github.com/bgg/go-file-gate/internal/usecase/repo"
	"github.com/bgg/go-file-gate/pkg/logger"
	"github.com/bgg/go-file-gate/pkg/postgres"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/ory/dockertest/v3"
)

func TestUserProfileRoute_Create(t *testing.T) {

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("could not connect to docker: %s", err)
	}

	dbResource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "alpine",
		Env:        []string{"POSTGRES_PASSWORD=secret"},
	})
	if err != nil {
		t.Fatalf("could not start db resource: %s", err)
	}
	defer pool.Purge(dbResource)

	dbPort := dbResource.GetPort("5432/tcp")
	dbURL := fmt.Sprintf("postgres://postgres:secret@localhost:%s/postgres?sslmode=disable", dbPort)
	var pg *postgres.Postgres
	err = pool.Retry(func() error {
		pg, err = postgres.New(dbURL, postgres.MaxPoolSize(1))
		if err != nil {
			return err
		}
		return pg.Ping(context.Background())
	})
	if err != nil {
		t.Fatalf("could not connnect to dockerized postgres: %s", err)
	}
	defer pg.Close()

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

	redisResource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "redis",
		Tag:        "alpine",
	})
	if err != nil {
		t.Fatalf("could not start resource: %s", err)
	}
	defer pool.Purge(redisResource)

	err = pool.Retry(func() error {
		conn, err := redigo.Dial("tcp", fmt.Sprintf("localhost:%s", redisResource.GetPort("6379/tcp")))
		if err != nil {
			return err
		}
		conn.Close()
		return nil
	})
	if err != nil {
		t.Fatalf("could not connect to redis: %s", err)
	}

	// // Optional: Add a small delay to ensure Redis is fully ready
	// time.Sleep(1 * time.Second)

	store, err := redis.NewStore(10, "tcp", fmt.Sprintf("localhost:%s", redisResource.GetPort("6379/tcp")), "", []byte("secret"))
	if err != nil {
		t.Fatalf("could not create redis store: %s", err)
	}

	router := gin.Default()
	router.Use(sessions.Sessions("user-auth", store))

	l := logger.New("debug")

	NewUserProfileRoutes(router.Group("/api/v1"), userProfileUseCase, l)

	// define a route to set session for testing
	router.GET("/set-session", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("userID", "testuser")
		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/set-session", nil)
	router.ServeHTTP(w, req)

	// extract the session cookie from the response
	var sessionCookie *http.Cookie
	for _, cookie := range w.Result().Cookies() {
		if cookie.Name == "user-auth" {
			sessionCookie = cookie
			break
		}
	}
	if sessionCookie == nil {
		t.Fatalf("could not find session cookie")
	}

	payload := map[string]interface{}{
		"UserID":      "testuser",
		"DisplayName": "Test User",
		"PictureURL":  "https://test.com/test.jpg",
	}

	// create a actual request with session cookie
	jsonPayload, _ := json.Marshal(payload)
	req, _ = http.NewRequest("POST", "/api/v1/user-profiles/", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(sessionCookie)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	if err := pool.Purge(dbResource); err != nil {
		t.Fatalf("could not purge resource: %s", err)
	}
}
