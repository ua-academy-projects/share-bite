package collection

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/guest/error/code"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	"github.com/ua-academy-projects/share-bite/pkg/notification"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (s *service) InviteCollaborator(ctx context.Context, in entity.InviteCollaboratorInput) error {
	var invitationID string
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
				timeLeft := resendInvitationCooldown - time.Since(invitation.LastSentAt)
				return apperror.InvitationCooldown(timeLeft)
			}

			if err := s.collectionRepo.RefreshInvitation(ctx, invitation.ID, time.Now().UTC().Add(invitationTTL)); err != nil {
				return fmt.Errorf("refresh invitation in repository: %w", err)
			}

			invitationID = invitation.ID
		}

		return nil
	}); txErr != nil {
		return txErr
	}

	// Push to Pub/Sub: triggers a real-time UI notification for the invitee if they are currently online.
	if s.publisher != nil {
		go func() {
			detachedCtx := context.WithoutCancel(ctx)
			publishCtx, cancel := context.WithTimeout(detachedCtx, time.Second*5)
			defer cancel()

			kvs := []zapcore.Field{
				zap.String("collection_id", in.CollectionID),
				zap.String("invitation_id", invitationID),
				zap.String("invitee_id", in.InviteeID),
			}
			publishCtx = logger.WithFields(publishCtx, kvs...)

			inviteeCustomer, err := s.customerRepo.GetByID(publishCtx, in.InviteeID)
			if err != nil {
				logger.WarnKV(publishCtx, "get invitee customer from repository failed", "error", err)
				return
			}
			inviteeUserID := inviteeCustomer.UserID

			logger.InfoKV(publishCtx, "publishing invitation event to redis..")

			msg := notification.NewMessage(
				notification.InvitationReceived,
				inviteeUserID,
				in.InviterID,
				"collection",
				in.CollectionID,
				time.Now().UTC(),
			)

			if err := s.publisher.Publish(publishCtx, inviteeUserID, msg); err != nil {
				logger.ErrorKV(publishCtx, "publishing invitation event to redis failed", "error", err)
				return
			}

			logger.InfoKV(publishCtx, "successfully published invitation event to redis")
		}()
	}

	// TODO: Push to Queue (e.g., SQS/RabbitMQ/Asynq)
	// Purpose: guaranteed background processing for offline delivery (e.g., sending an email to the invitee).

	return nil
}
