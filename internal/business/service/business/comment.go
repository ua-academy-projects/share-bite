package business

import (
	"context"
	"errors"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/internal/business/repository/business"
)

func (s *service) CreateComment(ctx context.Context, postID int64, authorID, content string) (*entity.Comment, error) {
	_, err := s.businessRepo.GetPostByID(ctx, postID)
	if err != nil {
		if errors.Is(err, business.ErrNotFound) {
			return nil, apperror.PostNotFound(postID)
		}
		return nil, fmt.Errorf("get post: %w", err)
	}

	comment, err := s.businessRepo.CreateComment(ctx, postID, authorID, content)
	if err != nil {
		return nil, fmt.Errorf("create comment: %w", err)
	}

	return comment, nil
}

func (s *service) UpdateComment(ctx context.Context, commentID int64, authorID, content string) (*entity.Comment, error) {
	comment, err := s.businessRepo.GetCommentByID(ctx, commentID)
	if err != nil {
		if errors.Is(err, business.ErrNotFound) {
			return nil, apperror.CommentNotFound(commentID)
		}
		return nil, fmt.Errorf("get comment: %w", err)
	}

	if comment.AuthorID != authorID {
		return nil, apperror.Forbidden("you can only edit your own comments")
	}

	updatedComment, err := s.businessRepo.UpdateComment(ctx, commentID, content)
	if err != nil {
		return nil, fmt.Errorf("update comment: %w", err)
	}

	return updatedComment, nil
}

func (s *service) DeleteComment(ctx context.Context, commentID int64, authorID string) error {
	comment, err := s.businessRepo.GetCommentByID(ctx, commentID)
	if err != nil {
		if errors.Is(err, business.ErrNotFound) {
			return apperror.CommentNotFound(commentID)
		}
		return fmt.Errorf("get comment: %w", err)
	}

	if comment.AuthorID != authorID {
		return apperror.Forbidden("you can only delete your own comments")
	}

	err = s.businessRepo.DeleteComment(ctx, commentID)
	if err != nil {
		return fmt.Errorf("delete comment: %w", err)
	}

	return nil
}

func (s *service) GetComments(ctx context.Context, postID int64, limit, offset int) ([]entity.CommentWithAuthor, error) {
	comments, err := s.businessRepo.ListCommentsWithAuthorsByPost(ctx, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list comments: %w", err)
	}

	return comments, nil
}

func (s *service) GetCommentCount(ctx context.Context, postID int64) (int, error) {
	count, err := s.businessRepo.CountCommentsByPost(ctx, postID)
	if err != nil {
		return 0, fmt.Errorf("count comments: %w", err)
	}

	return count, nil
}
