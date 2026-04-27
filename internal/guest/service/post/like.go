package post

import (
	"context"
	"fmt"
	"time"

	"github.com/ua-academy-projects/share-bite/pkg/notification"
)

func (s *service) Like(ctx context.Context, postID string, customerID string) error {
	post, err := s.Get(ctx, postID, customerID)
	if err != nil {
		return fmt.Errorf("validate post for like: %w", err)
	}

	if err := s.postRepo.Like(ctx, postID, customerID); err != nil {
		return fmt.Errorf("like post in repository: %w", err)
	}

	if s.publisher != nil && post.CustomerID != "" && post.CustomerID != customerID {
		authorUserID, err := s.postRepo.GetAuthorUserID(ctx, postID)
		if err == nil && authorUserID != "" {
			err = s.publisher.Publish(ctx, authorUserID, notification.Message{
				UserID:    authorUserID,
				Type:      notification.PostLiked,
				Data:      postID,
				CreatedAt: time.Now().UTC(),
			})
			if err != nil {
				return fmt.Errorf("publish post liked notification: %w", err)
			}
		}
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
