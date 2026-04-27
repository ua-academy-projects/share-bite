package post

import (
	"context"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
)

func (s *service) ExploreNearby(ctx context.Context, lat, lon float64, limit int) ([]dto.ExploreVenueItem, error) {
	venueIDs, err := s.venueProvider.GetNearbyVenues(ctx, lat, lon, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get nearby venues: %w", err)
	}

	if len(venueIDs) == 0 {
		return []dto.ExploreVenueItem{}, nil
	}

	posts, err := s.postRepo.GetPostsByVenueIDs(ctx, venueIDs, limit*3)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch posts: %w", err)
	}

	postsByVenue := make(map[int64][]dto.PostItem)
	for _, p := range posts {

		imageURLs := make([]string, 0, len(p.Images))

		for _, img := range p.Images {
			imageURLs = append(imageURLs, img.ObjectKey)
		}

		postsByVenue[p.VenueID] = append(postsByVenue[p.VenueID], dto.PostItem{
			ID:        p.ID,
			Content:   p.Text,
			CreatedAt: p.CreatedAt,
			Images:    imageURLs,
		})
	}

	result := make([]dto.ExploreVenueItem, 0, len(venueIDs))
	for _, vid := range venueIDs {
		venuePosts := postsByVenue[vid]

		if venuePosts == nil {
			venuePosts = make([]dto.PostItem, 0)
		}

		result = append(result, dto.ExploreVenueItem{
			VenueID: vid,
			Posts:   venuePosts,
		})
	}

	return result, nil
}
