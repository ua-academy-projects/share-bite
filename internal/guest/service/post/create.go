package post

import (
	"context"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/pkg/errwrap"
)

func (s *service) Create(ctx context.Context, in entity.CreatePostInput) (entity.Post, error) {
	exists, err := s.venueProvider.CheckExists(ctx, in.VenueID)
	if err != nil {
		return entity.Post{}, errwrap.Wrap("check venue existence via venue provider", apperror.ErrUpstreamError)
	}
	if !exists {
		return entity.Post{}, apperror.VenueNotFoundID(in.VenueID)
	}

	post, err := s.postRepo.Create(ctx, in)
	if err != nil {
		return entity.Post{}, errwrap.Wrap("create post in post repository", err)
	}

	return post, nil
}
