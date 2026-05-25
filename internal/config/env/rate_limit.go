package env

import (
	"errors"
	"time"
)

type rateLimitConfig struct {
	AuthRecoverRequestsValue int           `env:"RATE_LIMIT_AUTH_RECOVER_REQUESTS" envDefault:"5"`
	AuthRecoverDurationValue time.Duration `env:"RATE_LIMIT_AUTH_RECOVER_DURATION" envDefault:"10m"`
}

func NewRateLimitConfig(opts ...Options) (*rateLimitConfig, error) {
	config := new(rateLimitConfig)
	if err := Parse(config, opts...); err != nil {
		return nil, err
	}

	if config.AuthRecoverRequestsValue <= 0 || config.AuthRecoverDurationValue <= 0 {
		return nil, errors.New("recover rate-limit values must be greater than zero")
	}

	return config, nil
}

func (c *rateLimitConfig) AuthRecoverRequests() int {
	return c.AuthRecoverRequestsValue
}

func (c *rateLimitConfig) AuthRecoverDuration() time.Duration {
	return c.AuthRecoverDurationValue
}
