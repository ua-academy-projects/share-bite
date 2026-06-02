package service

import (
	"context"

	"github.com/ua-academy-projects/share-bite/pkg/notification"
)

type Processor struct {
	svc *Service
}

func NewProcessor(svc *Service) *Processor {
	return &Processor{svc: svc}
}

func (p *Processor) ProcessMessage(ctx context.Context, msg notification.Message) error {
	return p.svc.ProcessMessage(ctx, msg)
}
