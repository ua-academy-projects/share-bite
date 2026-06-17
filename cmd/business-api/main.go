package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sony/gobreaker"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/ua-academy-projects/share-bite/internal/config/env"
	"github.com/ua-academy-projects/share-bite/pkg/notification"
	"github.com/ua-academy-projects/share-bite/pkg/redis"
	"github.com/ua-academy-projects/share-bite/pkg/resilience"

	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/internal/business/error/code"
	"github.com/ua-academy-projects/share-bite/internal/business/handler/business"
	notifhandler "github.com/ua-academy-projects/share-bite/internal/business/handler/notification"
	businessrepo "github.com/ua-academy-projects/share-bite/internal/business/repository/business"
	businesssvc "github.com/ua-academy-projects/share-bite/internal/business/service/business"
	"github.com/ua-academy-projects/share-bite/internal/config"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
	"github.com/ua-academy-projects/share-bite/internal/storage"
	"github.com/ua-academy-projects/share-bite/pkg/closer"
	"github.com/ua-academy-projects/share-bite/pkg/database/pg"
	"github.com/ua-academy-projects/share-bite/pkg/database/txmanager"
	admingateway "github.com/ua-academy-projects/share-bite/pkg/gateway/admin"
	"github.com/ua-academy-projects/share-bite/pkg/outbox"
	"github.com/ua-academy-projects/share-bite/pkg/resilience"

	_ "github.com/ua-academy-projects/share-bite/docs/api/business"
	h3 "github.com/ua-academy-projects/share-bite/pkg/aws"
	"github.com/ua-academy-projects/share-bite/pkg/jwt"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	"go.uber.org/zap"
)

