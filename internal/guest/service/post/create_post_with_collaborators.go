package post

import (
	"context"
	"encoding/json"
	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/pkg/notification"
	"time"
)

const (
	postInvitationTTL          = 24 * time.Hour
	notificationPublishTimeout = 5 * time.Second
)

func (s *service) CreatePostWithCollaborators(ctx context.Context, in dto.CreatePostInput) (entity.Post, error) {
	var result entity.Post
	var invited []string

	err := s.txManager.ReadCommitted(ctx, func(txCtx context.Context) error {

		post, err := s.createPost(txCtx, in)
		if err != nil {
			return err
		}

		result = post

		invited = uniqueAndExcludeSelf(in.CustomerID, in.InvitedCustomerIDs)

		if len(invited) == 0 {

			err := s.postRepo.UpdateStatus(
				txCtx,
				post.ID,
				in.CustomerID,
				entity.PostStatusPublished,
			)
			if err != nil {
				return err
			}
			updatedPost, err := s.postRepo.GetByID(txCtx, post.ID)
			if err != nil {
				return err
			}
			result = updatedPost
			return nil

		}

		expiresAt := time.Now().Add(postInvitationTTL)

		return s.postRepo.CreatePostCollaborators(
			txCtx,
			post.ID,
			in.CustomerID,
			invited,
			expiresAt,
		)
	})

	if err == nil && s.publisher != nil && len(invited) > 0 {
		go func() {
			detached := context.WithoutCancel(ctx)
			publishCtx, cancel := context.WithTimeout(detached, notificationPublishTimeout)
			defer cancel()

			for _, customerID := range invited {
				customer, err := s.customerRepo.GetByID(publishCtx, customerID)
				if err != nil {
					continue
				}

				data := map[string]string{
					"post_id":    result.ID,
					"inviter_id": in.CustomerID,
				}

				dataBytes, _ := json.Marshal(data)

				msg := notification.Message{
					UserID:    customer.UserID,
					Type:      notification.InvitationReceived,
					Data:      string(dataBytes),
					CreatedAt: time.Now().UTC(),
				}

				_ = s.publisher.Publish(publishCtx, customer.UserID, msg)
			}
		}()
	}

	return result, err
}
