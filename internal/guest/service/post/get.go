package post

import (
	"context"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/pkg/errwrap"
)

func (s *service) Get(ctx context.Context, postID string, reqCustomerID string) (entity.Post, error) {
	post, err := s.postRepo.Get(ctx, postID, reqCustomerID)
	if err != nil {
		return entity.Post{}, errwrap.Wrap("get post from post repository", err)
	}

	return post, nil
}
