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
	venueID string,
) error {
	collection, err := s.collectionRepo.GetCollection(ctx, collectionID)
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
}
