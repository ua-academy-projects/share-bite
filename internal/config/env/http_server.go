package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type httpServerConfig struct {
	Host string `env:"HTTP_SERVER_HOST,required"`
	Port string `env:"HTTP_SERVER_PORT,required"`
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
