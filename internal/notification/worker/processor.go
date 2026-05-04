package worker

import (
	"context"

	"github.com/ua-academy-projects/share-bite/pkg/notification"
)

type NoOpProcessor struct{}

func (p *NoOpProcessor) Process(ctx context.Context, event notification.Message) error {
	return nil
}
