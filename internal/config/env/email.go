package env

import "github.com/caarlos0/env/v11"

type emailConfig struct {
	ResendAPIKey    string `env:"RESEND_API_KEY"`
	ResendFromEmail string `env:"RESEND_FROM_EMAIL"`
}

func NewEmailConfig() (*emailConfig, error) {
	config := new(emailConfig)
	if err := env.Parse(config); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *emailConfig) ResendAPIKeyValue() string {
	return c.ResendAPIKey
}

func (c *emailConfig) ResendFromEmailValue() string {
	return c.ResendFromEmail
}
