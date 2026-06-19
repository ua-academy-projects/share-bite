package comment

import (
	"context"
	"fmt"
	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	"github.com/ua-academy-projects/share-bite/pkg/outbox"
	"strconv"
	"time"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

func (s *service) Create(ctx context.Context, in dto.CreateCommentInput) (entity.Comment, error) {
	postIDStr := strconv.FormatInt(in.PostID, 10)

	var createdComment entity.Comment

	runCreate := func(txCtx context.Context) error {
		post, err := s.postSvc.Get(txCtx, postIDStr, "")
		if err != nil {
			return fmt.Errorf("check post existence in comment service: %w", err)
		}

		comment, err := s.commentRepo.Create(txCtx, in)
		if err != nil {
			return fmt.Errorf("create comment in repo: %w", err)
		}

		createdComment = comment

		if s.outboxWriter == nil || s.customerRepo == nil || post.CustomerID == "" || post.CustomerID == in.CustomerID {
			return nil
		}

		recipient, err := s.customerRepo.GetByID(txCtx, post.CustomerID)
		if err != nil {
			logger.ErrorKV(txCtx, "failed to get post author for comment notification, skipping", "post_id", in.PostID, "error", err)
			return nil
		}

		if recipient.UserID == "" {
			return nil
		}

		actor, err := s.customerRepo.GetByID(txCtx, in.CustomerID)
		if err != nil {
			return fmt.Errorf("get actor customer: %w", err)
		}

		var actorAvatar string
		if actor.AvatarObjectKey != nil && s.storage != nil {
			actorAvatar = s.storage.BuildURL(*actor.AvatarObjectKey)
		}

		eventType := outbox.EventTypePostCommented
		eventID := outbox.NewEventID(
			eventType,
			recipient.UserID,
			in.CustomerID,
			"comment",
			strconv.FormatInt(comment.ID, 10),
		)

		evt := outbox.Message{
			EventID:     eventID,
			EventType:   eventType,
			RecipientID: recipient.UserID,
			ActorID:     in.CustomerID,
			EntityType:  "comment",
			EntityID:    strconv.FormatInt(comment.ID, 10),
			Metadata: map[string]any{
				"post_id":         postIDStr,
				"actor_avatar":    actorAvatar,
				"actor_username":  actor.UserName,
				"comment_preview": comment.Text,
			},
			CreatedAt: time.Now().UTC(),
		}

		if err := s.outboxWriter.Enqueue(txCtx, outbox.Event{
			EventType:     eventType,
			Payload:       evt,
			SourceService: outbox.DefaultSourceService,
		}); err != nil {
			return fmt.Errorf("enqueue comment outbox event: %w", err)
		}

		return nil
	}

	if s.txManager != nil {
		if err := s.txManager.ReadCommitted(ctx, runCreate); err != nil {
			return entity.Comment{}, err
		}
	} else {
		if err := runCreate(ctx); err != nil {
			return entity.Comment{}, err
		}
	}

	return createdComment, nil
}
