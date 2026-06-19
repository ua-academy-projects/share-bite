package post

import (
	"context"
	"fmt"
	"time"

	"github.com/ua-academy-projects/share-bite/pkg/logger"
	"github.com/ua-academy-projects/share-bite/pkg/outbox"
)

func (s *service) Like(ctx context.Context, postID string, customerID string) error {
	post, err := s.Get(ctx, postID, customerID)
	if err != nil {
		return fmt.Errorf("validate post for like: %w", err)
	}

	return s.txManager.ReadCommitted(ctx, func(txCtx context.Context) error {
		inserted, err := s.postRepo.Like(txCtx, postID, customerID)
		if err != nil {
			return fmt.Errorf("like post in repository: %w", err)
		}
		if !inserted {
			logger.DebugKV(txCtx, "like already exists, skip outbox enqueue", "post_id", postID, "customer_id", customerID)
			return nil
		}

		if s.outboxWriter == nil || post.CustomerID == "" || post.CustomerID == customerID {
			return nil
		}

		authorUserID, err := s.postRepo.GetAuthorUserID(txCtx, postID)
		if err != nil {
			logger.ErrorKV(txCtx, "failed to get author user ID for notification, skipping", "post_id", postID, "error", err)
			return nil
		}
		if authorUserID == "" {
			return nil
		}

		// Get actor's profile for notification enrichment
		actor, err := s.customerRepo.GetByID(txCtx, customerID)
		if err != nil {
			return fmt.Errorf("get actor customer: %w", err)
		}

		var actorAvatar string
		if actor.AvatarObjectKey != nil {
			actorAvatar = s.storage.BuildURL(*actor.AvatarObjectKey)
		}

		eventType := outbox.EventTypePostLiked
		eventID := outbox.NewEventID(eventType, authorUserID, customerID, "post", postID)

		evt := outbox.Message{
			EventID:     eventID,
			EventType:   eventType,
			RecipientID: authorUserID,
			ActorID:     customerID,
			EntityType:  "post",
			EntityID:    postID,
			Metadata: map[string]any{
				"actor_avatar":   actorAvatar,
				"actor_username": actor.UserName,
			},
			CreatedAt: time.Now().UTC(),
		}

		if err := s.outboxWriter.Enqueue(txCtx, outbox.Event{
			EventType:     eventType,
			Payload:       evt,
			SourceService: outbox.DefaultSourceService,
		}); err != nil {
			return fmt.Errorf("enqueue outbox event: %w", err)
		}

		return nil
	})
}

func (s *service) Unlike(ctx context.Context, postID string, customerID string) error {
	if _, err := s.Get(ctx, postID, customerID); err != nil {
		return fmt.Errorf("validate post for unlike: %w", err)
	}

	if err := s.postRepo.Unlike(ctx, postID, customerID); err != nil {
		return fmt.Errorf("unlike post in repository: %w", err)
	}
	return nil
}
