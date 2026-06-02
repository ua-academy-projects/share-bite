package env

import (
	"github.com/caarlos0/env/v11"
)

type h3Config struct {
	H3Resolution      int `env:"H3_RESOLUTION" envDefault:"7"`
	H3RecommendRadius int `env:"H3_RECOMMEND_RADIUS" envDefault:"2"`
}

func NewH3Config() (*h3Config, error) {
	config := new(h3Config)
	if err := env.Parse(config); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *h3Config) Resolution() int {
	return c.H3Resolution
}

func (c *h3Config) RecommendRadius() int {
	return c.H3RecommendRadius
}
