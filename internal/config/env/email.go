package env

import (
	"time"
)

type emailConfig struct {
	EmailSenderProvider string        `env:"EMAIL_SENDER_PROVIDER"`
	ResendAPIKey        string        `env:"RESEND_API_KEY"`
	ResendFromEmail     string        `env:"RESEND_FROM_EMAIL"`
	PasswordResetTTL    time.Duration `env:"PASSWORD_RESET_TTL"`
}

func NewEmailConfig(opts ...Options) (*emailConfig, error) {
	config := new(emailConfig)
	if err := Parse(config, opts...); err != nil {
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

func (c *emailConfig) SenderProviderValue() string {
	if c.EmailSenderProvider == "" {
		return "resend"
	}

	return c.EmailSenderProvider
}

func (c *emailConfig) PasswordResetTTLValue() time.Duration {
	return c.PasswordResetTTL
}
