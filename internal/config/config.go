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

	Postgres Postgres

	JwtToken JwtToken
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

	postgresConfig, err := env.NewPostgresConfig()
	if err != nil {
		return fmt.Errorf("postgres config: %w", err)
	}

	jwtTokenConfig, err := env.NewJwtTokenConfig()
	if err != nil {
		return fmt.Errorf("jwt token config: %w", err)
	}

	cfg = &config{
		App: appConfig,

		GuestHttpServer:    guestHttpServerConfig,
		AdminHttpServer:    adminHttpServerConfig,
		BusinessHttpServer: businessHttpServerConfig,

		Postgres: postgresConfig,
		JwtToken: jwtTokenConfig,
	}

	return nil
}
