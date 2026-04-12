package post

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/ua-academy-projects/share-bite/internal/storage"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	"time"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

func (s *service) Create(ctx context.Context, in entity.CreatePostInput) (entity.Post, error) {
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

	var uploadedKeys []string

	var post entity.Post

	err = s.txManager.ReadCommitted(ctx, func(txCtx context.Context) error {
		createdPost, err := s.postRepo.Create(txCtx, in)
		if err != nil {
			return fmt.Errorf("create post in post repository: %w", err)
		}

		postImages := make([]entity.PostImage, 0, len(in.Images))
		for i, img := range in.Images {
			ext := extensionFromContentType(img.ContentType)
			if ext == "" {
				return apperror.ErrUnsupportedImageType
			}

			objectKey := generatePostImageKey(in.CustomerID, createdPost.ID, ext)
			uploadedKey, err := s.storage.Upload(ctx, objectKey, img.ContentType, img.File)
			if err != nil {
				return fmt.Errorf("upload post image to storage: %w", err)
			}

			uploadedKeys = append(uploadedKeys, uploadedKey)
			postImages = append(postImages, entity.PostImage{
				PostID:      createdPost.ID,
				ObjectKey:   uploadedKey,
				ContentType: img.ContentType,
				FileSize:    img.FileSize,
				SortOrder:   int16(i),
			})
		}

		if len(postImages) > 0 {
			if err := s.postRepo.CreateImages(txCtx, postImages); err != nil {
				return fmt.Errorf("create post images in post repository: %w", err)
			}

			createdPost.Images = postImages
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

func generatePostImageKey(customerID, postID, ext string) string {
	return fmt.Sprintf("posts/%s/%s/%s.%s", customerID, postID, uuid.New().String(), ext)
}

func extensionFromContentType(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return "jpg"
	case "image/png":
		return "png"
	default:
		return ""
	}
}
