package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/business/handler/business"
	businessrepo "github.com/ua-academy-projects/share-bite/internal/business/repository/business"
	businesssvc "github.com/ua-academy-projects/share-bite/internal/business/service/business"
	"github.com/ua-academy-projects/share-bite/internal/config"
	"github.com/ua-academy-projects/share-bite/pkg/closer"
	"github.com/ua-academy-projects/share-bite/pkg/database/pg"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
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
	router.Use(gin.Recovery())
	router.Use(ErrorMiddleware())

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

	// repos
	businessRepo := businessrepo.New(client)

	// services
	businessSvc := businesssvc.New(businessRepo)

	// handlers
	business.RegisterHandlers(router.Group("/business"), businessSvc)

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

		respCode := http.StatusInternalServerError
		resp := map[string]any{
			"message": "internal server error",
		}

		// TODO: handle custom errors

		c.JSON(respCode, resp)
	}
}
