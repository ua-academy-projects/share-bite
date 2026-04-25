package collection

import (
	"context"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

func (s *service) AcceptInvitation(ctx context.Context, invitationID string, customerID string) error {
	baseInvitation, err := s.collectionRepo.GetInvitation(ctx, invitationID)
	if err != nil {
		return fmt.Errorf("get invitation from repository: %w", err)
	}
	if baseInvitation.InviteeID != customerID {
		return apperror.InvitationNotFoundID(invitationID)
	}

	if txErr := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		collection, err := s.collectionRepo.GetCollectionForUpdate(ctx, baseInvitation.CollectionID)
		if err != nil {
			return fmt.Errorf("get collection from repository: %w", err)
		}

		count, err := s.collectionRepo.CountCollaborators(ctx, collection.ID)
		if err != nil {
			return fmt.Errorf("get collaborators count from repository: %w", err)
		}
		if count >= maxCollaboratorsPerCollection {
			return apperror.CollectionCollaboratorsLimitReached(maxCollaboratorsPerCollection)
		}

		invitation, err := s.collectionRepo.GetInvitationForUpdate(ctx, invitationID)
		if err != nil {
			return fmt.Errorf("get invitation for update from repository: %w", err)
		}
		if !invitation.CanBeAccepted() {
			if invitation.IsExpired() {
				return apperror.ErrInvitationExpired
			}

			return apperror.ErrInvitationAlreadyProcessed
		}

		if err := s.collectionRepo.CreateCollaborator(ctx, invitation.CollectionID, invitation.InviteeID); err != nil {
			return fmt.Errorf("create collaborator in repository: %w", err)
		}

		if err := s.collectionRepo.UpdateInvitationStatus(ctx, invitationID, entity.AcceptedInvitationStatus); err != nil {
			return fmt.Errorf("update invitation status in repository: %w", err)
		}

		return nil
	}); txErr != nil {
		return txErr
	}

	return nil
}
