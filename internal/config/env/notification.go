package env

type sqsConfig struct {
	QueueURLVal string `env:"SQS_QUEUE_URL"`
	RegionVal   string `env:"AWS_REGION"`
	EndpointVal string `env:"SQS_ENDPOINT_URL"`
}

func NewSQSConfig(prefix string, opts ...Options) (*sqsConfig, error) {
	cfg := new(sqsConfig)
	if err := Parse(cfg, append([]Options{{Prefix: prefix}}, opts...)...); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *sqsConfig) Queue() string {
	return c.QueueURLVal
}

func (c *sqsConfig) Region() string {
	return c.RegionVal
}

func (c *sqsConfig) Endpoint() string {
	return c.EndpointVal
}
