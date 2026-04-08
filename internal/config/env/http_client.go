package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type httpClientConfig struct {
	BaseURLVal             string        `env:"HTTP_CLIENT_BASE_URL"`
	TimeoutVal             time.Duration `env:"HTTP_CLIENT_TIMEOUT" envDefault:"10s"`
	MaxIdleConnsVal        int           `env:"HTTP_CLIENT_MAX_IDLE_CONNS" envDefault:"100"`
	MaxIdleConnsPerHostVal int           `env:"HTTP_CLIENT_MAX_IDLE_CONNS_PER_HOST" envDefault:"100"`
	IdleConnTimeoutVal     time.Duration `env:"HTTP_CLIENT_IDLE_CONN_TIMEOUT" envDefault:"90s"`
}

func NewHttpClientConfig(prefix string) (*httpClientConfig, error) {
	config := new(httpClientConfig)
	if err := env.ParseWithOptions(config, env.Options{
		Prefix: prefix,
	}); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *httpClientConfig) BaseURL() string {
	return c.BaseURLVal
}

func (c *httpClientConfig) Timeout() time.Duration {
	return c.TimeoutVal
}

func (c *httpClientConfig) MaxIdleConns() int {
	return c.MaxIdleConnsVal
}

func (c *httpClientConfig) MaxIdleConnsPerHost() int {
	return c.MaxIdleConnsPerHostVal
}

func (c *httpClientConfig) IdleConnTimeout() time.Duration {
	return c.IdleConnTimeoutVal
}
