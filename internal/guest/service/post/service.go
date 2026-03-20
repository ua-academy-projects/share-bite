package post

import (
	"context"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

type postRepository interface {
	List(ctx context.Context, in entity.ListPostsInput) (entity.ListPostsOutput, error)
	Get(ctx context.Context, postID string) (entity.Post, error)
}

type service struct {
	postRepo postRepository
}

func New(postRepo postRepository) *service {
	return &service{
		postRepo: postRepo,
	}
}
