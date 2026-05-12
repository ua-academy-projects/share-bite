package post

import (
	"context"
	"fmt"

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

	return post, nil
}
