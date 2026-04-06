package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type httpServerConfig struct {
	Host               string   `env:"HTTP_SERVER_HOST,required"`
	Port               string   `env:"HTTP_SERVER_PORT,required"`
	CorsAllowedOrigins []string `env:"CORS_ALLOWED_ORIGINS" envSeparator:","`
	CorsAllowedMethods []string `env:"CORS_ALLOWED_METHODS" envSeparator:","`
	CorsAllowedHeaders []string `env:"CORS_ALLOWED_HEADERS" envSeparator:","`
}

func NewHttpServerConfig(prefix string) (*httpServerConfig, error) {
	config := new(httpServerConfig)
	if err := env.ParseWithOptions(config, env.Options{
		Prefix: prefix,
	}); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *httpServerConfig) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func (c *httpServerConfig) AllowedOrigins() []string {
	if len(c.CorsAllowedOrigins) == 0 {
		return []string{"*"}
	}
	return c.CorsAllowedOrigins
}

func (c *httpServerConfig) AllowedMethods() []string {
	if len(c.CorsAllowedMethods) == 0 {
		return []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	}
	return c.CorsAllowedMethods
}

func (c *httpServerConfig) AllowedHeaders() []string {
	if len(c.CorsAllowedHeaders) == 0 {
		return []string{"Origin", "Content-Type", "Accept", "Authorization"}
	}
	return c.CorsAllowedHeaders
}
