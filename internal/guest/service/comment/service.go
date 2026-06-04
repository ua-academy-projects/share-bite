package comment

import (
	"context"
	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/internal/storage"
	"github.com/ua-academy-projects/share-bite/pkg/database"
	"github.com/ua-academy-projects/share-bite/pkg/outbox"
)

type commentRepository interface {
	Create(ctx context.Context, in dto.CreateCommentInput) (entity.Comment, error)
	GetByID(ctx context.Context, commentID int64) (entity.Comment, error)
	Update(ctx context.Context, in dto.UpdateCommentInput) (entity.Comment, error)
	Delete(ctx context.Context, commentID int64) error
	List(ctx context.Context, in dto.ListCommentsInput) (dto.ListCommentsOutput, error)
}

type postService interface {
	Get(ctx context.Context, postID string, reqCustomerID string) (entity.Post, error)
}

type customerRepo interface {
	GetByID(ctx context.Context, id string) (entity.Customer, error)
}

type Option func(*service)

func WithOutboxWriter(writer outbox.Writer) Option {
	return func(s *service) {
		s.outboxWriter = writer
	}
}

func WithCustomerRepo(repo customerRepo) Option {
	return func(s *service) {
		s.customerRepo = repo
	}
}

func WithTxManager(txManager database.TxManager) Option {
	return func(s *service) {
		s.txManager = txManager
	}
}

func WithStorage(storage storage.ObjectStorage) Option {
	return func(s *service) {
		s.storage = storage
	}
}

type service struct {
	commentRepo  commentRepository
	postSvc      postService
	customerRepo customerRepo
	txManager    database.TxManager
	outboxWriter outbox.Writer
	storage      storage.ObjectStorage
}

func New(commentRepo commentRepository, postSvc postService, opts ...Option) *service {
	svc := &service{
		commentRepo: commentRepo,
		postSvc:     postSvc,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(svc)
		}
	}

	return svc
}
