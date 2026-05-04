package post

import (
	"context"
	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"time"
)

func (s *service) CreatePostWithCollaborators(ctx context.Context, in dto.CreatePostInput) (entity.Post, error) {
	var result entity.Post

	err := s.txManager.ReadCommitted(ctx, func(txCtx context.Context) error {

		post, err := s.postRepo.Create(txCtx, in)
		if err != nil {
			return err
		}

		result = post

		invited := uniqueAndExcludeSelf(in.CustomerID, in.InvitedCustomerIDs)

		if len(invited) == 0 {
			return s.postRepo.UpdateStatus(
				txCtx,
				post.ID,
				in.CustomerID,
				entity.PostStatusPublished,
			)
		}

		expiresAt := time.Now().Add(24 * time.Hour)

		err = s.postRepo.CreatePostCollaborators(
			txCtx,
			post.ID,
			in.CustomerID,
			invited,
			expiresAt,
		)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
