package main

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/ua-academy-projects/share-bite/internal/config"
	"github.com/ua-academy-projects/share-bite/internal/guest/client/business"
	businessclient "github.com/ua-academy-projects/share-bite/internal/guest/client/business/api/client"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/guest/error/code"
	commenthandler "github.com/ua-academy-projects/share-bite/internal/guest/handler/comment"
	"github.com/ua-academy-projects/share-bite/internal/guest/handler/customer"
	"github.com/ua-academy-projects/share-bite/internal/guest/handler/post"
	commentrepo "github.com/ua-academy-projects/share-bite/internal/guest/repository/comment"
	customerrepo "github.com/ua-academy-projects/share-bite/internal/guest/repository/customer"
	postrepo "github.com/ua-academy-projects/share-bite/internal/guest/repository/post"
	commentsvc "github.com/ua-academy-projects/share-bite/internal/guest/service/comment"
	customersvc "github.com/ua-academy-projects/share-bite/internal/guest/service/customer"
	postsvc "github.com/ua-academy-projects/share-bite/internal/guest/service/post"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
	"github.com/ua-academy-projects/share-bite/pkg/closer"
	"github.com/ua-academy-projects/share-bite/pkg/database/pg"
	"github.com/ua-academy-projects/share-bite/pkg/database/txmanager"
	"github.com/ua-academy-projects/share-bite/pkg/jwt"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	common_middleware "github.com/ua-academy-projects/share-bite/pkg/middleware"
	"github.com/ua-academy-projects/share-bite/pkg/validator"
	"go.uber.org/zap"

	_ "github.com/ua-academy-projects/share-bite/docs/api"
)

// @title			ShareBite Guest API
// @version		1.0
// @description	API for guest customer profile, posts and comments.
//
// @securityDefinitions.apikey	BearerAuth
// @in			header
// @name		Authorization
//
// @BasePath		/

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
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("/swagger/doc.json"),
	))

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

	businessAPIClient := businessclient.New(
		httptransport.NewWithClient(
			businessAPIHost(clientCfg.BaseURL()),
			"/",
			[]string{"http"},
			httpClient,
		),
		strfmt.Default,
	)
	businessGateway := business.NewBusinessAPIClient(businessAPIClient)

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
	authMiddleware := middleware.Auth(tokenManager)
	optionalAuthMiddleware := middleware.OptionalAuth(tokenManager)

	// repos
	postRepo := postrepo.New(client)
	customerRepo := customerrepo.New(client)
	commentRepo := commentrepo.New(client)

	// services
	customerSvc := customersvc.New(customerRepo)
	txManager := txmanager.NewTransactionManager(client.DB())
	postSvc := postsvc.New(postRepo, businessGateway, storageClient, txManager)
	commentSvc := commentsvc.New(commentRepo, postSvc)
	// handlers
	customer.RegisterHandlers(router.Group("/customers"), customerSvc, authMiddleware)
	post.RegisterHandlers(router.Group("/posts", optionalAuthMiddleware), postSvc, customerSvc, authMiddleware, storageClient)
	commenthandler.RegisterHandlers(router.Group("/posts", optionalAuthMiddleware), commentSvc, customerSvc, authMiddleware)

	go func() {
		logger.Info(ctx, "guest http server is running")
		if err := router.Run(config.Config().GuestHttpServer.Address()); err != nil && err != http.ErrServerClosed {
			logger.Fatal(ctx, "run http server: ", err)
		}
	}()

	closer.Wait()
}

func businessAPIHost(baseURL string) string {
	host, port, err := net.SplitHostPort(baseURL)
	if err != nil {
		return baseURL
	}

	if host == "" || host == "0.0.0.0" || host == "::" {
		host = "localhost"
	}

	return net.JoinHostPort(host, port)
}

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()
		if err == nil {
			return
		}

		respCode := http.StatusInternalServerError
		resp := map[string]any{
			"message": "internal server error",
		}

		var valErr *validator.ValidationError
		if errors.As(err, &valErr) {
			respCode = http.StatusBadRequest
			resp = map[string]any{
				"message": valErr.Error(),
				"errors":  valErr.Errors,
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

			resp = map[string]any{
				"message": appErr.Error(),
			}
		}

		c.JSON(respCode, resp)
	}
}
