package env

import (
	"github.com/caarlos0/env/v11"
)

type sqsConfig struct {
	QueueURL string `env:"SQS_QUEUE_URL"`
	Region   string `env:"AWS_REGION"`
	Endpoint string `env:"SQS_ENDPOINT_URL"`
}

func NewSQSConfig(prefix string) (*sqsConfig, error) {
	cfg := new(sqsConfig)
	if err := env.ParseWithOptions(cfg, env.Options{
		Prefix: prefix,
	}); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *sqsConfig) Queue() string {
	return c.QueueURL
}

func (c *sqsConfig) Region() string {
	return c.Region
}

func (c *sqsConfig) Endpoint() string {
	return c.Endpoint
}
