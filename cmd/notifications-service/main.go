package main

import (
	"context"
	"fmt"
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
	hub := notification.NewHub(broker)

	notifRepo := notificationrepo.New(client.DB())
	notifSvc := notificationservice.New(notifRepo, broker)

	tokenManager := jwt.NewTokenManager(
		config.Config().JwtToken.AccessTokenSecretKey(),
		config.Config().JwtToken.RefreshTokenSecretKey(),
		config.Config().JwtToken.AccessTokenTTL(),
		config.Config().JwtToken.RefreshTokenTTL(),
	)
	authMw := middleware.Auth(tokenManager)

	router := gin.New()
	router.Use(gin.Recovery())

	notificationhandler.RegisterHandlers(router.Group("/notifications"), notifSvc, hub, authMw, streamAuthMiddleware(tokenManager))

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

	awsCfg, err := loadSQSConfig(ctx)
	if err != nil {
		logger.Fatal(ctx, "load aws config:", err)
	}
	sqsClient := sqs.NewFromConfig(awsCfg)

	queueRef := os.Getenv("NOTIFICATION_SQS_QUEUE_ARN")
	if queueRef == "" {
		queueRef = os.Getenv("NOTIFICATION_SQS_QUEUE_URL")
	}
	if queueRef == "" {
		logger.Fatal(ctx, "NOTIFICATION_SQS_QUEUE_ARN or NOTIFICATION_SQS_QUEUE_URL is required")
	}

	sqsQueueURL, err := resolveQueueURL(ctx, sqsClient, queueRef)
	if err != nil {
		logger.ErrorKV(ctx, "failed to resolve sqs queue url, consumer will not start", "error", err)
	} else {
		processor := notificationservice.NewProcessor(notifSvc)
		consumer := notificationconsumer.New(sqsClient, sqsQueueURL, processor)

		go func() {
			logger.InfoKV(ctx, "notifications sqs consumer starting", "queue_url", sqsQueueURL)
			if err := consumer.Run(ctx); err != nil && err != context.Canceled {
				logger.ErrorKV(ctx, "notifications sqs consumer stopped", "error", err)
			}
		}()
	}

	closer.Wait()
}

func resolveQueueURL(ctx context.Context, client *sqs.Client, queueRef string) (string, error) {
	if strings.HasPrefix(queueRef, "https://") || strings.HasPrefix(queueRef, "http://") {
		return queueRef, nil
	}

	parsedARN, err := arn.Parse(queueRef)
	if err != nil {
		return "", fmt.Errorf("invalid queue URL or ARN: %w", err)
	}

	// SQS ARN format: arn:aws:sqs:region:account-id:queue-name
	parts := strings.Split(parsedARN.Resource, ":")
	queueName := parts[len(parts)-1]

	out, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName:              aws.String(queueName),
		QueueOwnerAWSAccountId: aws.String(parsedARN.AccountID),
	})
	if err != nil {
		return "", fmt.Errorf("get queue URL from ARN: %w", err)
	}

	if out.QueueUrl == nil {
		return "", fmt.Errorf("get queue URL returned nil")
	}

	return *out.QueueUrl, nil
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

func loadSQSConfig(ctx context.Context) (aws.Config, error) {
	region := os.Getenv("NOTIFICATION_AWS_REGION")
	if region == "" {
		region = os.Getenv("AWS_REGION")
	}
	if region == "" {
		region = os.Getenv("AWS_DEFAULT_REGION")
	}
	if region == "" {
		region = "us-east-2"
	}

	opts := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(region),
	}

	if endpoint := os.Getenv("NOTIFICATION_SQS_ENDPOINT_URL"); endpoint != "" {
		opts = append(opts, awsconfig.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
				if !strings.EqualFold(service, sqs.ServiceID) {
					return aws.Endpoint{}, &aws.EndpointNotFoundError{}
				}
				return aws.Endpoint{URL: endpoint, HostnameImmutable: true}, nil
			}),
		))
	}

	return awsconfig.LoadDefaultConfig(ctx, opts...)
}

func streamAuthMiddleware(parser interface {
	ParseAccessToken(string) (string, string, error)
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.TrimSpace(c.GetHeader("Authorization"))
		if token != "" {
			if !strings.HasPrefix(token, "Bearer ") {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header"})
				return
			}
			token = strings.TrimSpace(strings.TrimPrefix(token, "Bearer "))
		}
		if token == "" {
			token = strings.TrimSpace(c.Query("access_token"))
		}
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing access token"})
			return
		}

		userID, role, err := parser.ParseAccessToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set(middleware.CtxUserID, userID)
		c.Set(middleware.CtxUserRole, role)
		c.Next()
	}
}
