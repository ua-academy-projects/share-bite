package comment

import (
	"context"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

type commentRepository interface {
	Create(ctx context.Context, in entity.CreateCommentInput) (entity.Comment, error)
	GetByID(ctx context.Context, commentID int64) (entity.Comment, error)
	Update(ctx context.Context, postID int64, in entity.UpdateCommentInput) (entity.Comment, error)
	Delete(ctx context.Context, commentID int64, postID int64) error
	List(ctx context.Context, in entity.ListCommentsInput) (entity.ListCommentsOutput, error)
}

type postService interface {
	Get(ctx context.Context, postID string) (entity.Post, error)
}

type service struct {
	commentRepo commentRepository
	postSvc     postService
}

func New(commentRepo commentRepository, postSvc postService) *service {
	return &service{
		commentRepo: commentRepo,
		postSvc:     postSvc,
	}
}
