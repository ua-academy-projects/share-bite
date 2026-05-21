package post

import (
	"context"
	"github.com/ua-academy-projects/share-bite/internal/imageprocessing"
	"github.com/ua-academy-projects/share-bite/internal/storage"
	"github.com/ua-academy-projects/share-bite/pkg/database"
	"github.com/ua-academy-projects/share-bite/pkg/outbox"
	"time"

	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

type postRepository interface {
	Create(ctx context.Context, in dto.CreatePostInput) (entity.Post, error)
	Update(ctx context.Context, in entity.UpdatePostInput) (entity.Post, error)
	List(ctx context.Context, in dto.ListPostsInput) (dto.ListPostsOutput, error)
	Get(ctx context.Context, postID string, reqCustomerID string) (entity.Post, error)
	GetByID(ctx context.Context, postID string) (entity.Post, error)
	GetAuthorUserID(ctx context.Context, postID string) (string, error)
	GetAuthorCustomerID(ctx context.Context, postID string) (string, error)
	DeleteImagesByPostIDReturningKeys(ctx context.Context, postID string) ([]string, error)
	Like(ctx context.Context, postID string, customerID string) (bool, error)
	Unlike(ctx context.Context, postID string, customerID string) error
	CreateImages(ctx context.Context, images []entity.PostImage) ([]entity.PostImage, error)
	DeleteImagesByPostID(ctx context.Context, postID string) error
	UpdateStatus(ctx context.Context, postID, customerID string, status entity.PostStatus) error
	GetPostsByVenueIDs(ctx context.Context, venueIDs []int64, limit int) ([]entity.Post, error)
	CreateMentions(ctx context.Context, mentions []entity.PostMention) error
	ListMentionsByPostIDs(ctx context.Context, postIDs []string) (map[string][]entity.PostMention, error)

	CreatePostCollaborators(ctx context.Context, postID string, invitedBy string, customerIDs []string, expiresAt time.Time) error
	GetPendingPostInvitations(ctx context.Context, customerID string) ([]entity.PostCollaborator, error)
	AcceptPostInvitation(ctx context.Context, collaboratorID string, customerID string) (string, error)
	DeclinePostInvitation(ctx context.Context, collaboratorID string, customerID string) error
	GetAcceptedPostCollaborators(ctx context.Context, postID string) ([]string, error)
	TryPublishPostIfAllAccepted(ctx context.Context, postID string) (bool, error)
	IsAcceptedCollaborator(ctx context.Context, postID string, customerID string) (bool, error)

	DeleteExpiredDraftPosts(ctx context.Context) error

	UpdateProcessedMetadata(ctx context.Context, imageID string, thumbnailKey string, width int, height int) error
	MarkProcessingFailed(ctx context.Context, imageID string, reason string) error
	IsAlreadyProcessed(ctx context.Context, imageID string) (bool, error)
	MarkProcessing(ctx context.Context, imageID string) error
}

type VenueProvider interface {
	CheckExists(ctx context.Context, venueID int64) (bool, error)
	GetNearbyVenues(ctx context.Context, lat, lon float64, limit int) ([]int64, error)
}

type followRepo interface {
	GetAllowedMentions(ctx context.Context, customerID string, mentions []string) ([]string, error)
}

type customerRepo interface {
	GetByIDs(ctx context.Context, ids []string) ([]entity.Customer, error)
	GetByID(ctx context.Context, id string) (entity.Customer, error)
}

type postCleanupService interface {
	CleanupExpiredPosts(ctx context.Context) error
}

type service struct {
	postRepo                postRepository
	venueProvider           VenueProvider
	followRepo              followRepo
	customerRepo            customerRepo
	storage                 storage.ObjectStorage
	txManager               database.TxManager
	outboxWriter            outbox.Writer
	imageProcessingProducer *imageprocessing.Producer
}

type Option func(*service)

func WithOutboxWriter(writer outbox.Writer) Option {
	return func(s *service) {
		s.outboxWriter = writer
	}
}

func WithImageProcessingProducer(producer *imageprocessing.Producer) Option {
	return func(s *service) {
		s.imageProcessingProducer = producer
	}
}

func New(postRepo postRepository, venueProvider VenueProvider, followRepo followRepo, customerRepo customerRepo, storage storage.ObjectStorage, txManager database.TxManager, opts ...Option) *service {
	if storage == nil {
		panic("post service: storage is not configured")
	}
	if txManager == nil {
		panic("post service: transaction manager is not configured")
	}
	svc := &service{
		postRepo:      postRepo,
		venueProvider: venueProvider,
		followRepo:    followRepo,
		customerRepo:  customerRepo,
		storage:       storage,
		txManager:     txManager,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(svc)
		}
	}

	return svc
}
