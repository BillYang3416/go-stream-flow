package app

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/bgg/go-flow-gateway/config"
	"github.com/bgg/go-flow-gateway/internal/adapter/event"
	v1 "github.com/bgg/go-flow-gateway/internal/adapter/rest/v1"
	"github.com/bgg/go-flow-gateway/internal/infra/email"
	"github.com/bgg/go-flow-gateway/internal/infra/messaging/rabbitmq"
	"github.com/bgg/go-flow-gateway/internal/infra/repo"
	"github.com/bgg/go-flow-gateway/internal/infra/utils"
	"github.com/bgg/go-flow-gateway/internal/usecase"
	"github.com/bgg/go-flow-gateway/pkg/logger"
	"github.com/bgg/go-flow-gateway/pkg/postgres"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	mail "github.com/xhit/go-simple-mail/v2"
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

	mailHogPort, err := strconv.Atoi(cfg.MailHog.Port)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - strconv.Atoi: %w", err))
	}
	// SMTP Client
	server := mail.NewSMTPClient()
	server.Host = cfg.MailHog.Host
	server.Port = mailHogPort
	server.Username = ""
	server.Password = ""
	server.Encryption = mail.EncryptionNone

	smtpClient, err := server.Connect()
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - server.Connect: %w", err))
	}
	defer smtpClient.Close()

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

	userUploadedFileCase := usecase.NewUserUploadedFileUseCase(
		repo.NewUserUploadedFileRepo(pg),
		rabbitmq.NewUserUploadedFilePublisher(l, ch),
		email.NewUserUploadedFileEmailSender(smtpClient, l),
	)
	// Consumer
	cs := event.NewUserUploadedFileConsumer(userUploadedFileCase, ch, l)
	go cs.StartConsume()

	// Use case
	userProfileUseCase := usecase.NewUserProfileUseCase(
		repo.NewUserProfileRepo(pg),
	)
	oauthDetailUseCase := usecase.NewOAuthDetailUseCase(
		repo.NewOAuthDetailRepo(pg),
	)
	userCredentialUseCase := usecase.NewUserCredentialUseCase(
		repo.NewUserCredentialRepo(pg),
		utils.NewBcryptHasher(),
	)

	// HTTP Server
	v1.NewRouter(cfg, handler, l, userProfileUseCase, userUploadedFileCase, oauthDetailUseCase, userCredentialUseCase)
	httpServer := &http.Server{
		Addr:    ":" + cfg.HTTP.Port,
		Handler: handler,
	}

	l.Info("Starting server on port %s\n", cfg.HTTP.Port)
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		l.Fatal(fmt.Errorf("app - Run - httpServer.ListenAndServe: %w", err))
	}
}
