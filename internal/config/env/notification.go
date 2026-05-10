package env

import (
	"os"
)

type sqsConfig struct {
	queue    string
	region   string
	endpoint string
}

func NewSQSConfig(prefix string) (*sqsConfig, error) {
	queue := os.Getenv(prefix + "SQS_QUEUE_URL")
	if queue == "" {
		queue = os.Getenv(prefix + "SQS_QUEUE")
	}
	region := os.Getenv(prefix + "AWS_REGION")
	endpoint := os.Getenv(prefix + "SQS_ENDPOINT_URL")

	return &sqsConfig{
		queue:    queue,
		region:   region,
		endpoint: endpoint,
	}, nil
}

func (c *sqsConfig) Queue() string {
	return c.queue
}

func (c *sqsConfig) Region() string {
	return c.region
}

func (c *sqsConfig) Endpoint() string {
	return c.endpoint
}
