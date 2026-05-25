package worker

import (
	"context"
	"fmt"

	"github.com/ua-academy-projects/share-bite/pkg/email"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	"github.com/ua-academy-projects/share-bite/pkg/notification"
)

type EmailEventHandler interface {
	Handle(ctx context.Context, event notification.Message, emailSender email.Sender) error
}

type EmailProcessor struct {
	emailSender email.Sender
	handlers    map[notification.EventType]EmailEventHandler
}

func NewEmailProcessor(emailSender email.Sender) *EmailProcessor {
	processor := &EmailProcessor{
		emailSender: emailSender,
		handlers:    make(map[notification.EventType]EmailEventHandler),
	}

	processor.RegisterHandler(notification.RegistrationConfirmed, NewRegistrationConfirmedHandler())

	return processor
}

func (p *EmailProcessor) RegisterHandler(eventType notification.EventType, handler EmailEventHandler) {
	p.handlers[eventType] = handler
}

func (p *EmailProcessor) Process(ctx context.Context, event notification.Message) error {
	handler, exists := p.handlers[event.EventType]
	if !exists {
		return fmt.Errorf("unsupported event type for email: %s", event.EventType)
	}

	if err := handler.Handle(ctx, event, p.emailSender); err != nil {
		logger.ErrorKV(ctx, "email handler failed",
			"event_type", event.EventType,
			"event_id", event.EventID,
			"error", err)
		return err
	}

	return nil
}
