package post

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

func (s *service) uploadPostImages(ctx context.Context, customerID string, images []dto.UploadImageInput) ([]entity.PostImage, []string, error) {
	if len(images) == 0 {
		return nil, nil, nil
	}

	uploadedKeys := make([]string, 0, len(images))
	postImages := make([]entity.PostImage, 0, len(images))

	uploadSessionID := uuid.New().String()

	for i, img := range images {
		ext := extensionFromContentType(img.ContentType)
		if ext == "" {
			rollbackUploadedImages(s.storage, uploadedKeys)
			return nil, nil, apperror.ErrUnsupportedImageType
		}

		objectKey := generatePostImageKey(customerID, uploadSessionID, ext)

		uploadedKey, err := s.storage.Upload(
			ctx,
			objectKey,
			img.ContentType,
			img.File,
		)
		if err != nil {
			rollbackUploadedImages(s.storage, uploadedKeys)

			return nil, nil, fmt.Errorf(
				"upload post image to storage: %w",
				err,
			)
		}

		uploadedKeys = append(uploadedKeys, uploadedKey)

		postImages = append(postImages, entity.PostImage{
			ObjectKey:   uploadedKey,
			ContentType: img.ContentType,
			FileSize:    img.FileSize,
			SortOrder:   int16(i),
		})
	}

	return postImages, uploadedKeys, nil
}
