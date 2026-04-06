package post

import (
	"context"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

func (s *service) Create(ctx context.Context, in entity.CreatePostInput) (entity.Post, error) {
	exists, err := s.venueProvider.CheckExists(ctx, in.VenueID)
	if err != nil {
		return entity.Post{}, fmt.Errorf("check venue existence via venue provider: %w", apperror.ErrUpstreamError)
	}
	if !exists {
		return entity.Post{}, apperror.VenueNotFoundID(in.VenueID)
	}

	post, err := s.postRepo.Create(ctx, in)
	if err != nil {
		return entity.Post{}, fmt.Errorf("create post in post repository: %w", err)
	}

	return post, nil
}
