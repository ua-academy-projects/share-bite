package post

import (
	"context"
	"fmt"
)

func (s *service) Like(ctx context.Context, postID string, customerID string) error {
	if _, err := s.Get(ctx, postID, customerID); err != nil {
		return fmt.Errorf("validate post for like: %w", err)
	}

	if err := s.postRepo.Like(ctx, postID, customerID); err != nil {
		return fmt.Errorf("like post in repository: %w", err)
	}
	return nil
}

func (s *service) Unlike(ctx context.Context, postID string, customerID string) error {
	if _, err := s.Get(ctx, postID, customerID); err != nil {
		return fmt.Errorf("validate post for unlike: %w", err)
	}

	if err := s.postRepo.Unlike(ctx, postID, customerID); err != nil {
		return fmt.Errorf("unlike post in repository: %w", err)
	}
	return nil
}
