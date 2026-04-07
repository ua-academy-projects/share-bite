package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/config/env"
	"go.uber.org/zap"

	apperr "github.com/ua-academy-projects/share-bite/internal/admin-auth/error"
	authhttp "github.com/ua-academy-projects/share-bite/internal/admin-auth/handler/auth"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/provider"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/provider/google"
	userrepo "github.com/ua-academy-projects/share-bite/internal/admin-auth/repository/user"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/routers"
	authsvc "github.com/ua-academy-projects/share-bite/internal/admin-auth/service/auth"
	"github.com/ua-academy-projects/share-bite/internal/config"

	"github.com/ua-academy-projects/share-bite/internal/middleware"
	pkgmw "github.com/ua-academy-projects/share-bite/pkg/middleware"

	"github.com/ua-academy-projects/share-bite/pkg/closer"
	"github.com/ua-academy-projects/share-bite/pkg/database/pg"
	"github.com/ua-academy-projects/share-bite/pkg/jwt"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

// @title           Share Bite Admin Auth API
// @version         1.0
// @description     This is an authentication microservice for Share Bite.

// @host            localhost:3850
// @BasePath        /

// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
func main() {
	ctx := context.Background()

	if err := config.Load(".env"); err != nil {
		logger.Fatal(ctx, "load config:", err)
	}

	googleCfg, err := env.NewGoogleConfig()
	if err != nil {
		log.Fatalf("load google oauth config: %v", err)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(pkgmw.RequestID())
	router.Use(pkgmw.RequestLogger())
	router.Use(ErrorMiddleware())

	cfg := config.Config()

	if cfg.App.IsProd() {
		logger.SetLevel(zap.InfoLevel)
		gin.SetMode(gin.ReleaseMode)
	} else {
		logger.SetLevel(zap.DebugLevel)
	}
	closer.SetShutdownTimeout(cfg.App.GracefulShutdownTimeout())

	client, err := pg.NewClient(ctx, cfg.Postgres.Dsn())
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

	tokenManager := jwt.NewTokenManager(
		cfg.JwtToken.AccessTokenSecretKey(),
		cfg.JwtToken.RefreshTokenSecretKey(),
		cfg.JwtToken.AccessTokenTTL(),
		cfg.JwtToken.RefreshTokenTTL(),
	)

	authMw := middleware.Auth(tokenManager)

	userRepo := userrepo.New(client)
	authSvc := authsvc.New(userRepo, tokenManager)

	providerFactory := provider.NewFactory(google.Config{
		ClientID:     googleCfg.ClientID(),
		ClientSecret: googleCfg.ClientSecret(),
		RedirectURL:  googleCfg.RedirectURL(),
	})

	authHandler := authhttp.NewHandler(authSvc, providerFactory)

	routers.SetupRouter(router.Group("/"), authHandler, authMw)

	go func() {
		addr := cfg.AdminHttpServer.Address()
		logger.Info(ctx, "auth http server is running on "+addr)
		if err := router.Run(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
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

		var appErr *apperr.AppError
		if errors.As(err.Err, &appErr) {
			c.JSON(appErr.HTTPStatus(), gin.H{"error": appErr.Message})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
