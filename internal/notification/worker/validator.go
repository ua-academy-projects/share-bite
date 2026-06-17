package worker

import (
	"fmt"
	"strings"

	"github.com/ua-academy-projects/share-bite/pkg/notification"
)

type emailEventValidator func(event notification.Message) error

var defaultEmailEventValidators = map[notification.EventType]emailEventValidator{
	notification.RegistrationConfirmed:  validateRegistrationConfirmedEvent,
	notification.PasswordResetRequested: validatePasswordResetRequestedEvent,
}

type DefaultValidator struct {
	allowedEventTypes map[notification.EventType]bool
}

func NewDefaultValidator(allowedEventTypes ...notification.EventType) *DefaultValidator {
	allowed := make(map[notification.EventType]bool)
	if len(allowedEventTypes) == 0 {
		for et := range defaultEmailEventValidators {
			allowed[et] = true
		}
	} else {
		for _, et := range allowedEventTypes {
			allowed[et] = true
		}
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

	validator, ok := defaultEmailEventValidators[event.EventType]
	if !ok {
		return fmt.Errorf("unsupported email event type: %q", event.EventType)
	}

	if err := validator(event); err != nil {
		return err
	}

	return nil
}

func validateRegistrationConfirmedEvent(event notification.Message) error {
	if err := validateCommonEmailEvent(event); err != nil {
		return err
	}

	username, ok := event.Metadata["username"].(string)
	if !ok || strings.TrimSpace(username) == "" {
		return fmt.Errorf("missing metadata.username")
	}

	return nil
}

func validatePasswordResetRequestedEvent(event notification.Message) error {
	if err := validateCommonEmailEvent(event); err != nil {
		return err
	}

	resetToken, ok := event.Metadata["reset_token"].(string)
	if !ok || strings.TrimSpace(resetToken) == "" {
		return fmt.Errorf("missing metadata.reset_token")
	}

	return nil
}

func validateCommonEmailEvent(event notification.Message) error {

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
