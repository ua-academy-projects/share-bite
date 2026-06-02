package post

import (
	"context"
	"fmt"
	"github.com/ua-academy-projects/share-bite/pkg/outbox"
	"time"

	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

const (
	postInvitationTTL          = 24 * time.Hour
	notificationPublishTimeout = 5 * time.Second
)

func (s *service) CreatePostWithCollaborators(ctx context.Context, in dto.CreatePostInput) (entity.Post, error) {
	if err := s.validateCreateInput(
		ctx,
		in,
	); err != nil {
		return entity.Post{}, err
	}

	postImages, uploadedKeys, err := s.uploadPostImages(
		ctx,
		in.CustomerID,
		in.Images,
	)
	if err != nil {
		return entity.Post{}, err
	}

	var result entity.Post

	err = s.txManager.ReadCommitted(
		ctx,
		func(txCtx context.Context) error {

			post, err := s.createPostTx(
				txCtx,
				in,
				postImages,
			)
			if err != nil {
				return err
			}

			result = post

			invited := uniqueAndExcludeSelf(
				in.CustomerID,
				in.InvitedCustomerIDs,
			)

			// no collaborators -> publish immediately
			if len(invited) == 0 {

				if err := s.postRepo.UpdateStatus(
					txCtx,
					post.ID,
					in.CustomerID,
					entity.PostStatusPublished,
				); err != nil {
					return err
				}

				updatedPost, err := s.postRepo.GetByID(
					txCtx,
					post.ID,
				)
				if err != nil {
					return err
				}

				result = updatedPost

				return nil
			}

			expiresAt := time.Now().Add(
				postInvitationTTL,
			)

			if err := s.postRepo.CreatePostCollaborators(
				txCtx,
				post.ID,
				in.CustomerID,
				invited,
				expiresAt,
			); err != nil {
				return err
			}

			// enqueue notifications
			if s.outboxWriter != nil {

				actor, err := s.customerRepo.GetByID(
					txCtx,
					in.CustomerID,
				)
				if err != nil {
					return fmt.Errorf(
						"get inviter customer: %w",
						err,
					)
				}

				actorName := actor.UserName

				if actor.FirstName != "" || actor.LastName != "" {
					actorName = fmt.Sprintf(
						"%s %s",
						actor.FirstName,
						actor.LastName,
					)
				}

				var actorAvatar string

				if actor.AvatarObjectKey != nil && s.storage != nil {
					actorAvatar = s.storage.BuildURL(
						*actor.AvatarObjectKey,
					)
				}

				for _, collaboratorID := range invited {

					customer, err := s.customerRepo.GetByID(
						txCtx,
						collaboratorID,
					)
					if err != nil {
						continue
					}

					eventType := "post_invitation_received"

					eventID := outbox.NewEventID(
						eventType,
						customer.UserID,
						in.CustomerID,
						"post",
						post.ID,
					)

					evt := outbox.Message{
						EventID:     eventID,
						EventType:   eventType,
						RecipientID: customer.UserID,
						ActorID:     in.CustomerID,
						EntityType:  "post",
						EntityID:    post.ID,
						Metadata: map[string]any{
							"inviter_customer_id": in.CustomerID,
							"actor_name":          actorName,
							"actor_avatar":        actorAvatar,
							"actor_username":      actor.UserName,
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
			}

			return nil
		},
	)

	if err != nil {
		rollbackUploadedImages(
			s.storage,
			uploadedKeys,
		)

		return entity.Post{}, fmt.Errorf(
			"execute collaborative post creation transaction: %w",
			err,
		)
	}

	return result, nil
}
