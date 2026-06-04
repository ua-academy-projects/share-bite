package follow

import (
	"context"
	"fmt"
	"time"

	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/pkg/outbox"
)

func (s *service) Follow(ctx context.Context, userID, targetCustomerID string) error {
	currentCustomer, err := s.customerRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if currentCustomer.ID == targetCustomerID {
		return apperror.ErrCannotFollowYourself
	}

	alreadyFollowing, err := s.customerFollowRepo.IsFollowing(ctx, currentCustomer.ID, targetCustomerID)
	if err != nil {
		return err
	}
	if alreadyFollowing {
		return nil
	}

	runFollow := func(txCtx context.Context) error {
		err := s.customerFollowRepo.Follow(txCtx, currentCustomer.ID, targetCustomerID)
		if err != nil {
			return err
		}

		if s.outboxWriter == nil || s.customerRepo == nil {
			return nil
		}

		targetCustomer, err := s.customerRepo.GetByID(txCtx, targetCustomerID)
		if err != nil {
			return err
		}

		if targetCustomer.UserID == "" {
			return nil
		}

		var actorAvatar string
		if currentCustomer.AvatarObjectKey != nil && s.storage != nil {
			actorAvatar = s.storage.BuildURL(*currentCustomer.AvatarObjectKey)
		}

		eventType := outbox.EventTypeFollowAdded
		eventID := outbox.NewEventID(
			eventType,
			targetCustomer.UserID,
			currentCustomer.ID,
		)

		evt := outbox.Message{
			EventID:     eventID,
			EventType:   eventType,
			RecipientID: targetCustomer.UserID,
			ActorID:     currentCustomer.ID,
			EntityType:  "customer",
			EntityID:    targetCustomerID,
			Metadata: map[string]any{
				"actor_avatar":   actorAvatar,
				"actor_username": currentCustomer.UserName,
			},
			CreatedAt: time.Now().UTC(),
		}

		if err := s.outboxWriter.Enqueue(txCtx, outbox.Event{
			EventType:     eventType,
			Payload:       evt,
			SourceService: outbox.DefaultSourceService,
		}); err != nil {
			return fmt.Errorf("enqueue follow outbox event: %w", err)
		}

		return nil
	}

	if s.txManager != nil {
		return s.txManager.ReadCommitted(ctx, runFollow)
	}
	return runFollow(ctx)
}
