package worker

import (
	"context"
	"fmt"

	"github.com/ua-academy-projects/share-bite/pkg/email"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	"github.com/ua-academy-projects/share-bite/pkg/notification"
)

type RegistrationConfirmedHandler struct{}

func NewRegistrationConfirmedHandler() *RegistrationConfirmedHandler {
	return &RegistrationConfirmedHandler{}
}

func (s *RegistrationConfirmedHandler) Handle(ctx context.Context, event notification.Message, emailSender email.Sender) error {
	emailAddr, ok := event.Metadata["email"].(string)
	if !ok || emailAddr == "" {
		return fmt.Errorf("invalid or missing email in metadata")
	}

	username, ok := event.Metadata["username"].(string)
	if !ok || username == "" {
		return fmt.Errorf("invalid or missing username in metadata")
	}

	if err := emailSender.SendEmail(ctx, emailAddr, "Welcome to Share Bite!", "registration_confirmed", map[string]any{
		"email":    emailAddr,
		"username": username,
	}); err != nil {
		return fmt.Errorf("send registration confirmation email: %w", err)
	}

	logger.InfoKV(ctx, "registration confirmation email sent",
		"email", emailAddr,
		"username", username,
		"event_id", event.EventID)
	return nil
}
