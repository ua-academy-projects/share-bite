package comment

import (
	"context"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"fmt"
)

func (s *service) Update(ctx context.Context, postID int64, in entity.UpdateCommentInput) (entity.Comment, error) {
	comment, err := s.commentRepo.GetByID(ctx, in.CommentID)
	if err != nil {
		return entity.Comment{}, fmt.Errorf("get comment by id: %w", err)
	}

	if comment.CustomerID != in.CustomerID {
		return entity.Comment{}, apperror.ErrForbidden
	}

	updatedComment, err := s.commentRepo.Update(ctx, postID, in)
	if err != nil {
		return entity.Comment{}, fmt.Errorf("update comment in repo: %w", err)
	}

	return updatedComment, nil
}

func (s *service) Delete(ctx context.Context, postID int64, commentID int64, customerID string) error {
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return fmt.Errorf("get comment by id: %w", err)
	}

	if comment.CustomerID != customerID {
		return apperror.ErrForbidden
	}

	if err := s.commentRepo.Delete(ctx, postID, commentID); err != nil {
		return fmt.Errorf("delete comment from repo: %w", err)
	}

	return s.commentRepo.Delete(ctx, commentID, postID)
}
