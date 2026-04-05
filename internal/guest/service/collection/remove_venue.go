package collection

import (
	"context"
	"fmt"

	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

func (s *service) RemoveVenue(
	ctx context.Context,
	collectionID string,
	customerID string,
	venueID string,
) error {
	collection, err := s.collectionRepo.GetCollection(ctx, collectionID)
	if err != nil {
		return fmt.Errorf("get collection from repository: %w", err)
	}
	if collection.CustomerID != customerID {
		return apperror.ErrCollectionAccessDenied
	}

	exists, err := s.collectionRepo.CheckIfVenueInCollection(ctx, collectionID, venueID)
	if err != nil {
		return fmt.Errorf("check if venue is in collection from repository: %w", err)
	}
	if !exists {
		return apperror.VenueNotFoundInCollection(venueID)
	}

	if err := s.collectionRepo.RemoveVenue(ctx, collectionID, venueID); err != nil {
		return fmt.Errorf("remove venue from collection in repository: %w", err)
	}

	return nil
}
