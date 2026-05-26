package config

import (
	"context"
	"fmt"
	"time"

	"github.com/ua-academy-projects/share-bite/pkg/logger"

	"github.com/joho/godotenv"
	"github.com/ua-academy-projects/share-bite/internal/config/env"
)

const (
	adminPrefix           = "ADMIN_"
	guestPrefix           = "GUEST_"
	businessPrefix        = "BUSINESS_"
	notificationPrefix    = "NOTIFICATION_"
	imageProcessingPrefix = "IMAGE_PROCESSING_"
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
	H3        H3
	Email     Email
	RateLimit RateLimit
	Github    GitHub

	Storage Storage

	Auth AuthConfig

	NotificationHttpServer HttpServer
	NotificationSQS        SQS

	ImageProcessingSQS SQS
}

type SQS interface {
	Queue() string
	Region() string
	Endpoint() string
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

type H3 interface {
	Resolution() int
	RecommendRadius() int
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
	PresignTTL() time.Duration
}

func Load(paths ...string) error {
	return LoadWithSecrets(nil, paths...)
}

// LoadWithSecrets loads configuration from .env files and merges them with provided secrets.
// Caller-provided secrets take precedence over .env entries to allow runtime overrides (e.g. from AWS Secrets Manager).
func LoadWithSecrets(secrets map[string]string, paths ...string) error {
	allSecrets := make(map[string]string)
	if len(paths) > 0 {
		dotEnv, err := godotenv.Read(paths...)
		if err == nil {
			for k, v := range dotEnv {
				allSecrets[k] = v
			}
		} else {
			logger.Info(context.Background(), "No .env file found, relying on system environment variables")
		}
	}

	for k, v := range secrets {
		allSecrets[k] = v
	}

	env.Init(allSecrets)

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
	h3Config, err := env.NewH3Config()
	if err != nil {
		return fmt.Errorf("h3 config: %w", err)
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
		return fmt.Errorf("github config: %w", err)
	}

	notificationHttpServerConfig, err := env.NewHttpServerConfig(notificationPrefix)
	if err != nil {
		return fmt.Errorf("notification http server config: %w", err)
	}

	notificationSQSConfig, err := env.NewSQSConfig(notificationPrefix)
	if err != nil {
		return fmt.Errorf("notification sqs config: %w", err)
	}

	imageProcessingSQSConfig, err := env.NewSQSConfig(imageProcessingPrefix)
	if err != nil {
		return fmt.Errorf("image processing sqs config: %w", err)
	}

	cfg = &config{
		App:    appConfig,
		Auth:   authConfig,
		Github: ghcfg,

		GuestHttpServer:    guestHttpServerConfig,
		AdminHttpServer:    adminHttpServerConfig,
		BusinessHttpServer: businessHttpServerConfig,

		BusinessHttpClient: businessHttpClientConfig,

		Postgres:  postgresConfig,
		Redis:     redisConfig,
		JwtToken:  jwtTokenConfig,
		H3:        h3Config,
		Email:     emailConfig,
		RateLimit: rateLimitConfig,

		Storage: storageConfig,

		NotificationHttpServer: notificationHttpServerConfig,
		NotificationSQS:        notificationSQSConfig,

		ImageProcessingSQS: imageProcessingSQSConfig,
	}

	return nil
}

func GetSecret(key string) string {
	return env.GetSecret(key)
}
