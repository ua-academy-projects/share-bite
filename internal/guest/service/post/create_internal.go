package post

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	"github.com/ua-academy-projects/share-bite/pkg/outbox"

	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"

	"github.com/ua-academy-projects/share-bite/internal/storage/key"
	"github.com/ua-academy-projects/share-bite/internal/storage/mediatype"
)

const cleanupTimeout = 5 * time.Second

func (s *service) validateCreateInput(ctx context.Context, in dto.CreatePostInput) error {
	exists, err := s.venueProvider.CheckExists(
		ctx,
		in.VenueID,
	)
	if err != nil {
		return fmt.Errorf(
			"check venue exists: %w",
			err,
		)
	}

	if !exists {
		return apperror.VenueNotFoundID(
			in.VenueID,
		)
	}

	if len(in.Images) > 0 && s.storage == nil {
		return apperror.Internal(
			"storage is not configured",
		)
	}

	// mentions validation
	if len(in.Mentions) > 0 {
		mentions := UniqueStrings(in.Mentions)

		if len(mentions) > 10 {
			return apperror.BadRequest(
				"too many mentions",
			)
		}

		for _, mention := range mentions {
			if mention == in.CustomerID {
				return apperror.BadRequest(
					"cannot mention yourself",
				)
			}
		}

		followedIDs, err := s.followRepo.GetAllowedMentions(
			ctx,
			in.CustomerID,
			mentions,
		)
		if err != nil {
			return err
		}

		allowedSet := make(
			map[string]struct{},
			len(followedIDs),
		)

		for _, id := range followedIDs {
			allowedSet[id] = struct{}{}
		}

		for _, mention := range mentions {
			if _, ok := allowedSet[mention]; !ok {
				return apperror.ErrForbiddenMention
			}
		}

		in.Mentions = mentions
	}

	return nil
}

func (s *service) uploadPostImages(ctx context.Context, customerID string, images []dto.UploadImageInput) ([]entity.PostImage, []string, error) {
	if len(images) == 0 {
		return nil, nil, nil
	}

	uploadedKeys := make(
		[]string,
		0,
		len(images),
	)

	postImages := make(
		[]entity.PostImage,
		0,
		len(images),
	)

	uploadSessionID := uuid.NewString()

	for i, img := range images {
		ext, ok := mediatype.ExtFromContentType(
			img.ContentType,
		)
		if !ok {
			rollbackUploadedImages(
				s.storage,
				uploadedKeys,
			)

			return nil, nil,
				apperror.ErrUnsupportedImageType
		}

		objectKey := key.CustomerPostImageKey(
			customerID,
			uploadSessionID,
			uuid.NewString(),
			ext,
		)

		if err := s.storage.Upload(
			ctx,
			objectKey,
			img.ContentType,
			img.File,
		); err != nil {

			rollbackUploadedImages(
				s.storage,
				uploadedKeys,
			)

			return nil, nil, fmt.Errorf(
				"upload post image to storage: %w",
				err,
			)
		}

		uploadedKeys = append(
			uploadedKeys,
			objectKey,
		)

		postImages = append(postImages, entity.PostImage{
			ObjectKey:   objectKey,
			ContentType: img.ContentType,
			FileSize:    img.FileSize,
			SortOrder:   int16(i),
		})
	}

	return postImages, uploadedKeys, nil
}

func (s *service) createPostTx(ctx context.Context, in dto.CreatePostInput, postImages []entity.PostImage) (entity.Post, error) {
	createdPost, err := s.postRepo.Create(
		ctx,
		in,
	)
	if err != nil {
		return entity.Post{}, fmt.Errorf(
			"create post in post repository: %w",
			err,
		)
	}

	// attach images
	if len(postImages) > 0 {
		for i := range postImages {
			postImages[i].PostID = createdPost.ID
		}

		createdImages, err := s.postRepo.CreateImages(
			ctx,
			postImages,
		)
		if err != nil {
			return entity.Post{}, fmt.Errorf(
				"create post images in post repository: %w",
				err,
			)
		}

		createdPost.Images = createdImages
	}

	// create mentions
	if len(in.Mentions) > 0 {
		mentions := make(
			[]entity.PostMention,
			0,
			len(in.Mentions),
		)

		for _, mention := range in.Mentions {
			mentions = append(
				mentions,
				entity.PostMention{
					PostID:     createdPost.ID,
					CustomerID: mention,
				},
			)
		}

		if err := s.postRepo.CreateMentions(
			ctx,
			mentions,
		); err != nil {
			return entity.Post{}, fmt.Errorf(
				"create mentions: %w",
				err,
			)
		}

		if s.outboxWriter != nil && s.customerRepo != nil {
			actor, err := s.customerRepo.GetByID(ctx, in.CustomerID)
			if err != nil {
				return entity.Post{}, fmt.Errorf("get actor customer: %w", err)
			}

			var actorAvatar string
			if actor.AvatarObjectKey != nil {
				actorAvatar = s.storage.BuildURL(*actor.AvatarObjectKey)
			}

			for _, mention := range in.Mentions {
				mentionedCustomer, err := s.customerRepo.GetByID(ctx, mention)
				if err != nil {
					logger.ErrorKV(ctx, "failed to get mentioned customer, skipping mention notification", "customer_id", mention, "error", err)
					continue
				}

				if mentionedCustomer.UserID == "" {
					continue
				}

				eventType := outbox.EventTypePostMentioned
				eventID := outbox.NewEventID(eventType, mentionedCustomer.UserID, in.CustomerID, "post", createdPost.ID)

				evt := outbox.Message{
					EventID:     eventID,
					EventType:   eventType,
					RecipientID: mentionedCustomer.UserID,
					ActorID:     in.CustomerID,
					EntityType:  "post",
					EntityID:    createdPost.ID,
					Metadata: map[string]any{
						"actor_avatar":   actorAvatar,
						"actor_username": actor.UserName,
					},
					CreatedAt: time.Now().UTC(),
				}

				if err := s.outboxWriter.Enqueue(ctx, outbox.Event{
					EventType:     eventType,
					Payload:       evt,
					SourceService: outbox.DefaultSourceService,
				}); err != nil {
					return entity.Post{}, fmt.Errorf("enqueue mention outbox event: %w", err)
				}
			}
		}
	}

	return createdPost, nil
}
