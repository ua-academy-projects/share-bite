package post

import (
	"context"
	"fmt"
	"time"

	"github.com/ua-academy-projects/share-bite/pkg/logger"
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
		go func() {
			pubCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			authorUserID, err := s.postRepo.GetAuthorUserID(pubCtx, postID)
			if err == nil && authorUserID != "" {
				err = s.publisher.Publish(pubCtx, authorUserID, notification.Message{
					UserID:    authorUserID,
					Type:      notification.PostLiked,
					Data:      postID,
					CreatedAt: time.Now().UTC(),
				})
				if err != nil {
					logger.ErrorKV(pubCtx, "publish post liked notification failed", "error", err)
				}
			} else if err != nil {
				logger.ErrorKV(pubCtx, "failed to get author user ID for notification", "error", err)
			}
		}()
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
