package v1

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bgg/go-flow-gateway/pkg/logger"
	"github.com/bgg/go-flow-gateway/pkg/postgres"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	redigo "github.com/gomodule/redigo/redis"
	mail "github.com/xhit/go-simple-mail/v2"

	"github.com/ory/dockertest/v3"
	amqp "github.com/rabbitmq/amqp091-go"
)

func setupDatabase(t *testing.T) (*postgres.Postgres, func()) {
	t.Helper()

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

	return pg, func() {
		pg.Close()
		pool.Purge(dbResource)
	}

}

func setupRedis(t *testing.T) (redis.Store, func()) {
	t.Helper()

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("could not connect to docker: %s", err)
	}

	redisResource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "redis",
		Tag:        "alpine",
	})
	if err != nil {
		t.Fatalf("could not start resource: %s", err)
	}

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

	// Optional: Add a small delay to ensure Redis is fully ready
	time.Sleep(1 * time.Second)

	store, err := redis.NewStore(10, "tcp", fmt.Sprintf("localhost:%s", redisResource.GetPort("6379/tcp")), "", []byte("secret"))
	if err != nil {
		t.Fatalf("could not create redis store: %s", err)
	}

	return store, func() {
		pool.Purge(redisResource)
	}
}

func setupRouter(t *testing.T) (*gin.Engine, func()) {
	t.Helper()

	store, redisTeardown := setupRedis(t)

	router := gin.Default()
	router.Use(sessions.Sessions("user-auth", store))

	return router, redisTeardown
}

func setupLogger(t *testing.T) logger.Logger {
	t.Helper()

	l := logger.New("debug")

	return l
}

func setupSessions(t *testing.T, router *gin.Engine) *http.Cookie {
	t.Helper()

	// define a route to set session for testing
	router.GET("/set-session", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("userID", 1)
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

	return sessionCookie
}

func setupRabbitMQ(t *testing.T) (*amqp.Channel, func()) {
	t.Helper()

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("could not connect to docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "rabbitmq",
		Tag:        "3.8.9-alpine",
		Env:        []string{"RABBITMQ_DEFAULT_USER=guest", "RABBITMQ_DEFAULT_PASS=guest"},
	})

	if err != nil {
		t.Fatalf("could not start resource: %s", err)
	}

	var conn *amqp.Connection
	if err = pool.Retry(func() error {
		conn, err = amqp.Dial(fmt.Sprintf("amqp://guest:guest@localhost:%s/", resource.GetPort("5672/tcp")))
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		t.Fatalf("could not connect to dockerized rabbitmq: %s", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		t.Fatalf("could not create channel: %s", err)
	}

	return ch, func() {
		ch.Close()
		conn.Close()
		pool.Purge(resource)
	}
}

func setupMailhog(t *testing.T) (*mail.SMTPClient, func()) {
	t.Helper()

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("could not connect to docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository:   "mailhog/mailhog",
		Tag:          "latest",
		ExposedPorts: []string{"1025", "8025"},
	})
	if err != nil {
		t.Fatalf("could not start resource: %s", err)
	}

	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.Username = ""
	server.Password = ""
	server.Encryption = mail.EncryptionNone

	var smtpClient *mail.SMTPClient
	if err := pool.Retry(func() error {
		smtpClient, err = server.Connect()
		if err != nil {
			return err
		}
		return nil

	}); err != nil {
		t.Fatalf("could not connect to mailhog: %s", err)
	}

	return smtpClient, func() {
		smtpClient.Close()
		pool.Purge(resource)
	}
}
