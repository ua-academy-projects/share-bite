package comment

import (
	"context"
	"fmt"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

func (s *service) Delete(ctx context.Context, postID int64, commentID int64, customerID string) error {
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return fmt.Errorf("get comment to delete: %w", err)
	}

	if comment.PostID != postID {
		return apperror.CommentNotFoundID(commentID)
	}

	if comment.CustomerID != customerID {
		return apperror.ErrForbidden
	}

	if err := s.commentRepo.Delete(ctx, commentID); err != nil {
		return fmt.Errorf("delete comment from repo: %w", err)
	}

	return nil
}
