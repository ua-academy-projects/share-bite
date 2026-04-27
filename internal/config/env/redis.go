package env

import "github.com/caarlos0/env/v11"

type RedisConfig struct {
	RedisHost     string `env:"REDIS_HOST" envDefault:"localhost"`
	RedisPort     string `env:"REDIS_PORT" envDefault:"6379"`
	RedisPassword string `env:"REDIS_PASSWORD" envDefault:""`
	RedisTLS      bool   `env:"REDIS_TLS_ENABLED" envDefault:"false"`
	RedisDB       int    `env:"REDIS_DB" envDefault:"0"`
}

func NewRedisConfig() (*RedisConfig, error) {
	cfg := new(RedisConfig)
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *RedisConfig) Addr() string {
	return c.RedisHost + ":" + c.RedisPort
}

func (c *RedisConfig) Password() string {
	return c.RedisPassword
}

func (c *RedisConfig) TLS() bool {
	return c.RedisTLS
}

func (c *RedisConfig) DB() int {
	return c.RedisDB
}
