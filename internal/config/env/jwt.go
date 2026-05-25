package env

import (
	"time"
)

type jwtTokenConfig struct {
	JwtAccessTokenSecretKey  string        `env:"JWT_ACCESS_TOKEN_SECRET_KEY,required"`
	JwtAccessTokenTTL        time.Duration `env:"JWT_ACCESS_TOKEN_TTL,required"`
	JwtRefreshTokenSecretKey string        `env:"JWT_REFRESH_TOKEN_SECRET_KEY,required"`
	JwtRefreshTokenTTL       time.Duration `env:"JWT_REFRESH_TOKEN_TTL,required"`
}

func NewJwtTokenConfig(opts ...Options) (*jwtTokenConfig, error) {
	config := new(jwtTokenConfig)
	if err := Parse(config, opts...); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *jwtTokenConfig) AccessTokenSecretKey() string {
	return c.JwtAccessTokenSecretKey
}

func (c *jwtTokenConfig) RefreshTokenSecretKey() string {
	return c.JwtRefreshTokenSecretKey
}

func (c *jwtTokenConfig) AccessTokenTTL() time.Duration {
	return c.JwtAccessTokenTTL
}

func (c *jwtTokenConfig) RefreshTokenTTL() time.Duration {
	return c.JwtRefreshTokenTTL
}
