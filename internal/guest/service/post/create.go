package post

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
	"github.com/ua-academy-projects/share-bite/internal/storage"
	"github.com/ua-academy-projects/share-bite/internal/storage/key"
	"github.com/ua-academy-projects/share-bite/internal/storage/mediatype"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	"time"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

func (s *service) Create(ctx context.Context, in dto.CreatePostInput) (entity.Post, error) {
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

		followedIDs, err := s.followRepo.GetAllowedMentions(ctx, in.CustomerID, mentions)
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

	// upload images BEFORE tx
	var uploadedKeys []string
	postImages := make([]entity.PostImage, 0, len(in.Images))
	uploadSessionID := uuid.New().String()

	for i, img := range in.Images {
		ext, ok := mediatype.ExtFromContentType(img.ContentType)
		if !ok {
			rollbackUploadedImages(s.storage, uploadedKeys)
			return entity.Post{}, apperror.ErrUnsupportedImageType
		}

		objectKey := key.CustomerPostImageKey(in.CustomerID, uploadSessionID, uuid.NewString(), ext)
		if err := s.storage.Upload(ctx, objectKey, img.ContentType, img.File); err != nil {
			rollbackUploadedImages(s.storage, uploadedKeys)
			return entity.Post{}, fmt.Errorf("upload post image to storage: %w", err)
		}

		uploadedKeys = append(uploadedKeys, objectKey)
		postImages = append(postImages, entity.PostImage{
			ObjectKey:   objectKey,
			ContentType: img.ContentType,
			FileSize:    img.FileSize,
			SortOrder:   int16(i),
		})
	}

	var post entity.Post

	err = s.txManager.ReadCommitted(ctx, func(txCtx context.Context) error {
		createdPost, err := s.postRepo.Create(txCtx, in)
		if err != nil {
			return fmt.Errorf("create post in post repository: %w", err)
		}

		// attach images
		if len(postImages) > 0 {
			for i := range postImages {
				postImages[i].PostID = createdPost.ID
			}

			if err := s.postRepo.CreateImages(txCtx, postImages); err != nil {
				return fmt.Errorf("create post images in post repository: %w", err)
			}

			createdPost.Images = postImages
		}

		// create mentions
		if len(in.Mentions) > 0 {
			mentions := make([]entity.PostMention, 0, len(in.Mentions))
			for _, m := range in.Mentions {
				mentions = append(mentions, entity.PostMention{
					PostID:     createdPost.ID,
					CustomerID: m,
				})
			}

			if err := s.postRepo.CreateMentions(txCtx, mentions); err != nil {
				return fmt.Errorf("create mentions: %w", err)
			}
		}

		post = createdPost
		return nil
	})

	if err != nil {
		rollbackUploadedImages(s.storage, uploadedKeys)
		return entity.Post{}, fmt.Errorf("execute post creation transaction: %w", err)
	}

	return post, nil
}

func rollbackUploadedImages(storage storage.ObjectStorage, keys []string) {
	for _, key := range keys {
		cleanupDelete(storage, key)
	}
}

func cleanupDelete(objectStorage storage.ObjectStorage, key string) {
	if objectStorage == nil || key == "" {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := objectStorage.Delete(ctx, key); err != nil {
		logger.WarnKV(ctx, "failed to cleanup post image object",
			"key", key,
			"error", err,
		)
	}
}

func UniqueStrings(input []string) []string {
	set := make(map[string]struct{}, len(input))
	res := make([]string, 0, len(input))
	for _, v := range input {
		if _, ok := set[v]; ok {
			continue
		}
		set[v] = struct{}{}
		res = append(res, v)
	}
	return res
}
