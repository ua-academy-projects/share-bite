package config

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/ua-academy-projects/share-bite/internal/config/env"
	"github.com/ua-academy-projects/share-bite/pkg/errwrap"
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
	Google   Google
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

type Google interface {
	ClientID() string
	ClientSecret() string
	RedirectURL() string
}

func Load(paths ...string) error {
	if len(paths) > 0 {
		if err := godotenv.Load(paths...); err != nil {
			return errwrap.Wrap("load config", err)
		}
	}

	appConfig, err := env.NewAppConfig()
	if err != nil {
		return errwrap.Wrap("app config", err)
	}

	guestHttpServerConfig, err := env.NewHttpServerConfig(guestPrefix)
	if err != nil {
		return errwrap.Wrap("guest http server config", err)
	}

	adminHttpServerConfig, err := env.NewHttpServerConfig(adminPrefix)
	if err != nil {
		return errwrap.Wrap("admin http server config", err)
	}

	businessHttpServerConfig, err := env.NewHttpServerConfig(businessPrefix)
	if err != nil {
		return errwrap.Wrap("business http server config", err)
	}

	postgresConfig, err := env.NewPostgresConfig()
	if err != nil {
		return errwrap.Wrap("postgres config", err)
	}

	jwtTokenConfig, err := env.NewJwtTokenConfig()
	if err != nil {
		return errwrap.Wrap("jwt token config", err)
	}

	googleConfig, err := env.NewGoogleConfig()
	if err != nil {
		return errwrap.Wrap("google config", err)
	}

	cfg = &config{
		App: appConfig,

		GuestHttpServer:    guestHttpServerConfig,
		AdminHttpServer:    adminHttpServerConfig,
		BusinessHttpServer: businessHttpServerConfig,

		Postgres: postgresConfig,
		JwtToken: jwtTokenConfig,
		Google:   googleConfig,
	}

	return nil
}
