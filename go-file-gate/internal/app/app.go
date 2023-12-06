package app

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bgg/go-file-gate/config"
	v1 "github.com/bgg/go-file-gate/internal/adapter/rest/v1"
	"github.com/bgg/go-file-gate/internal/infra/messaging/rabbitmq"
	"github.com/bgg/go-file-gate/internal/infra/repo"
	"github.com/bgg/go-file-gate/internal/usecase"
	"github.com/bgg/go-file-gate/pkg/logger"
	"github.com/bgg/go-file-gate/pkg/postgres"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Run(cfg *config.Config) {

	l := logger.New(cfg.Log.Level)

	// PostgreSQL Repository
	pg, err := postgres.New(cfg.Postgres.URL, postgres.MaxPoolSize(cfg.Postgres.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	handler := gin.New()
	// Redis Session
	store, err := redis.NewStore(10, "tcp", cfg.Redis.Host+":"+cfg.Redis.Port, cfg.Redis.Password, []byte(os.Getenv("SESSION_SECRET")))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - redis.NewStore: %w", err))
	}
	handler.Use(sessions.Sessions("user-auth", store))

	// RabbitMQ
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.RabbitMQ.Username,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port)
	conn, err := amqp.Dial(url)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - amqp.Dial: %w", err))
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - conn.Channel: %w", err))
	}
	defer ch.Close()

	// Use case
	userProfileUseCase := usecase.NewUserProfileUseCase(
		repo.NewUserProfileRepo(pg),
	)
	userUploadedFileCase := usecase.NewUserUploadedFileUseCase(
		repo.NewUserUploadedFileRepo(pg),
		rabbitmq.NewUserUploadedFilePublisher(ch),
	)

	// HTTP Server
	v1.NewRouter(cfg, handler, l, userProfileUseCase, userUploadedFileCase)
	httpServer := &http.Server{
		Addr:    ":" + cfg.HTTP.Port,
		Handler: handler,
	}

	l.Info("Starting server on port %s\n", cfg.HTTP.Port)
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		l.Fatal(fmt.Errorf("app - Run - httpServer.ListenAndServe: %w", err))
	}
}
