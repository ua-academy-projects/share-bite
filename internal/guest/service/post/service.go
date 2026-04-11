package post

import (
	"context"
	"github.com/ua-academy-projects/share-bite/internal/storage"
	"github.com/ua-academy-projects/share-bite/pkg/database"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

type postRepository interface {
	Create(ctx context.Context, in entity.CreatePostInput) (entity.Post, error)
	List(ctx context.Context, in entity.ListPostsInput) (entity.ListPostsOutput, error)
	Get(ctx context.Context, postID string) (entity.Post, error)

	CreateImages(ctx context.Context, images []entity.PostImage) error
}

type VenueProvider interface {
	CheckExists(ctx context.Context, venueID string) (bool, error)
}

type service struct {
	postRepo      postRepository
	venueProvider VenueProvider
	storage       storage.ObjectStorage
	txManager     database.TxManager
}

func New(postRepo postRepository, venueProvider VenueProvider, storage storage.ObjectStorage, txManager database.TxManager) *service {
	return &service{
		postRepo:      postRepo,
		venueProvider: venueProvider,
		storage:       storage,
		txManager:     txManager,
	}
}
