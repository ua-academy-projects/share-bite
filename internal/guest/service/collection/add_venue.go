package collection

import (
	"context"
	"fmt"

	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

const (
	maxVenuesPerCollection = 100
	sortOrderGap           = 100.0
)

func (s *service) AddVenue(
	ctx context.Context,
	collectionID string,
	customerID string,
	venueID int64,
) error {
	// TODO: check whether this venue exists
	// exists, err := s.businessClient.CheckExists(ctx, venueID)
	// if err != nil {
	// 	return fmt.Errorf("check venue existence: %w", err)
	// }
	// if !exists {
	// 	// TODO: replace it with return apperror.VenueNotFoundID(venueID)
	// 	// now this will break another thing
	// 	return apperror.VenueNotFoundID(fmt.Sprintf("%d", venueID))
	// }

	if txErr := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		collection, err := s.collectionRepo.GetCollectionForUpdate(ctx, collectionID)
		if err != nil {
			return fmt.Errorf("get collection from repository: %w", err)
		}
		if collection.CustomerID != customerID {
			return apperror.ErrCollectionAccessDenied
		}

		count, err := s.collectionRepo.CountVenues(ctx, collectionID)
		if err != nil {
			return fmt.Errorf("get count of collection venues from repository: %w", err)
		}
		if count >= maxVenuesPerCollection {
			return apperror.ErrCollectionFull
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
