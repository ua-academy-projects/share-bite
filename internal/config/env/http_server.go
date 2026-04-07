package env

import (
	"net"
	"strings"

	"github.com/caarlos0/env/v11"
)

type httpServerConfig struct {
	Host string `env:"HTTP_SERVER_HOST,required"`
	Port string `env:"HTTP_SERVER_PORT,required"`

	AllowedOriginsRaw string `env:"HTTP_SERVER_ALLOWED_ORIGINS"`
	AllowedMethodsRaw string `env:"HTTP_SERVER_ALLOWED_METHODS"`
	AllowedHeadersRaw string `env:"HTTP_SERVER_ALLOWED_HEADERS"`
	ExposeHeadersRaw  string `env:"HTTP_SERVER_EXPOSE_HEADERS"`
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
	if c.AllowedOriginsRaw == "" {
		return []string{}
	}
	return strings.Split(c.AllowedOriginsRaw, ",")
}

func (c *httpServerConfig) AllowedMethods() []string {
	if c.AllowedMethodsRaw == "" {
		return []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	}
	return strings.Split(c.AllowedMethodsRaw, ",")
}

func (c *httpServerConfig) AllowedHeaders() []string {
	if c.AllowedHeadersRaw == "" {
		return []string{"Origin", "Content-Type", "Accept", "Authorization"}
	}
	return strings.Split(c.AllowedHeadersRaw, ",")
}

func (c *httpServerConfig) ExposeHeaders() []string {
	if c.ExposeHeadersRaw == "" {
		return []string{"Content-Length"}
	}
	return strings.Split(c.ExposeHeadersRaw, ",")
}
