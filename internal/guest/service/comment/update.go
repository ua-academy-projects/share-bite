package comment

import (
	"context"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/pkg/errwrap"
)

func (s *service) Update(ctx context.Context, in entity.UpdateCommentInput) (entity.Comment, error) {
	comment, err := s.commentRepo.GetByID(ctx, in.CommentID)
	if err != nil {
		return entity.Comment{}, errwrap.Wrap("get comment by id", err)
	}

	if comment.CustomerID != in.CustomerID {
		return entity.Comment{}, apperror.ErrForbidden
	}

	updatedComment, err := s.commentRepo.Update(ctx, in)
	if err != nil {
		return entity.Comment{}, errwrap.Wrap("update comment in repo", err)
	}

	return updatedComment, nil
}

func (s *service) Delete(ctx context.Context, commentID int64, customerID string) error {
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return errwrap.Wrap("get comment by id", err)
	}

	if comment.CustomerID != customerID {
		return apperror.ErrForbidden
	}

	if err := s.commentRepo.Delete(ctx, commentID); err != nil {
		return errwrap.Wrap("delete comment from repo", err)
	}

	return nil
}
