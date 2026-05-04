package post

import (
	"context"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

func (s *service) AcceptInvitation(ctx context.Context, collaboratorID string, customerID string) error {

	return s.txManager.ReadCommitted(ctx, func(txCtx context.Context) error {

		postID, err := s.postRepo.AcceptPostInvitation(txCtx, collaboratorID, customerID)
		if err != nil {
			return err
		}

		allAccepted, err := s.postRepo.AreAllPostCollaboratorsAccepted(txCtx, postID)
		if err != nil {
			return err
		}

		if allAccepted {
			return s.postRepo.UpdateStatus(
				txCtx,
				postID,
				customerID,
				entity.PostStatusPublished,
			)
		}

		return nil
	})
}
