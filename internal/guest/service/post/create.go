package post

import (
	"context"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/pkg/errwrap"
)

func (s *service) Create(ctx context.Context, in entity.CreatePostInput) (entity.Post, error) {
	post, err := s.postRepo.Create(ctx, in)
	if err != nil {
		return entity.Post{}, errwrap.Wrap("create post in post repository", err)
	}

	return post, nil
}
