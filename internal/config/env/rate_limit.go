package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type rateLimitConfig struct {
	AuthRecoverRequestsValue int           `env:"RATE_LIMIT_AUTH_RECOVER_REQUESTS" envDefault:"6"`
	AuthRecoverDurationValue time.Duration `env:"RATE_LIMIT_AUTH_RECOVER_DURATION" envDefault:"10m"`
}

func NewRateLimitConfig() (*rateLimitConfig, error) {
	config := new(rateLimitConfig)
	if err := env.Parse(config); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *rateLimitConfig) AuthRecoverRequests() int {
	return c.AuthRecoverRequestsValue
}

func (c *rateLimitConfig) AuthRecoverDuration() time.Duration {
	return c.AuthRecoverDurationValue
}
