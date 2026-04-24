package comment

import (
	"context"
	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
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
