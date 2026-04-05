package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
	_ "github.com/ua-academy-projects/share-bite/docs/api/guest"
	"github.com/ua-academy-projects/share-bite/internal/config"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/guest/error/code"
	"github.com/ua-academy-projects/share-bite/internal/guest/gateway/business"
	"github.com/ua-academy-projects/share-bite/internal/guest/handler/collection"
	commenthandler "github.com/ua-academy-projects/share-bite/internal/guest/handler/comment"
	"github.com/ua-academy-projects/share-bite/internal/guest/handler/customer"
	"github.com/ua-academy-projects/share-bite/internal/guest/handler/post"
	collectionrepo "github.com/ua-academy-projects/share-bite/internal/guest/repository/collection"
	commentrepo "github.com/ua-academy-projects/share-bite/internal/guest/repository/comment"
	customerrepo "github.com/ua-academy-projects/share-bite/internal/guest/repository/customer"
	postrepo "github.com/ua-academy-projects/share-bite/internal/guest/repository/post"
	collectionsvc "github.com/ua-academy-projects/share-bite/internal/guest/service/collection"
	commentsvc "github.com/ua-academy-projects/share-bite/internal/guest/service/comment"
	customersvc "github.com/ua-academy-projects/share-bite/internal/guest/service/customer"
	postsvc "github.com/ua-academy-projects/share-bite/internal/guest/service/post"
	"github.com/ua-academy-projects/share-bite/internal/guest/util/response"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
	"github.com/ua-academy-projects/share-bite/pkg/closer"
	"github.com/ua-academy-projects/share-bite/pkg/database/pg"
	"github.com/ua-academy-projects/share-bite/pkg/database/txmanager"
	"github.com/ua-academy-projects/share-bite/pkg/jwt"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	common_middleware "github.com/ua-academy-projects/share-bite/pkg/middleware"
	"github.com/ua-academy-projects/share-bite/pkg/validator"
	"go.uber.org/zap"
	"net/http"
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
	router.Use(common_middleware.RequestLogger())
	router.Use(gin.Recovery())
	router.Use(ErrorMiddleware())

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

	businessGateway := business.NewBusinessAPIClient(config.Config().BusinessHttpClient.BaseURL(), httpClient)

	storageClient, err := newStorageClient(ctx, config.Config().Storage)
	if err != nil {
		logger.Fatal(ctx, "init storage client:", err)
	}

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
	txManager := txmanager.NewTransactionManager(client.DB())
	postSvc := postsvc.New(postRepo, businessGateway, storageClient, txManager)
	commentSvc := commentsvc.New(commentRepo, postSvc)
	collectionSvc := collectionsvc.New(collectionRepo, businessGateway)

	// middlewares
	authMiddleware := middleware.Auth(tokenManager)
	optionalAuthMiddleware := middleware.OptionalAuth(tokenManager)
	customerMiddleware := middleware.CustomerID(customerSvc)

	// handlers
	customer.RegisterHandlers(router.Group("/customers"), customerSvc, authMiddleware)
	post.RegisterHandlers(router.Group("/posts", optionalAuthMiddleware), postSvc, customerSvc, authMiddleware, storageClient)
	commenthandler.RegisterHandlers(router.Group("/posts", optionalAuthMiddleware), commentSvc, customerSvc, authMiddleware)

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

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()
		if err == nil {
			return
		}

		respCode := http.StatusInternalServerError
		resp := response.ErrorResponse{
			Message: "internal server error",
		}

		var validationErr *validator.ValidationError
		if errors.As(err, &validationErr) {
			details := make([]response.ErrorDetail, 0, len(validationErr.Errors))
			for _, e := range validationErr.Errors {
				details = append(details, response.ErrorDetail{
					Field:   e.Field,
					Message: e.Message,
				})
			}
			resp = response.ErrorResponse{
				Message: validationErr.Error(),
				Details: details,
			}

			c.JSON(respCode, resp)
			return
		}

		var appErr *apperror.Error
		if errors.As(err, &appErr) {
			switch appErr.Code {
			case code.NotFound:
				respCode = http.StatusNotFound

			case code.InvalidJSON,
				code.InvalidRequest,
				code.BadRequest,
				code.EmptyUpdate:
				respCode = http.StatusBadRequest

			case code.UpstreamError:
				respCode = http.StatusBadGateway

			case code.AlreadyExists:
				respCode = http.StatusConflict

			case code.Forbidden:
				respCode = http.StatusForbidden

			default:
				respCode = http.StatusInternalServerError
			}

			resp.Message = appErr.Error()
		}

		c.JSON(respCode, resp)
	}
}
