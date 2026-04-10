package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/ua-academy-projects/share-bite/docs/api"
	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/internal/business/error/code"
	"github.com/ua-academy-projects/share-bite/internal/business/handler/business"
	businessrepo "github.com/ua-academy-projects/share-bite/internal/business/repository/business"
	businesssvc "github.com/ua-academy-projects/share-bite/internal/business/service/business"
	"github.com/ua-academy-projects/share-bite/internal/config"
	"github.com/ua-academy-projects/share-bite/pkg/closer"
	"github.com/ua-academy-projects/share-bite/pkg/database/pg"
	"github.com/ua-academy-projects/share-bite/pkg/database/txmanager"
	"github.com/ua-academy-projects/share-bite/pkg/jwt"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	"go.uber.org/zap"
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
		AllowAllOrigins:  true,
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

	businessRepo := businessrepo.New(client)

	businessSvc := businesssvc.New(businessRepo, txManager)

	tokenManager := jwt.NewTokenManager(
		config.Config().JwtToken.AccessTokenSecretKey(),
		config.Config().JwtToken.RefreshTokenSecretKey(),
		config.Config().JwtToken.AccessTokenTTL(),
		config.Config().JwtToken.RefreshTokenTTL(),
	)

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