// @title                  ShareBite Business API
// @version             1.0
// @description          API for discovering brand locations (venues).
//
// @securityDefinitions.apikey  BearerAuth
// @in                    header
// @name                   Authorization
//
// @BasePath                /
func main() {
	ctx := context.Background()

	if err := config.Load(".env"); err != nil {
		logger.Fatal(ctx, err)
	}

	cfg := config.Config()

	redisCfg, err := env.NewRedisConfig()
	if err != nil {
		logger.Fatal(ctx, "load redis config: ", err)
	}

	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.BusinessHttpServer.AllowedOrigins(),
		AllowMethods:     cfg.BusinessHttpServer.AllowedMethods(),
		AllowHeaders:     cfg.BusinessHttpServer.AllowedHeaders(),
		ExposeHeaders:    cfg.BusinessHttpServer.ExposeHeaders(),
		AllowCredentials: true,
	}))
	router.Use(gin.Recovery())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("/swagger/doc.json"),
	))

	router.Use(ErrorMiddleware())

	if config.Config().App.IsProd() {
		logger.SetLevel(zap.InfoLevel)
		gin.SetMode(gin.ReleaseMode)
	} else {
		logger.SetLevel(zap.DebugLevel)
	}
	closer.SetShutdownTimeout(config.Config().App.GracefulShutdownTimeout())

	client, err := pg.NewClient(ctx, config.Config().Postgres.Dsn())
	if err != nil {
		logger.Fatal(ctx, "new database client: ", err)
	}
	if err := client.DB().Ping(ctx); err != nil {
		logger.Fatal(ctx, "ping database: ", err)
	}
	closer.Add(func(ctx context.Context) error {
		client.Close()
		return nil
	})

	rdb, err := redis.NewClient(
		redisCfg.Addr(),
		redisCfg.Password(),
		redisCfg.DB(),
		redisCfg.TLS(),
	)
	if err != nil {
		logger.Fatal(ctx, "new redis client: ", err)
	}

	ctxPing, cancelPing := context.WithTimeout(ctx, 3*time.Second)
	_, err = rdb.Ping(ctxPing).Result()
	cancelPing()
	if err != nil {
		logger.Fatal(ctx, "ping redis: ", err)
	}
	closer.Add(func(ctx context.Context) error {
		return rdb.Close()
	})

	notificationResiliencePolicy := resilience.Policy{
		RetryConfig: resilience.RetryConfig{
			InitialInterval:     10 * time.Millisecond,
			RandomizationFactor: 0.2,
			Multiplier:          2.0,
			MaxInterval:         200 * time.Millisecond,
			MaxElapsedTime:      1500 * time.Millisecond,
		},
		Breaker: resilience.NewCircuitBreaker(resilience.CircuitBreakerConfig{
			Name:        "business-redis-sub",
			MaxRequests: 1,
			Interval:    10 * time.Second,
			Timeout:     5 * time.Second,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				return counts.ConsecutiveFailures >= 20
			},
			OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
				logger.WarnKV(ctx, "redis circuit breaker state changed", "name", name, "from", from.String(), "to", to.String())
			},
			IsSuccessful: func(err error) bool {
				if err == nil || errors.Is(err, context.Canceled) {
					return true
				}
				return redis.IsPermanentRedisError(err)
			},
		}),
		RetryNotify: func(err error, nextRetryIn time.Duration) {
			logger.Debugf(ctx, "redis subscribe/publish retry scheduled in %v: %v", nextRetryIn, err)
		},
	}

	broker := notification.NewBroker(rdb, notification.WithPublishPolicy(notificationResiliencePolicy))
	notifHub := notification.NewHub(broker)

	txManager := txmanager.NewTransactionManager(client.DB())

	storageClient, err := storage.NewStorageClient(ctx, config.Config().Storage)
	if err != nil {
		logger.Fatal(ctx, "init storage client:", err)
	}

	// repos
	businessRepo := businessrepo.New(client)

	// services
	h3Service := h3.NewH3Service()
	h3Settings := businesssvc.H3Settings{
		Resolution:      config.Config().H3.Resolution(),
		RecommendRadius: config.Config().H3.RecommendRadius(),
	}
	outboxWriter := outbox.NewWriter(client.DB())

	adminResiliencePolicy := resilience.Policy{
		RetryConfig: resilience.RetryConfig{
			InitialInterval:     200 * time.Millisecond,
			RandomizationFactor: 0.25,
			Multiplier:          2,
			MaxInterval:         2 * time.Second,
			MaxElapsedTime:      8 * time.Second,
		},
	}

	adminGateway := admingateway.New(
		config.Config().AdminHttpServer.Address(),
		"/",
		"http",
		&adminResiliencePolicy,
	)

	businessSvc := businesssvc.New(businessRepo, txManager, storageClient, h3Service, h3Settings, outboxWriter, adminGateway)

	tokenManager := jwt.NewTokenManager(
		config.Config().JwtToken.AccessTokenSecretKey(),
		config.Config().JwtToken.RefreshTokenSecretKey(),
		config.Config().JwtToken.AccessTokenTTL(),
		config.Config().JwtToken.RefreshTokenTTL(),
	)

	authMw := middleware.Auth(tokenManager)

	// handlers
	business.RegisterHandlers(router.Group("/business"), businessSvc, tokenManager, storageClient)
	notifhandler.RegisterHandlers(router.Group("/business"), notifHub, authMw)

	go func() {
		logger.Info(ctx, "business http server is running")
		if err := router.Run(config.Config().BusinessHttpServer.Address()); err != nil && err != http.ErrServerClosed {
			logger.Fatal(ctx, "run http server: ", err)
		}
	}()

	closer.Wait()
}

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()
		if err == nil {
			return
		}

		ctx := c.Request.Context()

		var appErr *apperror.Error
		if errors.As(err.Err, &appErr) {
			switch appErr.Code {
			case code.NotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": appErr.Error()})
				return
			case code.BadRequest:
				c.JSON(http.StatusBadRequest, gin.H{"error": appErr.Error()})
				return
			case code.Forbidden:
				c.JSON(http.StatusForbidden, gin.H{"error": appErr.Error()})
				return
			case code.Unauthorized:
				c.JSON(http.StatusUnauthorized, gin.H{"error": appErr.Error()})
				return
			case code.Conflict:
				c.JSON(http.StatusConflict, gin.H{"error": appErr.Error()})
				return
			}
		}

		logger.ErrorKV(ctx, "internal error", "error", err.Err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
