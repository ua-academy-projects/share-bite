package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	businessclient "github.com/ua-academy-projects/share-bite/internal/admin-auth/adapter/business"
	guestclient "github.com/ua-academy-projects/share-bite/internal/admin-auth/adapter/guest"
	apperr "github.com/ua-academy-projects/share-bite/internal/admin-auth/error"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/handler"
	adminhttp "github.com/ua-academy-projects/share-bite/internal/admin-auth/handler/admin"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/worker"
	"github.com/ua-academy-projects/share-bite/internal/config/env"
	"go.uber.org/zap"

	authhttp "github.com/ua-academy-projects/share-bite/internal/admin-auth/handler/auth"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/provider"
	gh "github.com/ua-academy-projects/share-bite/internal/admin-auth/provider/github"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/provider/google"
	userrepo "github.com/ua-academy-projects/share-bite/internal/admin-auth/repository/user"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/routers"
	adminsvc "github.com/ua-academy-projects/share-bite/internal/admin-auth/service/admin"
	authsvc "github.com/ua-academy-projects/share-bite/internal/admin-auth/service/auth"
	"github.com/ua-academy-projects/share-bite/internal/config"

	"github.com/ua-academy-projects/share-bite/internal/middleware"
	pkgmw "github.com/ua-academy-projects/share-bite/pkg/middleware"

	"github.com/ua-academy-projects/share-bite/pkg/closer"
	"github.com/ua-academy-projects/share-bite/pkg/database/pg"
	"github.com/ua-academy-projects/share-bite/pkg/database/txmanager"
	"github.com/ua-academy-projects/share-bite/pkg/jwt"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	"github.com/ua-academy-projects/share-bite/pkg/outbox"

	adminmw "github.com/ua-academy-projects/share-bite/internal/admin-auth/middleware"
)

// @title			Share Bite Admin Auth API
// @version		1.0
// @description	Admin authentication API documentation.
// @BasePath		/
func main() {
	ctx := context.Background()

	if err := config.Load(".env"); err != nil {
		logger.Fatal(ctx, err)
	}

	googleCfg, err := env.NewGoogleConfig()
	if err != nil {
		logger.Fatal(ctx, "load google oauth config: ", err)
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
	txManager := txmanager.NewTransactionManager(client.DB())
	userRepo := userrepo.New(client)
	adminRepo := userrepo.NewAdmin(client)

	workerManager := worker.NewManager(userRepo)
	workerManager.Start(ctx)
	closer.Add(func(ctx context.Context) error {
		workerManager.Stop()
		return nil
	})

	providerFactory := provider.NewFactory(google.Config{
		ClientID:     googleCfg.ClientID(),
		ClientSecret: googleCfg.ClientSecret(),
		RedirectURL:  googleCfg.RedirectURL(),
	})
	outboxWriter := outbox.NewWriter(client.DB())
	authSvc := authsvc.New(userRepo, tokenManager, txManager, cfg.Email.PasswordResetTTLValue(), outboxWriter, cfg.Auth.MaxSessions())
	authHandler := authhttp.NewHandler(authSvc, providerFactory)

	customerClient := guestclient.NewClient(client)
	businessClient := businessclient.NewClient(client)

	adminSvc := adminsvc.NewService(adminRepo, userRepo, customerClient, businessClient, txManager)
	adminHandler := adminhttp.NewHandler(adminSvc)

	limiter := adminmw.NewAuthRecoveryLimiter(
		cfg.RateLimit.AuthRecoverRequests(),
		cfg.RateLimit.AuthRecoverDuration(),
	)

	ghConfig := gh.Config{
		ClientID:           cfg.Github.GetClientID(),
		ClientSecret:       cfg.Github.GetClientSecret(),
		RedirectURL:        cfg.Github.GetRedirectURL(),
		SuccessRedirectURL: cfg.Github.GetSuccessRedirectURL(),
	}

	sessionStore := gh.NewJWTSessionStore(tokenManager)
	ghHandler := gh.NewHandler(ghConfig, userRepo, sessionStore, txManager)

	routers.SetupRouter(router.Group("/"), authHandler, adminHandler, *ghHandler, authMw, limiter)

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

		respCode := http.StatusInternalServerError
		resp := handler.ErrorResponse{Error: "internal server error"}

		var appErr *apperr.AppError
		if errors.As(err.Err, &appErr) {
			respCode = appErr.HTTPStatus()

			resp = handler.ErrorResponse{Error: appErr.Message}
		}

		c.JSON(respCode, resp)
	}
}
