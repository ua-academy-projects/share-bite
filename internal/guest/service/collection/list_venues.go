package collection

import (
	"context"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

func (s *service) ListVenues(
	ctx context.Context,
	collectionID string,
	customerID *string,
) ([]entity.EnrichedVenueItem, error) {
	collection, err := s.collectionRepo.GetCollection(ctx, collectionID)
	if err != nil {
		return nil, fmt.Errorf("get collection from repository: %w", err)
	}
	if !canAccessCollection(collection, customerID) {
		return nil, apperror.CollectionNotFoundID(collectionID)
	}

	collectionVenues, err := s.collectionRepo.ListCollectionVenues(ctx, collectionID)
	if err != nil {
		return nil, fmt.Errorf("get collection venues from repository: %w", err)
	}
	if len(collectionVenues) == 0 {
		return nil, nil
	}

	venueIDs := make([]string, 0, len(collectionVenues))
	for _, v := range collectionVenues {
		venueIDs = append(venueIDs, v.VenueID)
	}

	venues, err := s.businessClient.ListVenues(ctx, venueIDs)
	if err != nil {
		return nil, fmt.Errorf("get venues from business service: %w", err)
	}

	enrichedVenues := make([]entity.EnrichedVenueItem, 0, len(collectionVenues))
	for _, cv := range collectionVenues {
		venue, ok := venues[cv.VenueID]
		if !ok {
			// The venue has been deleted in another service so
			// we're just not showing it.
			continue
		}

		enrichedVenue := entity.EnrichedVenueItem{
			VenueItem: entity.Venue{
				ID:          venue.ID,
				Name:        venue.Name,
				Description: venue.Description,
				AvatarURL:   venue.AvatarURL,
			},
			SortOrder: cv.SortOrder,
			AddedAt:   cv.AddedAt,
		}
		enrichedVenues = append(enrichedVenues, enrichedVenue)
	}

	return enrichedVenues, nil
}
