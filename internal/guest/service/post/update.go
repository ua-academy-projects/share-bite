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
			return entity.Post{}, fmt.Errorf("check venue existence via venue provider: %w", apperror.ErrUpstreamError)
		}
		if !exists {
			return entity.Post{}, apperror.VenueNotFoundID(*in.VenueID)
		}
	}

	post, err := s.postRepo.Update(ctx, in)
	if err != nil {
		return entity.Post{}, fmt.Errorf("update post in post repository: %w", err)
	}

	return post, nil
}

func isValidPostStatusTransition(from, to entity.PostStatus) bool {
	switch from {
	case entity.PostStatusDraft:
		return to == entity.PostStatusDraft ||
			to == entity.PostStatusPublished ||
			to == entity.PostStatusArchived ||
			to == entity.PostStatusDeleted
	case entity.PostStatusPublished:
		return to == entity.PostStatusPublished ||
			to == entity.PostStatusArchived ||
			to == entity.PostStatusDeleted
	case entity.PostStatusArchived:
		return to == entity.PostStatusArchived ||
			to == entity.PostStatusPublished ||
			to == entity.PostStatusDeleted
	default:
		return false
	}
}
