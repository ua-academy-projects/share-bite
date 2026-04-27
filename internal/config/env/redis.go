package env

import "github.com/caarlos0/env/v11"

type redisConfig struct {
	RedisHost     string `env:"REDIS_HOST" envDefault:"localhost"`
	RedisPort     string `env:"REDIS_PORT" envDefault:"6379"`
	RedisPassword string `env:"REDIS_PASSWORD" envDefault:""`
}

func NewRedisConfig() (*redisConfig, error) {
	cfg := new(redisConfig)
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *redisConfig) Addr() string {
	return c.RedisHost + ":" + c.RedisPort
}

func (c *redisConfig) Password() string {
	return c.RedisPassword
}
