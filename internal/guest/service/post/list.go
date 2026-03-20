package post

import (
	"context"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/pkg/errwrap"
)

func (s *service) List(ctx context.Context, in entity.ListPostsInput) (entity.ListPostsOutput, error) {
	out, err := s.postRepo.List(ctx, in)
	if err != nil {
		return entity.ListPostsOutput{}, errwrap.Wrap("get list of posts from post repository", err)
	}

	return out, nil
}
