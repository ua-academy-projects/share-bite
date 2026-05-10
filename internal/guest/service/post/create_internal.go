package post

import (
	"context"
	"fmt"
	"time"

	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/storage"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

const cleanupTimeout = 5 * time.Second

func (s *service) createPost(ctx context.Context, in dto.CreatePostInput) (entity.Post, error) {
	exists, err := s.venueProvider.CheckExists(ctx, in.VenueID)
	if err != nil {
		return entity.Post{}, fmt.Errorf("check venue exists: %w", err)
	}

	if !exists {
		return entity.Post{}, apperror.VenueNotFoundID(in.VenueID)
	}

	if len(in.Images) > 0 && s.storage == nil {
		return entity.Post{}, apperror.Internal("storage is not configured")
	}

	// mentions validation
	if len(in.Mentions) > 0 {
		mentions := UniqueStrings(in.Mentions)

		if len(mentions) > 10 {
			return entity.Post{}, apperror.BadRequest("too many mentions")
		}

		for _, m := range mentions {
			if m == in.CustomerID {
				return entity.Post{}, apperror.BadRequest("cannot mention yourself")
			}
		}

		followedIDs, err := s.followRepo.GetAllowedMentions(
			ctx,
			in.CustomerID,
			mentions,
		)
		if err != nil {
			return entity.Post{}, err
		}

		allowedSet := make(map[string]struct{}, len(followedIDs))
		for _, id := range followedIDs {
			allowedSet[id] = struct{}{}
		}

		for _, m := range mentions {
			if _, ok := allowedSet[m]; !ok {
				return entity.Post{}, apperror.ErrForbiddenMention
			}
		}

		in.Mentions = mentions
	}

	postImages, uploadedKeys, err := s.uploadPostImages(
		ctx,
		in.CustomerID,
		in.Images,
	)
	if err != nil {
		return entity.Post{}, err
	}

	createdPost, err := s.postRepo.Create(ctx, in)
	if err != nil {
		rollbackUploadedImages(s.storage, uploadedKeys)

		return entity.Post{}, fmt.Errorf(
			"create post in repository: %w",
			err,
		)
	}

	// images
	if len(postImages) > 0 {
		for i := range postImages {
			postImages[i].PostID = createdPost.ID
		}

		if err := s.postRepo.CreateImages(ctx, postImages); err != nil {
			rollbackUploadedImages(s.storage, uploadedKeys)

			return entity.Post{}, fmt.Errorf(
				"create post images: %w",
				err,
			)
		}

		createdPost.Images = postImages
	}

	// mentions
	if len(in.Mentions) > 0 {
		mentions := make([]entity.PostMention, 0, len(in.Mentions))

		for _, m := range in.Mentions {
			mentions = append(mentions, entity.PostMention{
				PostID:     createdPost.ID,
				CustomerID: m,
			})
		}

		if err := s.postRepo.CreateMentions(ctx, mentions); err != nil {
			rollbackUploadedImages(s.storage, uploadedKeys)

			return entity.Post{}, fmt.Errorf(
				"create mentions: %w",
				err,
			)
		}
	}

	return createdPost, nil
}

func rollbackUploadedImages(objectStorage storage.ObjectStorage, keys []string) {
	if objectStorage == nil || len(keys) == 0 {
		return
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		cleanupTimeout,
	)
	defer cancel()

	for _, key := range keys {
		cleanupDelete(ctx, objectStorage, key)
	}
}

func cleanupDelete(ctx context.Context, objectStorage storage.ObjectStorage, key string) {
	if key == "" {
		return
	}

	if err := objectStorage.Delete(ctx, key); err != nil {
		logger.WarnKV(
			ctx,
			"failed to cleanup post image object",
			"key",
			key,
			"error",
			err,
		)
	}
}
