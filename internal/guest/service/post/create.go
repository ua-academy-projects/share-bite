package post

import (
	"context"
	"fmt"
	"github.com/ua-academy-projects/share-bite/internal/imageprocessing"
	"github.com/ua-academy-projects/share-bite/pkg/logger"

	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

func (s *service) Create(ctx context.Context, in dto.CreatePostInput) (entity.Post, error) {
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

	var post entity.Post

	err = s.txManager.ReadCommitted(
		ctx,
		func(txCtx context.Context) error {
			createdPost, err := s.createPostTx(
				txCtx,
				in,
				postImages,
			)
			if err != nil {
				return err
			}

			post = createdPost

			return nil
		},
	)

	if err != nil {
		rollbackUploadedImages(
			s.storage,
			uploadedKeys,
		)

		return entity.Post{}, fmt.Errorf(
			"execute post creation transaction: %w",
			err,
		)
	}

	if s.imageProcessingProducer != nil {
		for _, image := range post.Images {
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

				markErr := s.postRepo.MarkProcessingFailed(
					ctx,
					image.ID,
					"failed to enqueue image processing task",
				)

				if markErr != nil {
					logger.ErrorKV(
						ctx,
						"failed to mark image processing as failed",
						"image_id", image.ID,
						"error", markErr,
					)
				}
			}
		}
	}

	return post, nil
}
