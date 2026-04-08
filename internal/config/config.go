package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/ua-academy-projects/share-bite/internal/config/env"
)

const (
	adminPrefix    = "ADMIN_"
	guestPrefix    = "GUEST_"
	businessPrefix = "BUSINESS_"
)

type config struct {
	App App

	GuestHttpServer    HttpServer
	AdminHttpServer    HttpServer
	BusinessHttpServer HttpServer

	BusinessHttpClient HttpClient

	Postgres Postgres

	JwtToken  JwtToken
	Email     Email
	RateLimit RateLimit
}

var cfg *config

func Config() *config {
	return cfg
}

type App interface {
	Name() string
	Stage() string
	IsProd() bool
	GracefulShutdownTimeout() time.Duration
}

type HttpServer interface {
	Address() string
	AllowedOrigins() []string
	AllowedMethods() []string
	AllowedHeaders() []string
	ExposeHeaders() []string
}

type HttpClient interface {
	BaseURL() string
	Timeout() time.Duration
	MaxIdleConns() int
	MaxIdleConnsPerHost() int
	IdleConnTimeout() time.Duration
}

type Postgres interface {
	Dsn() string
	MigrationsDir() string
}

type JwtToken interface {
	AccessTokenSecretKey() string
	AccessTokenTTL() time.Duration

	RefreshTokenSecretKey() string
	RefreshTokenTTL() time.Duration
}

type Email interface {
	SenderProviderValue() string
	ResendAPIKeyValue() string
	ResendFromEmailValue() string
}

type RateLimit interface {
	AuthRecoverRequests() int
	AuthRecoverDuration() time.Duration
}

func Load(paths ...string) error {
	if len(paths) > 0 {
		if err := godotenv.Load(paths...); err != nil {
			return fmt.Errorf("load config: %w", err)
		}
	}

	appConfig, err := env.NewAppConfig()
	if err != nil {
		return fmt.Errorf("app config: %w", err)
	}

	guestHttpServerConfig, err := env.NewHttpServerConfig(guestPrefix)
	if err != nil {
		return fmt.Errorf("guest http server config: %w", err)
	}

	adminHttpServerConfig, err := env.NewHttpServerConfig(adminPrefix)
	if err != nil {
		return fmt.Errorf("admin http server config: %w", err)
	}

	businessHttpServerConfig, err := env.NewHttpServerConfig(businessPrefix)
	if err != nil {
		return fmt.Errorf("business http server config: %w", err)
	}

	businessHttpClientConfig, err := env.NewHttpClientConfig(businessPrefix)
	if err != nil {
		return fmt.Errorf("business http client config: %w", err)
	}

	postgresConfig, err := env.NewPostgresConfig()
	if err != nil {
		return fmt.Errorf("postgres config: %w", err)
	}

	jwtTokenConfig, err := env.NewJwtTokenConfig()
	if err != nil {
		return errwrap.Wrap("jwt token config", err)
	}

	emailConfig, err := env.NewEmailConfig()
	if err != nil {
		return errwrap.Wrap("email config", err)
	}

	rateLimitConfig, err := env.NewRateLimitConfig()
	if err != nil {
		return errwrap.Wrap("rate limit config", err)
	}

	cfg = &config{
		App: appConfig,

		GuestHttpServer:    guestHttpServerConfig,
		AdminHttpServer:    adminHttpServerConfig,
		BusinessHttpServer: businessHttpServerConfig,

		BusinessHttpClient: businessHttpClientConfig,

		Postgres:  postgresConfig,
		JwtToken:  jwtTokenConfig,
		Email:     emailConfig,
		RateLimit: rateLimitConfig,
	}

	return nil
}
