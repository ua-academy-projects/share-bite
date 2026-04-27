package post

import (
	"context"
	"github.com/ua-academy-projects/share-bite/internal/storage"
	"github.com/ua-academy-projects/share-bite/pkg/database"

	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

type postRepository interface {
	Create(ctx context.Context, in dto.CreatePostInput) (entity.Post, error)
	Update(ctx context.Context, in entity.UpdatePostInput) (entity.Post, error)
	List(ctx context.Context, in dto.ListPostsInput) (dto.ListPostsOutput, error)
	Get(ctx context.Context, postID string, reqCustomerID string) (entity.Post, error)
	GetByID(ctx context.Context, postID string) (entity.Post, error)
	Like(ctx context.Context, postID string, customerID string) error
	Unlike(ctx context.Context, postID string, customerID string) error
	CreateImages(ctx context.Context, images []entity.PostImage) error
	DeleteImagesByPostID(ctx context.Context, postID string) error
	UpdateStatus(ctx context.Context, postID, customerID string, status entity.PostStatus) error
	GetPostsByVenueIDs(ctx context.Context, venueIDs []int64, limit int) ([]entity.Post, error)
}

type VenueProvider interface {
	CheckExists(ctx context.Context, venueID int64) (bool, error)
	GetNearbyVenues(ctx context.Context, lat, lon float64, limit int) ([]int64, error)
}

type service struct {
	postRepo      postRepository
	venueProvider VenueProvider
	storage       storage.ObjectStorage
	txManager     database.TxManager
}

func New(postRepo postRepository, venueProvider VenueProvider, storage storage.ObjectStorage, txManager database.TxManager) *service {
	if storage == nil {
		panic("post service: storage is not configured")
	}
	if txManager == nil {
		panic("post service: transaction manager is not configured")
	}
	return &service{
		postRepo:      postRepo,
		venueProvider: venueProvider,
		storage:       storage,
		txManager:     txManager,
	}
}
