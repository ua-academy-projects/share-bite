package post

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/pkg/notification"
	"time"
)

func (s *service) AcceptInvitation(ctx context.Context, collaboratorID string, customerID string) error {
	var postID string
	var authorID string
	var allAccepted bool

	err := s.txManager.ReadCommitted(ctx, func(txCtx context.Context) error {
		var err error

		postID, authorID, err = s.postRepo.AcceptPostInvitation(txCtx, collaboratorID, customerID)
		if err != nil {
			return fmt.Errorf("accept invitation: %w", err)
		}

		allAccepted, err = s.postRepo.AreAllPostCollaboratorsAccepted(txCtx, postID)
		if err != nil {
			return fmt.Errorf("check collaborators accepted: %w", err)
		}

		if allAccepted {
			err := s.postRepo.UpdateStatus(txCtx, postID, authorID, entity.PostStatusPublished)
			if err != nil {
				return fmt.Errorf("update post status: %w", err)
			}
		}

		return nil
	})

	if err != nil || s.publisher == nil {
		return err
	}

	go func() {
		detached := context.WithoutCancel(ctx)
		publishCtx, cancel := context.WithTimeout(detached, 5*time.Second)
		defer cancel()

		// TODO: notify author that user accepted invitation

		if allAccepted {
			collaborators, err := s.postRepo.GetAcceptedPostCollaborators(publishCtx, postID)
			if err != nil {
				return
			}

			authorID, err := s.postRepo.GetAuthorUserID(publishCtx, postID)
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

				dataBytes, _ := json.Marshal(data)

				msg := notification.Message{
					UserID:    customer.UserID,
					Type:      notification.PostPublished,
					Data:      string(dataBytes),
					CreatedAt: time.Now().UTC(),
				}

				_ = s.publisher.Publish(publishCtx, customer.UserID, msg)
			}
		}
	}()

	return nil
}
