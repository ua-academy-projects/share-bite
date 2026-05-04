package worker

import (
	"fmt"

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
		return fmt.Errorf("missing event_id")
	}

	if event.EventType == "" {
		return fmt.Errorf("missing event_type")
	}

	if len(v.allowedEventTypes) > 0 && !v.allowedEventTypes[event.EventType] {
		return fmt.Errorf("unknown event type: %q", event.EventType)
	}

	if event.RecipientID == "" {
		return fmt.Errorf("missing recipient_id")
	}

	if event.CreatedAt.IsZero() {
		return fmt.Errorf("missing created_at")
	}

	return nil
}
