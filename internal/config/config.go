package config

import (
	"context"
	"fmt"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
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
	Redis    Redis

	JwtToken  JwtToken
	Email     Email
	RateLimit RateLimit
	Github    GitHub

	Storage Storage

	Auth AuthConfig
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

type GitHub interface {
	GetClientID() string
	GetClientSecret() string
	GetRedirectURL() string
	GetSuccessRedirectURL() string
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
	Scheme() string
	Timeout() time.Duration
	MaxIdleConns() int
	MaxIdleConnsPerHost() int
	IdleConnTimeout() time.Duration
}

type Postgres interface {
	Dsn() string
	MigrationsDir() string
}

type Redis interface {
	Addr() string
	Password() string
	TLS() bool
	DB() int
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
	PasswordResetTTLValue() time.Duration
}

type AuthConfig interface {
	MaxSessions() int
}

type RateLimit interface {
	AuthRecoverRequests() int
	AuthRecoverDuration() time.Duration
}

type Storage interface {
	Endpoint() string
	Region() string
	AccessKey() string
	SecretKey() string
	Bucket() string
	UsePathStyle() bool
}

func Load(paths ...string) error {
	if len(paths) > 0 {
		if err := godotenv.Load(paths...); err != nil {
			logger.Info(context.Background(), "No .env file found, relying on system environment variables")
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

	redisConfig, err := env.NewRedisConfig()
	if err != nil {
		return fmt.Errorf("redis config: %w", err)
	}

	jwtTokenConfig, err := env.NewJwtTokenConfig()
	if err != nil {
		return fmt.Errorf("jwt token config: %w", err)
	}

	emailConfig, err := env.NewEmailConfig()
	if err != nil {
		return fmt.Errorf("email config: %w", err)
	}

	rateLimitConfig, err := env.NewRateLimitConfig()
	if err != nil {
		return fmt.Errorf("rate limit config: %w", err)
	}

	storageConfig, err := env.NewS3StorageConfig()
	if err != nil {
		return fmt.Errorf("storage config: %w", err)
	}

	authConfig, err := env.NewAuthConfig()
	if err != nil {
		return fmt.Errorf("auth config: %w", err)
	}
	
	ghcfg, err := env.NewGitHubConfig()
	if err != nil {	
		return fmt.Errorf("Errorl load github config: %w", err)
	}

	cfg = &config{
		App:  appConfig,
		Auth: authConfig,
		Github: ghcfg,

		GuestHttpServer:    guestHttpServerConfig,
		AdminHttpServer:    adminHttpServerConfig,
		BusinessHttpServer: businessHttpServerConfig,

		BusinessHttpClient: businessHttpClientConfig,

		Postgres:  postgresConfig,
		Redis:     redisConfig,
		JwtToken:  jwtTokenConfig,
		Email:     emailConfig,
		RateLimit: rateLimitConfig,

		Storage: storageConfig,
	}

	return nil
}
