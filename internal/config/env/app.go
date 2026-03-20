package env

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

const (
	prodStage = "prod"
	devStage  = "dev"
)

type appConfig struct {
	AppName                    string        `env:"APP_NAME,required"`
	AppStage                   string        `env:"APP_STAGE,required"`
	AppGracefulShutdownTimeout time.Duration `env:"APP_GRACEFUL_SHUTDOWN_TIMEOUT,required"`
}

func NewAppConfig() (*appConfig, error) {
	config := new(appConfig)
	if err := env.Parse(config); err != nil {
		return nil, err
	}

	if config.AppStage != prodStage && config.AppStage != devStage {
		return nil, fmt.Errorf(`unknown stage option: %s (only %s or %s)`, config.AppStage, prodStage, devStage)
	}

	return config, nil
}

func (c *appConfig) Name() string {
	return c.AppName
}

func (c *appConfig) Stage() string {
	return c.AppStage
}

func (c *appConfig) IsProd() bool {
	return c.AppStage == prodStage
}

func (c *appConfig) GracefulShutdownTimeout() time.Duration {
	return c.AppGracefulShutdownTimeout
}
