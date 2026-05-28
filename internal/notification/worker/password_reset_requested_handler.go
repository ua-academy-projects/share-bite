package worker

import (
	"context"
	"fmt"

	"github.com/ua-academy-projects/share-bite/pkg/email"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	"github.com/ua-academy-projects/share-bite/pkg/notification"
)

type PasswordResetRequestedHandler struct{}

func NewPasswordResetRequestedHandler() *PasswordResetRequestedHandler {
	return &PasswordResetRequestedHandler{}
}

func (h *PasswordResetRequestedHandler) Handle(ctx context.Context, event notification.Message, emailSender email.Sender) error {
	emailAddr, ok := event.Metadata["email"].(string)
	if !ok || emailAddr == "" {
		return fmt.Errorf("invalid or missing email in metadata")
	}

	resetToken, ok := event.Metadata["reset_token"].(string)
	if !ok || resetToken == "" {
		return fmt.Errorf("invalid or missing reset_token in metadata")
	}

	if err := emailSender.SendPasswordResetToken(ctx, emailAddr, resetToken); err != nil {
		return fmt.Errorf("send password reset email: %w", err)
	}

	logger.InfoKV(ctx, "password reset email sent", "event_id", event.EventID)
	return nil
}
