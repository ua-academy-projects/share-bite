package worker

import (
	"context"
	"time"

	"github.com/ua-academy-projects/share-bite/pkg/logger"
	"github.com/ua-academy-projects/share-bite/pkg/notification"
)

type NoOpProcessor struct{}

func (p *NoOpProcessor) Process(ctx context.Context, event notification.Message) error {
	return nil
}

// PublisherProcessor publishes incoming notification events to a Publisher (Redis broker).
type PublisherProcessor struct {
	publisher      notification.Publisher
	publishTimeout time.Duration
}

func NewPublisherProcessor(p notification.Publisher, timeout time.Duration) *PublisherProcessor {
	if p == nil {
		panic("publisher cannot be nil")
	}
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &PublisherProcessor{
		publisher:      p,
		publishTimeout: timeout,
	}
}

func (p *PublisherProcessor) Process(ctx context.Context, event notification.Message) error {
	pubCtx, cancel := context.WithTimeout(ctx, p.publishTimeout)
	defer cancel()

	if err := p.publisher.Publish(pubCtx, event.RecipientID, event); err != nil {
		logger.ErrorKV(pubCtx, "failed to publish notification", "recipient_id", event.RecipientID, "error", err)
		return err
	}

	logger.DebugKV(ctx, "published notification", "recipient_id", event.RecipientID, "event_id", event.EventID)
	return nil
}
