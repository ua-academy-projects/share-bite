package collection

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/guest/error/code"
)

func (s *service) InviteCollaborator(ctx context.Context, in entity.InviteCollaboratorInput) error {
	if txErr := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		collection, err := s.collectionRepo.GetCollectionForUpdate(ctx, in.CollectionID)
		if err != nil {
			return fmt.Errorf("get collection from repository: %w", err)
		}
		if err := s.requireOwner(ctx, in.CollectionID, in.InviterID, collection.CustomerID); err != nil {
			return err
		}
		if in.InviteeID == collection.CustomerID {
			return apperror.CustomerAlreadyCollaborator(in.InviteeID)
		}

		isCollaborator, err := s.collectionRepo.CheckIfCollaborator(ctx, in.CollectionID, in.InviteeID)
		if err != nil {
			return fmt.Errorf("check if customer is already a collaborator: %w", err)
		}
		if isCollaborator {
			return apperror.CustomerAlreadyCollaborator(in.InviteeID)
		}

		count, err := s.collectionRepo.CountCollaborators(ctx, in.CollectionID)
		if err != nil {
			return fmt.Errorf("get count of collection collaborators from repository: %w", err)
		}
		if count >= maxCollaboratorsPerCollection {
			return apperror.CollectionCollaboratorsLimitReached(maxCollaboratorsPerCollection)
		}

		var invitationID string

		invitation, err := s.collectionRepo.GetInvitationByInvitee(ctx, in.CollectionID, in.InviteeID)
		if err != nil {
			// create new invitation
			var appErr *apperror.Error
			if errors.As(err, &appErr) && appErr.Code == code.NotFound {
				in.Expiry = time.Now().UTC().Add(invitationTTL)

				invitationID, err = s.collectionRepo.CreateInvitation(ctx, in)
				if err != nil {
					return fmt.Errorf("create invitation in repository: %w", err)
				}
			} else {
				return fmt.Errorf("get invitation by invitee from repository: %w", err)
			}
		} else {
			// refresh already sent invitation
			if !invitation.CanResend(resendInvitationCooldown) {
				return apperror.InvitationCooldown(resendInvitationCooldown)
			}

			if err := s.collectionRepo.RefreshInvitation(ctx, invitation.ID, time.Now().UTC().Add(invitationTTL)); err != nil {
				return fmt.Errorf("refresh invitation in repository: %w", err)
			}

			invitationID = invitation.ID
		}

		// TODO: push an event into the queue
		// then this event will be processed and
		// information about invitation will be sent to invitee customer
		_ = invitationID

		return nil
	}); txErr != nil {
		return txErr
	}

	return nil
}
