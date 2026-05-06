package collection

import (
	"context"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

func (s *service) DeclineInvitation(ctx context.Context, invitationID string, customerID string) error {
	invitation, err := s.collectionRepo.GetInvitation(ctx, invitationID)
	if err != nil {
		return fmt.Errorf("get invitation from repository: %w", err)
	}
	if invitation.InviteeID != customerID {
		return apperror.InvitationNotFoundID(invitationID)
	}

	if !invitation.CanBeDeclined() {
		return apperror.ErrInvitationAlreadyProcessed
	}

	if err := s.collectionRepo.UpdateInvitationStatus(ctx, invitationID, entity.DeclinedInvitationStatus); err != nil {
		return fmt.Errorf("update invitation status in repository: %w", err)
	}

	return nil
}
