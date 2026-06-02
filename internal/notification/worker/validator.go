package worker

import (
	"fmt"
	"strings"

	"github.com/ua-academy-projects/share-bite/pkg/notification"
)

type DefaultValidator struct {
	allowedEventTypes map[notification.EventType]bool
}

func NewDefaultValidator(allowedEventTypes ...notification.EventType) *DefaultValidator {
	allowed := make(map[notification.EventType]bool)
	for _, et := range allowedEventTypes {
		allowed[et] = true
	}
	return &DefaultValidator{
		allowedEventTypes: allowed,
	}
}

func (v *DefaultValidator) Validate(event notification.Message) error {
	if event.EventID == "" {
		return fmt.Errorf("missing eventID")
	}

	if event.EventType == "" {
		return fmt.Errorf("missing eventType")
	}

	if len(v.allowedEventTypes) > 0 && !v.allowedEventTypes[event.EventType] {
		return fmt.Errorf("unknown event type: %q", event.EventType)
	}

	if event.EventType != notification.RegistrationConfirmed {
		return fmt.Errorf("unsupported email event type: %q", event.EventType)
	}

	if event.RecipientID == "" {
		return fmt.Errorf("missing recipientID")
	}

	if event.CreatedAt.IsZero() {
		return fmt.Errorf("missing createdAt")
	}

	if event.Metadata == nil {
		return fmt.Errorf("missing metadata")
	}

	emailValue, ok := event.Metadata["email"].(string)
	if !ok || strings.TrimSpace(emailValue) == "" {
		return fmt.Errorf("missing metadata.email")
	}

	return nil
}
