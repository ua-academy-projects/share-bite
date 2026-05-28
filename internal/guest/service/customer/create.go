package customer

import (
	"context"
	"fmt"
	"time"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
	"github.com/ua-academy-projects/share-bite/pkg/outbox"
)

func (s *service) Create(ctx context.Context, in entity.CreateCustomer) (string, error) {
	var customerID string
	err := s.txManager.ReadCommitted(ctx, func(txCtx context.Context) error {
		var err error
		customerID, err = s.customerRepo.Create(txCtx, in)
		if err != nil {
			return fmt.Errorf("create customer in repo: %w", err)
		}

		authToken, ok := txCtx.Value(middleware.CtxAccessToken).(string)
		if !ok || authToken == "" {
			return fmt.Errorf("missing access token to resolve customer email")
		}

		email, err := s.adminClient.GetUserEmail(txCtx, in.UserID, authToken)
		if err != nil {
			return fmt.Errorf("get customer email: %w", err)
		}

		metadata := map[string]any{
			"username": in.UserName,
			"email":    email,
		}

		event := outbox.Event{
			EventType: outbox.EventTypeRegistrationConfirmed,
			Payload: outbox.Message{
				EventID:     outbox.NewEventID(customerID, email),
				EventType:   outbox.EventTypeRegistrationConfirmed,
				RecipientID: in.UserID,
				ActorID:     in.UserID,
				EntityType:  "customer",
				EntityID:    customerID,
				Metadata:    metadata,
				CreatedAt:   time.Now().UTC(),
			},
			SourceService: outbox.DefaultSourceService,
		}

		if err := s.outboxWriter.Enqueue(txCtx, event); err != nil {
			return fmt.Errorf("failed to enqueue registration_confirmed outbox event: %w", err)
		}

		return nil
	})
	if err != nil {
		return "", fmt.Errorf("create customer in repo: %w", err)
	}

	return customerID, nil
}
