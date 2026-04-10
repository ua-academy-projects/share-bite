package email

import (
	"context"
	"errors"

	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

type fakeSender struct{}

func NewFakeSender() Sender {
	return &fakeSender{}
}

func (s *fakeSender) SendPasswordResetToken(ctx context.Context, toEmail, token string) error {
	if toEmail == "" {
		return errors.New("recipient email is empty")
	}
	if token == "" {
		return errors.New("reset token is empty")
	}

	logger.InfoKV(
		ctx,
		"password reset email sent by fake sender",
		"to_email", toEmail,
		"reset_token", token,
	)

	return nil
}
