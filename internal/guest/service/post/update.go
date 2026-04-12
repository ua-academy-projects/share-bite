package post

import (
	"context"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

func (s *service) Update(ctx context.Context, in entity.UpdatePostInput) (entity.Post, error) {
	currentPost, err := s.postRepo.GetByID(ctx, in.ID)
	if err != nil {
		return entity.Post{}, fmt.Errorf("get post by id in post repository: %w", err)
	}

	if currentPost.CustomerID != in.CustomerID {
		return entity.Post{}, apperror.PostNotFoundID(in.ID)
	}

	nextStatus := currentPost.Status
	if in.Status != nil {
		nextStatus = *in.Status
	}

	if !isValidPostStatusTransition(currentPost.Status, nextStatus) {
		return entity.Post{}, apperror.InvalidPostStatusTransition(string(currentPost.Status), string(nextStatus))
	}

	if in.VenueID != nil {
		exists, err := s.venueProvider.CheckExists(ctx, *in.VenueID)
		if err != nil {
			return entity.Post{}, fmt.Errorf("check venue existence via venue provider: %w: %w", apperror.ErrUpstreamError, err)
		}
		if !exists {
			return entity.Post{}, apperror.VenueNotFoundID(*in.VenueID)
		}
	}

	if !in.RewriteImages {
		post, err := s.postRepo.Update(ctx, in)
		if err != nil {
			return entity.Post{}, fmt.Errorf("update post in post repository: %w", err)
		}

		return post, nil
	}

	if s.storage == nil {
		return entity.Post{}, apperror.Internal("storage is not configured")
	}

	if s.txManager == nil {
		return entity.Post{}, apperror.Internal("transaction manager is not configured")
	}

	var uploadedKeys []string
	newImages := make([]entity.PostImage, 0, len(in.Images))

	for i, img := range in.Images {
		ext := extensionFromContentType(img.ContentType)
		if ext == "" {
			return entity.Post{}, apperror.ErrUnsupportedImageType
		}

		objectKey := generatePostImageKey(in.CustomerID, in.ID, ext)
		uploadedKey, err := s.storage.Upload(ctx, objectKey, img.ContentType, img.File)
		if err != nil {
			rollbackUploadedImages(s.storage, uploadedKeys)
			return entity.Post{}, fmt.Errorf("upload post image to storage: %w", err)
		}

		uploadedKeys = append(uploadedKeys, uploadedKey)
		newImages = append(newImages, entity.PostImage{
			ObjectKey:   uploadedKey,
			ContentType: img.ContentType,
			FileSize:    img.FileSize,
			SortOrder:   int16(i),
		})
	}

	oldKeys := extractImageKeys(currentPost.Images)

	var post entity.Post
	err = s.txManager.ReadCommitted(ctx, func(txCtx context.Context) error {
		updatedPost, err := s.postRepo.Update(txCtx, in)
		if err != nil {
			return fmt.Errorf("update post in post repository: %w", err)
		}

		if err := s.postRepo.DeleteImagesByPostID(txCtx, updatedPost.ID); err != nil {
			return fmt.Errorf("delete post images in post repository: %w", err)
		}

		if len(newImages) > 0 {
			for i := range newImages {
				newImages[i].PostID = updatedPost.ID
			}

			if err := s.postRepo.CreateImages(txCtx, newImages); err != nil {
				return fmt.Errorf("create post images in post repository: %w", err)
			}
		}

		updatedPost.Images = newImages
		post = updatedPost

		return nil
	})
	if err != nil {
		rollbackUploadedImages(s.storage, uploadedKeys)
		return entity.Post{}, fmt.Errorf("execute post update transaction: %w", err)
	}

	for _, key := range oldKeys {
		cleanupDelete(s.storage, key)
	}

	post.Images = newImages

	return post, nil
}

func extractImageKeys(images []entity.PostImage) []string {
	keys := make([]string, 0, len(images))
	for _, image := range images {
		if image.ObjectKey == "" {
			continue
		}

		keys = append(keys, image.ObjectKey)
	}

	return keys
}

func isValidPostStatusTransition(from, to entity.PostStatus) bool {
	switch from {
	case entity.PostStatusDraft:
		return to == entity.PostStatusDraft ||
			to == entity.PostStatusPublished ||
			to == entity.PostStatusArchived
	case entity.PostStatusPublished:
		return to == entity.PostStatusPublished ||
			to == entity.PostStatusArchived
	case entity.PostStatusArchived:
		return to == entity.PostStatusArchived ||
			to == entity.PostStatusPublished
	default:
		return false
	}
}
