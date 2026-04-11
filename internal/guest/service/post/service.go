package post

import (
	"context"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

type postRepository interface {
	Create(ctx context.Context, in entity.CreatePostInput) (entity.Post, error)
	Update(ctx context.Context, in entity.UpdatePostInput) (entity.Post, error)
	List(ctx context.Context, in entity.ListPostsInput) (entity.ListPostsOutput, error)
	Get(ctx context.Context, postID string) (entity.Post, error)
	GetByID(ctx context.Context, postID string) (entity.Post, error)
}

type VenueProvider interface {
	CheckExists(ctx context.Context, venueID int64) (bool, error)
}

type service struct {
	postRepo      postRepository
	venueProvider VenueProvider
}

func New(postRepo postRepository, venueProvider VenueProvider) *service {
	return &service{
		postRepo:      postRepo,
		venueProvider: venueProvider,
	}
}
