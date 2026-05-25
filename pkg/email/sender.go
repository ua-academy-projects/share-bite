package email

import (
	"context"
)

type Sender interface {
	SendPasswordResetToken(ctx context.Context, toEmail, token string) error
	SendEmail(ctx context.Context, toEmail, subject, templateName string, data map[string]any) error
}
