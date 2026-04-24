package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
	_ "github.com/ua-academy-projects/share-bite/docs/api/guest"
	"github.com/ua-academy-projects/share-bite/internal/config"
	businessgateway "github.com/ua-academy-projects/share-bite/internal/guest/gateway/business"
	"github.com/ua-academy-projects/share-bite/internal/guest/handler/collection"
	"github.com/ua-academy-projects/share-bite/internal/guest/handler/comment"
	"github.com/ua-academy-projects/share-bite/internal/guest/handler/customer"
	"github.com/ua-academy-projects/share-bite/internal/guest/handler/post"
	guest_middleware "github.com/ua-academy-projects/share-bite/internal/guest/middleware"
	collectionrepo "github.com/ua-academy-projects/share-bite/internal/guest/repository/collection"
	commentrepo "github.com/ua-academy-projects/share-bite/internal/guest/repository/comment"
	customerrepo "github.com/ua-academy-projects/share-bite/internal/guest/repository/customer"
	postrepo "github.com/ua-academy-projects/share-bite/internal/guest/repository/post"
	collectionsvc "github.com/ua-academy-projects/share-bite/internal/guest/service/collection"
	commentsvc "github.com/ua-academy-projects/share-bite/internal/guest/service/comment"
	customersvc "github.com/ua-academy-projects/share-bite/internal/guest/service/customer"
	postsvc "github.com/ua-academy-projects/share-bite/internal/guest/service/post"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
	"github.com/ua-academy-projects/share-bite/internal/storage"
	"github.com/ua-academy-projects/share-bite/pkg/closer"
	"github.com/ua-academy-projects/share-bite/pkg/database/pg"
	"github.com/ua-academy-projects/share-bite/pkg/database/txmanager"
	"github.com/ua-academy-projects/share-bite/pkg/jwt"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	common_middleware "github.com/ua-academy-projects/share-bite/pkg/middleware"
	"github.com/ua-academy-projects/share-bite/pkg/notification"
	redis "github.com/ua-academy-projects/share-bite/pkg/redis"
	"github.com/ua-academy-projects/share-bite/pkg/validator"
	"go.uber.org/zap"
)

// @title						Share Bite - Guest Service API
// @version					1.0
// @description				API for the Guest microservice. Manages customer profiles, their posts, collections, comments, likes etc.
//
// @host						localhost:3800
// @BasePath					/
//
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				Type "Bearer " followed by your JWT token.
func main() {
	ctx := context.Background()

	// for local development only
	if err := config.Load(".env"); err != nil {
		logger.Fatal(ctx, "load config:", err)
	}

	// docker variant
	// if err := config.Load(); err != nil {
	// 	logger.Fatal(ctx, "load config:", err)
	// }

	router := gin.New()
	router.Use(common_middleware.RequestID())
	router.Use(common_middleware.RequestLogger())
	router.Use(gin.Recovery())
	router.Use(guest_middleware.ErrorMiddleware())

	router.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))

	binding.Validator = validator.New("binding")

	if config.Config().App.IsProd() {
		logger.SetLevel(zap.InfoLevel)
		gin.SetMode(gin.ReleaseMode)
	} else {
		logger.SetLevel(zap.DebugLevel)
	}
	closer.SetShutdownTimeout(config.Config().App.GracefulShutdownTimeout())

	// db connection
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
	// redis connection
	rdb, err := redis.NewClient(config.Config().Redis.Addr(), config.Config().Redis.Password())
	if err != nil {
		logger.Fatal(ctx, "new redis client: ", err)
	}
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		logger.Fatal(ctx, "ping redis: ", err)
	}
	closer.Add(func(ctx context.Context) error {
		rdb.Close()
		return nil
	})
	// notifications
	broker := notification.NewBroker(rdb)

	notifHub := notification.NewHub(broker)

	// clients
	clientCfg := config.Config().BusinessHttpClient
	httpClient := &http.Client{
		Timeout: clientCfg.Timeout(),
		Transport: &http.Transport{
			MaxIdleConns:        clientCfg.MaxIdleConns(),
			MaxIdleConnsPerHost: clientCfg.MaxIdleConnsPerHost(),
			IdleConnTimeout:     clientCfg.IdleConnTimeout(),
		},
	}
	closer.Add(func(ctx context.Context) error {
		httpClient.CloseIdleConnections()
		return nil
	})

	businessGateway, err := businessgateway.NewBusinessAPIClient(clientCfg.BaseURL(), "/", httpClient)
	if err != nil {
		logger.Fatalf(ctx, "init business gateway: %v", err)
	}

	storageClient, err := storage.NewStorageClient(ctx, config.Config().Storage)
	if err != nil {
		logger.Fatal(ctx, "init storage client:", err)
	}

	txManager := txmanager.NewTransactionManager(client.DB())

	tokenManager := jwt.NewTokenManager(
		config.Config().JwtToken.AccessTokenSecretKey(),
		config.Config().JwtToken.RefreshTokenSecretKey(),
		config.Config().JwtToken.AccessTokenTTL(),
		config.Config().JwtToken.RefreshTokenTTL(),
	)
	// repos
	postRepo := postrepo.New(client)
	customerRepo := customerrepo.New(client)
	commentRepo := commentrepo.New(client)
	collectionRepo := collectionrepo.New(client)

	// services
	customerSvc := customersvc.New(customerRepo)
	postSvc := postsvc.New(postRepo, businessGateway, storageClient, txManager)
	commentSvc := commentsvc.New(commentRepo, postSvc)
	collectionSvc := collectionsvc.New(collectionRepo, txManager, businessGateway)

	authMiddleware := middleware.Auth(tokenManager)
	optionalAuthMiddleware := middleware.OptionalAuth(tokenManager)
	customerMiddleware := middleware.CustomerID(customerSvc)

	// handlers
	customer.RegisterHandlers(router.Group("/customers"), customerSvc, authMiddleware, storageClient)
	post.RegisterHandlers(router.Group("/posts", optionalAuthMiddleware), postSvc, customerSvc, authMiddleware, storageClient)
	comment.RegisterHandlers(router.Group("/posts", optionalAuthMiddleware), commentSvc, customerSvc, authMiddleware)
	collection.RegisterHandlers(
		router.Group("/collections"),
		collectionSvc,
		authMiddleware,
		optionalAuthMiddleware,
		customerMiddleware,
	)

	go func() {
		logger.Info(ctx, "guest http server is running")
		if err := router.Run(config.Config().GuestHttpServer.Address()); err != nil && err != http.ErrServerClosed {
			logger.Fatal(ctx, "run http server: ", err)
		}
	}()

	closer.Wait()
}
