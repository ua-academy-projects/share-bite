package post

import (
	"context"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

func (s *service) Get(ctx context.Context, postID string) (entity.Post, error) {
	post, err := s.postRepo.Get(ctx, postID)
	if err != nil {
		return entity.Post{}, fmt.Errorf("get post from post repository: %w", err)
	}

	return post, nil
}
