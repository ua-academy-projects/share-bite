package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gin-gonic/gin"
	"github.com/sony/gobreaker"
	"github.com/ua-academy-projects/share-bite/internal/config"
	guestcustomerrepo "github.com/ua-academy-projects/share-bite/internal/guest/repository/customer"
	guestcustomersvc "github.com/ua-academy-projects/share-bite/internal/guest/service/customer"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
	notificationconsumer "github.com/ua-academy-projects/share-bite/internal/notification/consumer"
	notificationhandler "github.com/ua-academy-projects/share-bite/internal/notification/handler"
	notificationrepo "github.com/ua-academy-projects/share-bite/internal/notification/repository"
	notificationservice "github.com/ua-academy-projects/share-bite/internal/notification/service"
	storageclient "github.com/ua-academy-projects/share-bite/internal/storage"
	"github.com/ua-academy-projects/share-bite/pkg/closer"
	"github.com/ua-academy-projects/share-bite/pkg/database/pg"
	"github.com/ua-academy-projects/share-bite/pkg/jwt"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	"github.com/ua-academy-projects/share-bite/pkg/notification"
	redispkg "github.com/ua-academy-projects/share-bite/pkg/redis"
	"github.com/ua-academy-projects/share-bite/pkg/resilience"
)

func main() {
	baseCtx := context.Background()
	ctx, stop := signal.NotifyContext(baseCtx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := config.Load(".env"); err != nil && !os.IsNotExist(err) {
		logger.Fatal(ctx, "config load:", err)
	}

	closer.SetShutdownTimeout(config.Config().App.GracefulShutdownTimeout())

	client, err := pg.NewClient(ctx, config.Config().Postgres.Dsn())
	if err != nil {
		logger.Fatal(ctx, "new database client:", err)
	}
	if err := client.DB().Ping(ctx); err != nil {
		logger.Fatal(ctx, "ping database:", err)
	}
	closer.Add(func(ctx context.Context) error {
		client.Close()
		return nil
	})

	storageClient, err := storageclient.NewStorageClient(ctx, config.Config().Storage)
	if err != nil {
		logger.Fatal(ctx, "new storage client:", err)
	}

	rdb, err := redispkg.NewClient(
		config.Config().Redis.Addr(),
		config.Config().Redis.Password(),
		config.Config().Redis.DB(),
		config.Config().Redis.TLS(),
	)
	if err != nil {
		logger.Fatal(ctx, "new redis client:", err)
	}
	closer.Add(func(ctx context.Context) error {
		rdb.Close()
		return nil
	})

	broker := notification.NewBroker(rdb, notification.WithPublishPolicy(resilience.Policy{
		RetryConfig: resilience.RetryConfig{
			InitialInterval:     25 * time.Millisecond,
			RandomizationFactor: 0.2,
			Multiplier:          2,
			MaxInterval:         500 * time.Millisecond,
			MaxElapsedTime:      3 * time.Second,
		},
		Breaker: resilience.NewCircuitBreaker(resilience.CircuitBreakerConfig{
			Name:        "notifications-redis-publish",
			MaxRequests: 1,
			Interval:    10 * time.Second,
			Timeout:     5 * time.Second,
			ReadyToTrip: func(counts gobreaker.Counts) bool { return counts.ConsecutiveFailures >= 10 },
		}),
	}))
	hub := notification.NewHub(broker)

	customerRepo := guestcustomerrepo.New(client)
	customerSvc := guestcustomersvc.New(customerRepo)
	notifRepo := notificationrepo.New(client.DB())
	notifSvc := notificationservice.New(notifRepo, customerSvc, broker, storageClient)

	tokenManager := jwt.NewTokenManager(
		config.Config().JwtToken.AccessTokenSecretKey(),
		config.Config().JwtToken.RefreshTokenSecretKey(),
		config.Config().JwtToken.AccessTokenTTL(),
		config.Config().JwtToken.RefreshTokenTTL(),
	)
	authMw := middleware.Auth(tokenManager)

	router := gin.New()
	router.Use(gin.Recovery())

	notificationhandler.RegisterHandlers(router.Group("/notifications"), notifSvc, hub, customerSvc, authMw)

	serverAddr := notificationsServerAddr()
	httpServer := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}
	closer.Add(func(ctx context.Context) error {
		return httpServer.Shutdown(ctx)
	})

	go func() {
		logger.InfoKV(ctx, "notifications http server starting", "addr", serverAddr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(ctx, "notifications http server:", err)
		}
	}()

	sqsQueueURL := os.Getenv("NOTIFICATION_SQS_QUEUE_URL")
	if sqsQueueURL == "" {
		logger.Fatal(ctx, "NOTIFICATION_SQS_QUEUE_URL is required")
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		logger.Fatal(ctx, "load aws config:", err)
	}
	sqsClient := sqs.NewFromConfig(awsCfg)
	processor := notificationservice.NewProcessor(notifSvc)
	consumer := notificationconsumer.New(sqsClient, sqsQueueURL, processor)

	go func() {
		logger.InfoKV(ctx, "notifications sqs consumer starting", "queue_url", sqsQueueURL)
		if err := consumer.Run(ctx); err != nil && err != context.Canceled {
			logger.ErrorKV(ctx, "notifications sqs consumer stopped", "error", err)
		}
	}()

	closer.Wait()
}

func notificationsServerAddr() string {
	host := os.Getenv("NOTIFICATIONS_HTTP_SERVER_HOST")
	if host == "" {
		host = "0.0.0.0"
	}
	port := os.Getenv("NOTIFICATIONS_HTTP_SERVER_PORT")
	if port == "" {
		port = "4005"
	}
	return fmt.Sprintf("%s:%s", host, port)
}
