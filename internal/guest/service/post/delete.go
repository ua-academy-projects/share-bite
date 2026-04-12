package post

import (
	"context"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

func (s *service) Delete(ctx context.Context, postID, customerID string) error {
	currentPost, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("get post by id in post repository: %w", err)
	}

	if currentPost.CustomerID != customerID {
		return apperror.PostNotFoundID(postID)
	}

	if currentPost.Status == entity.PostStatusDeleted {
		return nil
	}

	deletedStatus := entity.PostStatusDeleted
	_, err = s.postRepo.Update(ctx, entity.UpdatePostInput{
		ID:         postID,
		CustomerID: customerID,
		Status:     &deletedStatus,
	})
	if err != nil {
		return fmt.Errorf("delete post in post repository: %w", err)
	}

	return nil
}
