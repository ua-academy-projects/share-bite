// @title Share Bite Admin Auth API
// @version 1.0
// @description Admin authentication API documentation.
// @BasePath /
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	apperror "github.com/ua-academy-projects/share-bite/internal/admin-auth/error"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/error/code"
	"go.uber.org/zap"

	_ "github.com/ua-academy-projects/share-bite/docs/admin_auth"

	"github.com/ua-academy-projects/share-bite/internal/config"
	"github.com/ua-academy-projects/share-bite/pkg/closer"
	"github.com/ua-academy-projects/share-bite/pkg/database/pg"
	"github.com/ua-academy-projects/share-bite/pkg/database/txmanager"
	"github.com/ua-academy-projects/share-bite/pkg/logger"

	"github.com/ua-academy-projects/share-bite/pkg/jwt"

	authhttp "github.com/ua-academy-projects/share-bite/internal/admin-auth/handler/auth"
	adminmw "github.com/ua-academy-projects/share-bite/internal/admin-auth/middleware"
	userrepo "github.com/ua-academy-projects/share-bite/internal/admin-auth/repository/user"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/routers"
	authsvc "github.com/ua-academy-projects/share-bite/internal/admin-auth/service/auth"
	emailsvc "github.com/ua-academy-projects/share-bite/internal/admin-auth/service/email"
)

func main() {
	ctx := context.Background()

	if err := config.Load(".env"); err != nil {
		logger.Fatal(ctx, "load config:", err)
	}

	router := gin.New()
	router.Use(gin.Recovery())
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

	txManager := txmanager.NewTransactionManager(client.DB())
	userRepo := userrepo.New(client, txManager)

	provider := strings.ToLower(strings.TrimSpace(cfg.Email.SenderProviderValue()))
	var emailSender emailsvc.Sender

	switch provider {
	case "", "resend":
		emailSender, err = emailsvc.NewResendSender(
			cfg.Email.ResendAPIKeyValue(),
			cfg.Email.ResendFromEmailValue(),
		)
		if err != nil {
			logger.Fatal(ctx, "new resend email sender: ", err)
		}
	case "fake":
		emailSender = emailsvc.NewFakeSender()
	default:
		logger.Fatal(ctx, "new email sender: ", fmt.Errorf("unknown email sender provider: %s", provider))
	}

	authSvc := authsvc.New(userRepo, tokenManager, emailSender)

	authHandler := authhttp.NewHandler(authSvc)
	limiter := adminmw.NewAuthRecoveryLimiter(
		cfg.RateLimit.AuthRecoverRequests(),
		cfg.RateLimit.AuthRecoverDuration(),
	)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	routers.SetupRouter(router.Group("/"), authHandler, limiter)

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
		resp := map[string]any{
			"message": "internal server error",
		}

		var appErr *apperror.Error
		if errors.As(err, &appErr) {
			switch appErr.Code {
			case code.NotFound:
				respCode = http.StatusNotFound
			case code.InvalidRequest:
				respCode = http.StatusBadRequest

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
