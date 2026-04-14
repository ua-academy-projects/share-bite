package main

import (
	"context"
	"errors"
	"net/http"
	"os"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	awscred "github.com/aws/aws-sdk-go-v2/credentials"
	s3sdk "github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/ua-academy-projects/share-bite/docs/api/business"
	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/internal/business/error/code"
	"github.com/ua-academy-projects/share-bite/internal/business/handler/business"
	businessrepo "github.com/ua-academy-projects/share-bite/internal/business/repository/business"
	businesssvc "github.com/ua-academy-projects/share-bite/internal/business/service/business"
	"github.com/ua-academy-projects/share-bite/internal/config"
	"github.com/ua-academy-projects/share-bite/internal/storage/s3"
	"github.com/ua-academy-projects/share-bite/pkg/closer"
	"github.com/ua-academy-projects/share-bite/pkg/database/pg"
	"github.com/ua-academy-projects/share-bite/pkg/database/txmanager"

	"github.com/ua-academy-projects/share-bite/pkg/jwt"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	"go.uber.org/zap"
	_ "github.com/ua-academy-projects/share-bite/docs/api/business"
)

// @title			ShareBite Business API
// @version		1.0
// @description	API for discovering brand locations (venues).
//
// @securityDefinitions.apikey	BearerAuth
// @in			header
// @name		Authorization
//
// @BasePath		/
func main() {
	ctx := context.Background()

	if err := config.Load(".env"); err != nil {
		logger.Fatal(ctx, "load config:", err)
	}

	cfg := config.Config()

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

	txManager := txmanager.NewTransactionManager(client.DB())

	var storageClient *s3.S3Storage
	{
		s3Endpoint := os.Getenv("S3_ENDPOINT")
		s3Region := os.Getenv("S3_REGION")
		s3AccessKey := os.Getenv("S3_ACCESS_KEY")
		s3SecretKey := os.Getenv("S3_SECRET_KEY")
		s3Bucket := os.Getenv("S3_BUCKET")
		if s3Bucket == "" {
			logger.Fatal(ctx, "S3_BUCKET is required but not set") 
		}
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

	// repos
	businessRepo := businessrepo.New(client)

	// services
	businessSvc := businesssvc.New(businessRepo, txManager, storageClient)

	tokenManager := jwt.NewTokenManager(
		config.Config().JwtToken.AccessTokenSecretKey(),
		config.Config().JwtToken.RefreshTokenSecretKey(),
		config.Config().JwtToken.AccessTokenTTL(),
		config.Config().JwtToken.RefreshTokenTTL(),
	)


	// handlers
	business.RegisterHandlers(router.Group("/business"), businessSvc, tokenManager)


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
			}
		}

		logger.ErrorKV(ctx, "internal error", "error", err.Err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
