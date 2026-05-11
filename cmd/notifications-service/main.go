package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gin-gonic/gin"
	"github.com/sony/gobreaker"
	"github.com/ua-academy-projects/share-bite/internal/config"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
	notificationconsumer "github.com/ua-academy-projects/share-bite/internal/notification/consumer"
	notificationhandler "github.com/ua-academy-projects/share-bite/internal/notification/handler"
	notificationrepo "github.com/ua-academy-projects/share-bite/internal/notification/repository"
	notificationservice "github.com/ua-academy-projects/share-bite/internal/notification/service"
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
	notificationHub := notification.NewHub(broker)

	notificationRepo := notificationrepo.New(client.DB())
	notificationService := notificationservice.New(notificationRepo, broker)

	tokenManager := jwt.NewTokenManager(
		config.Config().JwtToken.AccessTokenSecretKey(),
		config.Config().JwtToken.RefreshTokenSecretKey(),
		config.Config().JwtToken.AccessTokenTTL(),
		config.Config().JwtToken.RefreshTokenTTL(),
	)
	authMiddleware := middleware.Auth(tokenManager)

	router := gin.New()
	router.Use(gin.Recovery(), ErrorMiddleware())
	notificationhandler.RegisterHandlers(router.Group("/"), notificationService, notificationHub, authMiddleware)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	serverAddr := config.Config().NotificationHttpServer.Address()
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
			logger.Fatal(ctx, "notifications http server error:", err)
		}
	}()

	// SQS Configuration
	sqsCfg := config.Config().NotificationSQS
	region := sqsCfg.Region()
	if region == "" {
		logger.Fatal(ctx, "NOTIFICATION_AWS_REGION is required")
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(region))
	if err != nil {
		logger.Fatal(ctx, "load aws config:", err)
	}

	sqsClient := sqs.NewFromConfig(awsCfg, func(o *sqs.Options) {
		if endpoint := sqsCfg.Endpoint(); endpoint != "" {
			o.BaseEndpoint = &endpoint
		}
	})

	queueRef := sqsCfg.Queue()
	if queueRef != "" {
		sqsQueueURL := queueRef
		if strings.HasPrefix(queueRef, "arn:aws:sqs") {
			parsedARN, err := arn.Parse(queueRef)
			if err != nil {
				logger.Fatal(ctx, "invalid SQS ARN:", err)
			}

			_, queueName, _ := strings.Cut(parsedARN.Resource, ":")
			if queueName == "" {
				queueName = parsedARN.Resource
			}

			out, err := sqsClient.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
				QueueName:              aws.String(queueName),
				QueueOwnerAWSAccountId: aws.String(parsedARN.AccountID),
			})
			if err != nil {
				logger.Fatal(ctx, "failed to resolve SQS queue URL from ARN:", err)
			}
			sqsQueueURL = *out.QueueUrl
		}

		processor := notificationservice.NewProcessor(notificationService)
		consumer := notificationconsumer.New(sqsClient, sqsQueueURL, processor)

		go func() {
			logger.InfoKV(ctx, "notifications sqs consumer starting", "queue_url", sqsQueueURL)
			if err := consumer.Run(ctx); err != nil && err != context.Canceled {
				logger.ErrorKV(ctx, "notifications sqs consumer stopped", "error", err)
			}
		}()
	} else {
		logger.Warn(ctx, "SQS queue not provided")
	}

	closer.Wait()
}

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()
		if err == nil {
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
