package comment

import (
	"context"
	"github.com/ua-academy-projects/share-bite/internal/guest/dto"

	"fmt"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

func (s *service) Update(ctx context.Context, postID int64, in dto.UpdateCommentInput) (entity.Comment, error) {
	comment, err := s.commentRepo.GetByID(ctx, in.CommentID)
	if err != nil {
		return entity.Comment{}, fmt.Errorf("get comment for update: %w", err)
	}

	if comment.PostID != postID {
		return entity.Comment{}, apperror.CommentNotFoundID(in.CommentID)
	}

	if comment.CustomerID != in.CustomerID {
		return entity.Comment{}, apperror.ErrForbidden
	}

	updatedComment, err := s.commentRepo.Update(ctx, in)
	if err != nil {
		return entity.Comment{}, fmt.Errorf("update comment in repo: %w", err)
	}

	return updatedComment, nil
}
