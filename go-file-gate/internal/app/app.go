package app

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bgg/go-file-gate/config"
	v1 "github.com/bgg/go-file-gate/internal/controllers/http/v1"
	"github.com/bgg/go-file-gate/internal/usecase"
	"github.com/bgg/go-file-gate/internal/usecase/repo"
	"github.com/bgg/go-file-gate/pkg/logger"
	"github.com/bgg/go-file-gate/pkg/postgres"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
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

	// Use case
	userProfileUseCase := usecase.NewUserProfileUseCase(
		repo.NewUserProfileRepo(pg),
	)
	// HTTP Server
	v1.NewRouter(handler, l, userProfileUseCase)
	httpServer := &http.Server{
		Addr:    ":" + cfg.HTTP.Port,
		Handler: handler,
	}

	l.Info("Starting server on port %s\n", cfg.HTTP.Port)
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		l.Fatal(fmt.Errorf("app - Run - httpServer.ListenAndServe: %w", err))
	}
}