package customer

import (
	"context"
	"fmt"
	"time"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/pkg/outbox"
)

func (s *service) Create(ctx context.Context, in entity.CreateCustomer) (string, error) {
	var customerID string
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var err error
		customerID, err = s.customerRepo.Create(ctx, in)
		if err != nil {
			return fmt.Errorf("create customer in repo: %w", err)
		}

		if s.outboxWriter != nil {
			event := outbox.Event{
				EventType: outbox.EventTypeRegistrationConfirmed,
				Payload: outbox.Message{
					EventID:     outbox.NewEventID(customerID, in.Email),
					EventType:   outbox.EventTypeRegistrationConfirmed,
					RecipientID: in.UserID,
					// TODO: Enrich email from admin-auth endpoint instead of passing it here.
					// This avoids leaking email in the customer struct and keeps auth data centralized.
					Metadata: map[string]any{
						"email":    in.Email,
						"username": in.UserName,
					},
					CreatedAt: time.Now().UTC(),
				},
				SourceService: outbox.DefaultSourceService,
			}

			if err := s.outboxWriter.Enqueue(ctx, event); err != nil {
				return fmt.Errorf("failed to enqueue registration_confirmed outbox event: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return customerID, nil
}
