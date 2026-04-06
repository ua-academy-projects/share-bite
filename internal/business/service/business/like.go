package business

import (
	"context"
	"errors"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/internal/business/repository/business"
)

func (s *service) ToggleLike(ctx context.Context, postID int64, customerID string) (bool, error) {
	_, err := s.businessRepo.GetPostByID(ctx, postID)
	if err != nil {
		if errors.Is(err, business.ErrNotFound) {
			return false, apperror.PostNotFound(postID)
		}
		return false, fmt.Errorf("get post: %w", err)
	}

	liked, err := s.businessRepo.CheckUserLiked(ctx, postID, customerID)
	if err != nil {
		return false, fmt.Errorf("check user liked: %w", err)
	}

	if liked {
		err = s.businessRepo.DeleteLike(ctx, postID, customerID)
		if err != nil {
			return false, fmt.Errorf("delete like: %w", err)
		}
		return false, nil
	}

	_, err = s.businessRepo.CreateLike(ctx, postID, customerID)
	if err != nil {
		return false, fmt.Errorf("create like: %w", err)
	}
	return true, nil
}

func (s *service) GetLikes(ctx context.Context, postID int64, limit, offset int) ([]entity.Like, error) {
	likes, err := s.businessRepo.GetLikesByPost(ctx, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get likes: %w", err)
	}

	return likes, nil
}

func (s *service) GetLikeCount(ctx context.Context, postID int64) (int, error) {
	count, err := s.businessRepo.CountLikesByPost(ctx, postID)
	if err != nil {
		return 0, fmt.Errorf("count likes: %w", err)
	}

	return count, nil
}
