package main

import (
	"context"
	"errors"
	"github.com/ua-academy-projects/share-bite/internal/storage/s3"
	"github.com/ua-academy-projects/share-bite/pkg/database/txmanager"
	"net/http"
	"os"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	awscred "github.com/aws/aws-sdk-go-v2/credentials"
	s3sdk "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/ua-academy-projects/share-bite/internal/config"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/guest/error/code"
	"github.com/ua-academy-projects/share-bite/internal/guest/gateway/business"
	"github.com/ua-academy-projects/share-bite/internal/guest/handler/customer"
	"github.com/ua-academy-projects/share-bite/internal/guest/handler/post"
	customerrepo "github.com/ua-academy-projects/share-bite/internal/guest/repository/customer"
	postrepo "github.com/ua-academy-projects/share-bite/internal/guest/repository/post"
	customersvc "github.com/ua-academy-projects/share-bite/internal/guest/service/customer"
	postsvc "github.com/ua-academy-projects/share-bite/internal/guest/service/post"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
	"github.com/ua-academy-projects/share-bite/pkg/closer"
	"github.com/ua-academy-projects/share-bite/pkg/database/pg"
	"github.com/ua-academy-projects/share-bite/pkg/jwt"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	common_middleware "github.com/ua-academy-projects/share-bite/pkg/middleware"
	"github.com/ua-academy-projects/share-bite/pkg/validator"
	"go.uber.org/zap"
)

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

	var storageClient *s3.S3Storage
	{
		s3Endpoint := os.Getenv("S3_ENDPOINT")
		s3Region := os.Getenv("S3_REGION")
		s3AccessKey := os.Getenv("S3_ACCESS_KEY")
		s3SecretKey := os.Getenv("S3_SECRET_KEY")
		s3Bucket := os.Getenv("S3_BUCKET")
		usePathStyle := os.Getenv("S3_USE_PATH_STYLE") == "true"

		if s3Bucket != "" {
			loaderOpts := []func(*awscfg.LoadOptions) error{
				awscfg.WithRegion(s3Region),
			}
			if s3AccessKey != "" && s3SecretKey != "" {
				loaderOpts = append(loaderOpts, awscfg.WithCredentialsProvider(
					awscred.NewStaticCredentialsProvider(s3AccessKey, s3SecretKey, ""),
				))
			}
			if s3Endpoint != "" {
				loaderOpts = append(loaderOpts, awscfg.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
					func(service, region string, options ...interface{}) (aws.Endpoint, error) {
						return aws.Endpoint{URL: s3Endpoint, HostnameImmutable: true}, nil
					},
				)))
			}

			awsCfg, err := awscfg.LoadDefaultConfig(ctx, loaderOpts...)
			if err != nil {
				logger.Fatal(ctx, "aws load config: ", err)
			}
			s3Client := s3sdk.NewFromConfig(awsCfg, func(o *s3sdk.Options) {
				o.UsePathStyle = usePathStyle
			})

			storageClient = s3.NewS3Storage(s3Client, s3Bucket, s3Endpoint)
		}
	}

	tokenManager := jwt.NewTokenManager(
		config.Config().JwtToken.AccessTokenSecretKey(),
		config.Config().JwtToken.RefreshTokenSecretKey(),
		config.Config().JwtToken.AccessTokenTTL(),
		config.Config().JwtToken.RefreshTokenTTL(),
	)
	authMiddleware := middleware.Auth(tokenManager)

	// repos
	postRepo := postrepo.New(client)
	customerRepo := customerrepo.New(client)

	// services
	customerSvc := customersvc.New(customerRepo)
	txManager := txmanager.NewTransactionManager(client.DB())
	postSvc := postsvc.New(postRepo, businessGateway, storageClient, txManager)
	// handlers
	customer.RegisterHandlers(router.Group("/customers"), customerSvc, authMiddleware)
	post.RegisterHandlers(router.Group("/posts"), postSvc, customerSvc, authMiddleware, storageClient)

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
