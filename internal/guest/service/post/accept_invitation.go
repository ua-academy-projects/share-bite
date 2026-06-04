package post

import (
	"context"
	"fmt"
	"time"

	"github.com/ua-academy-projects/share-bite/pkg/outbox"
)

func (s *service) AcceptInvitation(ctx context.Context, collaboratorID string, customerID string) error {
	return s.txManager.ReadCommitted(
		ctx,
		func(txCtx context.Context) error {

			postID, err := s.postRepo.AcceptPostInvitation(
				txCtx,
				collaboratorID,
				customerID,
			)
			if err != nil {
				return fmt.Errorf(
					"accept invitation: %w",
					err,
				)
			}

			allAccepted, err := s.postRepo.TryPublishPostIfAllAccepted(
				txCtx,
				postID,
			)
			if err != nil {
				return fmt.Errorf(
					"try publish post: %w",
					err,
				)
			}

			if !allAccepted || s.outboxWriter == nil {
				return nil
			}

			collaborators, err := s.postRepo.GetAcceptedPostCollaborators(
				txCtx,
				postID,
			)
			if err != nil {
				return err
			}

			authorID, err := s.postRepo.GetAuthorCustomerID(
				txCtx,
				postID,
			)
			if err == nil {
				collaborators = append(
					collaborators,
					authorID,
				)
			}

			actor, err := s.customerRepo.GetByID(
				txCtx,
				customerID,
			)
			if err != nil {
				return fmt.Errorf(
					"get actor customer: %w",
					err,
				)
			}

			var actorAvatar string

			if actor.AvatarObjectKey != nil {
				actorAvatar = s.storage.BuildURL(
					*actor.AvatarObjectKey,
				)
			}

			seen := make(map[string]struct{})

			for _, collaboratorCustomerID := range collaborators {

				if _, ok := seen[collaboratorCustomerID]; ok {
					continue
				}

				seen[collaboratorCustomerID] = struct{}{}

				customer, err := s.customerRepo.GetByID(
					txCtx,
					collaboratorCustomerID,
				)
				if err != nil {
					continue
				}

				eventType := outbox.EventTypePostPublished

				eventID := outbox.NewEventID(
					eventType,
					customer.UserID,
					customerID,
					"post",
					postID,
				)

				evt := outbox.Message{
					EventID:     eventID,
					EventType:   eventType,
					RecipientID: customer.UserID,
					ActorID:     customerID,
					EntityType:  "post",
					EntityID:    postID,
					Metadata: map[string]any{
						"actor_avatar":   actorAvatar,
						"actor_username": actor.UserName,
					},
					CreatedAt: time.Now().UTC(),
				}

				if err := s.outboxWriter.Enqueue(
					txCtx,
					outbox.Event{
						EventType:     eventType,
						Payload:       evt,
						SourceService: outbox.DefaultSourceService,
					},
				); err != nil {
					return fmt.Errorf(
						"enqueue outbox event: %w",
						err,
					)
				}
			}

			return nil
		},
	)
}
