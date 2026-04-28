package collection

import (
	"context"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

func (s *service) ListInvitations(ctx context.Context, in entity.ListInvitationsInput) (entity.ListInvitationsOutput, error) {
	// check if the caller has sufficient permissions.
	// only the owner can view the full list of invitations for a collection.
	if in.CollectionID != nil {
		collection, err := s.collectionRepo.GetCollection(ctx, *in.CollectionID)
		if err != nil {
			return entity.ListInvitationsOutput{}, fmt.Errorf("get collection from repository: %w", err)
		}

		if err := s.requireOwner(ctx, *in.CollectionID, in.CallerID, collection.CustomerID); err != nil {
			return entity.ListInvitationsOutput{}, err
		}
	}

	invitations, err := s.collectionRepo.ListInvitations(ctx, in)
	if err != nil {
		return entity.ListInvitationsOutput{}, fmt.Errorf("get list of invitations from repository: %w", err)
	}

	var nextCursorID string
	requestLimit := in.Limit - 1

	if len(invitations) > requestLimit {
		invitations = invitations[:requestLimit]

		lastItem := invitations[len(invitations)-1]
		nextCursorID = lastItem.ID
	}

	return entity.ListInvitationsOutput{
		Invitations:  invitations,
		NextCursorID: nextCursorID,
	}, nil
}
