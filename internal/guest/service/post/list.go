package post

import (
	"context"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
)

func (s *service) List(ctx context.Context, in dto.ListPostsInput) (dto.ListPostsOutput, error) {
	out, err := s.postRepo.List(ctx, in)
	if err != nil {
		return dto.ListPostsOutput{}, fmt.Errorf("get list of posts from post repository: %w", err)
	}

	return out, nil
}
