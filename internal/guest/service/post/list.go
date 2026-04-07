package post

import (
	"context"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

func (s *service) List(ctx context.Context, in entity.ListPostsInput) (entity.ListPostsOutput, error) {
	out, err := s.postRepo.List(ctx, in)
	if err != nil {
		return entity.ListPostsOutput{}, fmt.Errorf("get list of posts from post repository: %w", err)
	}

	return out, nil
}
