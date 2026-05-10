package post

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	"github.com/ua-academy-projects/share-bite/pkg/notification"
	"time"
)

func (s *service) AcceptInvitation(ctx context.Context, collaboratorID string, customerID string) error {
	var postID string
	var allAccepted bool

	err := s.txManager.ReadCommitted(ctx, func(txCtx context.Context) error {
		var err error

		postID, err = s.postRepo.AcceptPostInvitation(txCtx, collaboratorID, customerID)
		if err != nil {
			return fmt.Errorf("accept invitation: %w", err)
		}

		allAccepted, err = s.postRepo.TryPublishPostIfAllAccepted(txCtx, postID)
		if err != nil {
			return fmt.Errorf("try publish post: %w", err)
		}

		return nil
	})

	if err != nil || s.publisher == nil {
		return err
	}

	go func() {
		detached := context.WithoutCancel(ctx)
		publishCtx, cancel := context.WithTimeout(detached, notificationPublishTimeout)
		defer cancel()

		// TODO: notify author that user accepted invitation

		if allAccepted {
			collaborators, err := s.postRepo.GetAcceptedPostCollaborators(publishCtx, postID)
			if err != nil {
				return
			}

			authorID, err := s.postRepo.GetAuthorCustomerID(publishCtx, postID)
			if err == nil {
				collaborators = append(collaborators, authorID)
			}

			seen := make(map[string]struct{})

			for _, cid := range collaborators {
				if _, ok := seen[cid]; ok {
					continue
				}
				seen[cid] = struct{}{}

				customer, err := s.customerRepo.GetByID(publishCtx, cid)
				if err != nil {
					continue
				}

				data := map[string]string{
					"post_id": postID,
				}

				dataBytes, err := json.Marshal(data)
				if err != nil {
					logger.ErrorKV(
						publishCtx,
						"marshal post published notification failed",
						"post_id",
						postID,
						"error",
						err,
					)
					continue
				}

				msg := notification.Message{
					UserID:    customer.UserID,
					Type:      notification.PostPublished,
					Data:      string(dataBytes),
					CreatedAt: time.Now().UTC(),
				}

				if err := s.publisher.Publish(
					publishCtx,
					customer.UserID,
					msg,
				); err != nil {
					logger.ErrorKV(
						publishCtx,
						"publish post published notification failed",
						"user_id",
						customer.UserID,
						"post_id",
						postID,
						"error",
						err,
					)
				}
			}
		}
	}()

	return nil
}
