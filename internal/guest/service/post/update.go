package post

import (
	"context"
	"errors"
	"fmt"
	"github.com/ua-academy-projects/share-bite/internal/imageprocessing"
	"github.com/ua-academy-projects/share-bite/pkg/logger"

	"github.com/google/uuid"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/storage/key"
	"github.com/ua-academy-projects/share-bite/internal/storage/mediatype"
)

func (s *service) Update(ctx context.Context, in entity.UpdatePostInput) (entity.Post, error) {
	currentPost, err := s.postRepo.GetByID(ctx, in.ID)
	if err != nil {
		return entity.Post{}, fmt.Errorf(
			"get post by id in post repository: %w",
			err,
		)
	}

	isOwner := currentPost.CustomerID == in.CustomerID
	canEdit := isOwner

	if !canEdit {
		canEdit, err = s.postRepo.IsAcceptedCollaborator(
			ctx,
			currentPost.ID,
			in.CustomerID,
		)
		if err != nil {
			return entity.Post{}, fmt.Errorf(
				"check accepted collaborator: %w",
				err,
			)
		}
	}

	if !canEdit {
		return entity.Post{}, apperror.PostNotFoundID(in.ID)
	}

	if !isOwner {
		if in.Status != nil {
			return entity.Post{}, apperror.PostChangesNotAllowed("status")
		}

		if in.Rating != nil {
			return entity.Post{}, apperror.PostChangesNotAllowed("rating")
		}

		if in.VenueID != nil {
			return entity.Post{}, apperror.PostChangesNotAllowed("venue")
		}
	}

	nextStatus := currentPost.Status
	if in.Status != nil {
		nextStatus = *in.Status
	}

	if !isValidPostStatusTransition(
		currentPost.Status,
		nextStatus,
	) {
		return entity.Post{}, apperror.InvalidPostStatusTransition(
			string(currentPost.Status),
			string(nextStatus),
		)
	}

	if in.VenueID != nil {
		exists, err := s.venueProvider.CheckExists(
			ctx,
			*in.VenueID,
		)
		if err != nil {
			return entity.Post{}, fmt.Errorf(
				"check venue existence via venue provider: %w: %w",
				apperror.ErrUpstreamError,
				err,
			)
		}

		if !exists {
			return entity.Post{}, apperror.VenueNotFoundID(
				*in.VenueID,
			)
		}
	}

	if !in.RewriteImages {
		post, err := s.postRepo.Update(ctx, in)
		if err != nil {
			return entity.Post{}, fmt.Errorf(
				"update post in post repository: %w",
				err,
			)
		}

		return post, nil
	}

	if s.storage == nil {
		return entity.Post{}, apperror.Internal(
			"storage is not configured",
		)
	}

	if s.txManager == nil {
		return entity.Post{}, apperror.Internal(
			"transaction manager is not configured",
		)
	}

	var preservedImages []entity.PostImage
	var keysToDelete []string

	keptMap := make(map[string]bool, len(in.KeptImages))
	for _, key := range in.KeptImages {
		keptMap[key] = true
	}

	for _, oldImg := range currentPost.Images {
		if keptMap[oldImg.ObjectKey] {
			preservedImages = append(
				preservedImages,
				oldImg,
			)
		} else {
			keysToDelete = append(
				keysToDelete,
				oldImg.ObjectKey,
			)
		}
	}

	var uploadedKeys []string

	newImages := make(
		[]entity.PostImage,
		0,
		len(preservedImages)+len(in.Images),
	)

	uploadSessionID := uuid.NewString()

	// Preserve existing images
	for i, img := range preservedImages {
		newImages = append(newImages, entity.PostImage{
			ObjectKey:   img.ObjectKey,
			ContentType: img.ContentType,
			FileSize:    img.FileSize,
			SortOrder:   int16(i),
		})
	}

	// Upload new images
	for i, img := range in.Images {
		ext, ok := mediatype.ExtFromContentType(
			img.ContentType,
		)
		if !ok {
			rollbackUploadedImages(
				s.storage,
				uploadedKeys,
			)

			return entity.Post{}, apperror.ErrUnsupportedImageType
		}

		objectKey := key.CustomerPostImageKey(
			in.CustomerID,
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

			return entity.Post{}, fmt.Errorf(
				"upload post image to storage: %w",
				err,
			)
		}

		uploadedKeys = append(
			uploadedKeys,
			objectKey,
		)

		newImage := entity.PostImage{
			ObjectKey:   objectKey,
			ContentType: img.ContentType,
			FileSize:    img.FileSize,
			SortOrder: int16(
				len(preservedImages) + i,
			),
		}
		newImages = append(newImages, newImage)
	}

	var post entity.Post

	err = s.txManager.ReadCommitted(
		ctx,
		func(txCtx context.Context) error {
			updatedPost, err := s.postRepo.Update(
				txCtx,
				in,
			)
			if err != nil {
				if errors.Is(
					err,
					apperror.ErrEmptyUpdate,
				) && in.RewriteImages {

					updatedPost, err = s.postRepo.GetByID(
						txCtx,
						in.ID,
					)
					if err != nil {
						return fmt.Errorf(
							"get post by id in post repository: %w",
							err,
						)
					}
				} else {
					return fmt.Errorf(
						"update post in post repository: %w",
						err,
					)
				}
			}

			if err := s.postRepo.DeleteImagesByPostID(
				txCtx,
				updatedPost.ID,
			); err != nil {
				return fmt.Errorf(
					"delete post images in post repository: %w",
					err,
				)
			}

			if len(newImages) > 0 {
				for i := range newImages {
					newImages[i].PostID = updatedPost.ID
				}
				createdImages, err := s.postRepo.CreateImages(
					txCtx,
					newImages,
				)
				if err != nil {
					return fmt.Errorf(
						"create post images in post repository: %w",
						err,
					)
				}
				updatedPost.Images = createdImages
			}

			post = updatedPost

			return nil
		},
	)
	if err != nil {
		rollbackUploadedImages(
			s.storage,
			uploadedKeys,
		)

		return entity.Post{}, fmt.Errorf(
			"execute post update transaction: %w",
			err,
		)
	}

	for _, key := range keysToDelete {
		cleanupDelete(s.storage, key)
	}
	if s.imageProcessingProducer != nil {
		for _, image := range post.Images {
			if keptMap[image.ObjectKey] {
				continue
			}

			err := s.imageProcessingProducer.SendMessage(
				ctx,
				imageprocessing.ProcessImageMessage{
					ImageID: image.ID,
					S3Key:   image.ObjectKey,
				},
			)
			if err != nil {
				logger.ErrorKV(
					ctx,
					"failed to send image processing message",
					"image_id", image.ID,
					"object_key", image.ObjectKey,
					"error", err,
				)
			}
		}
	}

	return post, nil
}
