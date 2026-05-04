package collection

import (
	"context"
	"fmt"

	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

func (s *service) AddVenue(
	ctx context.Context,
	collectionID string,
	customerID string,
	venueID int64,
) error {
	// TODO: check whether this venue exists

	if txErr := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		collection, err := s.collectionRepo.GetCollectionForUpdate(ctx, collectionID)
		if err != nil {
			return fmt.Errorf("get collection from repository: %w", err)
		}
		if err := s.requireCollaborator(ctx, collectionID, customerID, collection.CustomerID); err != nil {
			return err
		}

		count, err := s.collectionRepo.CountVenues(ctx, collectionID)
		if err != nil {
			return fmt.Errorf("get count of collection venues from repository: %w", err)
		}
		if count >= maxVenuesPerCollection {
			return apperror.CollectionVenuesLimitReached(maxVenuesPerCollection)
		}

		maxSortOrder, err := s.collectionRepo.GetMaxSortOrder(ctx, collectionID)
		if err != nil {
			return fmt.Errorf("get max sort order from repository: %w", err)
		}
		sortOrder := maxSortOrder + sortOrderGap

		if err := s.collectionRepo.AddVenue(ctx, collectionID, venueID, sortOrder); err != nil {
			return fmt.Errorf("add venue to collection in repository: %w", err)
		}

		return nil
	}); txErr != nil {
		return txErr
	}

	return nil
}
