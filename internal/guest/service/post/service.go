package post

import (
	"context"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

type postRepository interface {
	Create(ctx context.Context, in entity.CreatePostInput) (entity.Post, error)
	List(ctx context.Context, in entity.ListPostsInput) (entity.ListPostsOutput, error)
	Get(ctx context.Context, postID string, reqCustomerID string) (entity.Post, error)
	Like(ctx context.Context, postID string, customerID string) error
	Unlike(ctx context.Context, postID string, customerID string) error
}

type VenueProvider interface {
	CheckExists(ctx context.Context, venueID string) (bool, error)
}

type VenueProvider interface {
	CheckExists(ctx context.Context, venueID string) (bool, error)
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
