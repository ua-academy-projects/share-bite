package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	apperr "github.com/ua-academy-projects/share-bite/internal/admin-auth/error"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/provider/email"
	"github.com/ua-academy-projects/share-bite/internal/config/env"
	"go.uber.org/zap"
	"github.com/ua-academy-projects/share-bite/internal/config"
	"github.com/ua-academy-projects/share-bite/pkg/closer"
	"github.com/ua-academy-projects/share-bite/pkg/database/pg"
	"github.com/ua-academy-projects/share-bite/pkg/jwt"
	"github.com/ua-academy-projects/share-bite/pkg/logger"

	gh "github.com/ua-academy-projects/share-bite/internal/admin-auth/ghAuth"
	authhttp "github.com/ua-academy-projects/share-bite/internal/admin-auth/handler/auth"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/provider"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/provider/google"
	userrepo "github.com/ua-academy-projects/share-bite/internal/admin-auth/repository/user"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/routers"
	authsvc "github.com/ua-academy-projects/share-bite/internal/admin-auth/service/auth"

	"github.com/ua-academy-projects/share-bite/internal/middleware"
	pkgmw "github.com/ua-academy-projects/share-bite/pkg/middleware"

	"github.com/ua-academy-projects/share-bite/pkg/database/txmanager"
	commonmiddleware "github.com/ua-academy-projects/share-bite/pkg/middleware"

	// @title           Share Bite Admin Auth API
	// @version         1.0
	// @description     This is an authentication microservice for Share Bite.

	// @host            localhost:3850
	// @BasePath        /

	adminmw "github.com/ua-academy-projects/share-bite/internal/admin-auth/middleware"
)

// @title Share Bite Admin Auth API
// @version 1.0
// @description Admin authentication API documentation.
// @BasePath /
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
	router.Use(commonmiddleware.RequestID())
	router.Use(commonmiddleware.RequestLogger())
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

	providerFactory := provider.NewFactory(google.Config{
		ClientID:     googleCfg.ClientID(),
		ClientSecret: googleCfg.ClientSecret(),
		RedirectURL:  googleCfg.RedirectURL(),
	})

	providerStr := strings.ToLower(strings.TrimSpace(cfg.Email.SenderProviderValue()))
	var emailSender email.Sender

	switch providerStr {
	case "", "resend":
		emailSender, err = email.NewResendSender(
			cfg.Email.ResendAPIKeyValue(),
			cfg.Email.ResendFromEmailValue(),
		)
		if err != nil {
			logger.Fatal(ctx, "new resend email sender: ", err)
		}
	case "fake":
		emailSender = email.NewFakeSender()
	default:
		logger.Fatal(ctx, "new email sender: ", fmt.Errorf("unknown email sender provider: %s", providerStr))
	}
	authSvc := authsvc.New(userRepo, tokenManager, emailSender, txManager)
	authHandler := authhttp.NewHandler(authSvc, providerFactory)

	limiter := adminmw.NewAuthRecoveryLimiter(
		cfg.RateLimit.AuthRecoverRequests(),
		cfg.RateLimit.AuthRecoverDuration(),
	)

	ghConfig := gh.Config{
		ClientID:           cfg.GitHub.GetClientID(),
		ClientSecret:       cfg.GitHub.GetClientSecret(),
		RedirectURL:        cfg.GitHub.GetRedirectURL(),
		SuccessRedirectURL: cfg.GitHub.GetSuccessRedirectURL(),
	}
	
	sessionStore := gh.NewJWTSessionStore(tokenManager)
	ghHandler := gh.NewHandler(ghConfig, userRepo, sessionStore)

	routers.SetupRouter(router.Group("/"), authHandler,authMw, limiter, ghHandler)

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
		resp := authhttp.ErrorResponse{Error: "internal server error"}

		var appErr *apperr.AppError
		if errors.As(err.Err, &appErr) {
			respCode = appErr.HTTPStatus()

			resp = authhttp.ErrorResponse{Error: appErr.Error()}
		}

		c.JSON(respCode, resp)
	}
}