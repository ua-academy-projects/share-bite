package business

import (
	"context"
	"errors"
	"fmt"

	repository "github.com/ua-academy-projects/share-bite/internal/business/repository/business"
)

func (s *service) GetOnboardingContext(ctx context.Context, userID string) (brandID int, venueID int, err error) {
	const op = "service.business.GetOnboardingContext"

	brandID, err = s.businessRepo.GetBrandIDByOwnerUserID(ctx, userID)
	if err != nil {
		if !errors.Is(err, repository.ErrNotFound) {
			return 0, 0, fmt.Errorf("%s: get brand id: %w", op, err)
		}
		brandID = 0
	}

	if brandID == 0 {
		return 0, 0, nil
	}

	venueID, err = s.businessRepo.GetFirstVenueIDByOwnerUserID(ctx, userID)
	if err != nil {
		if !errors.Is(err, repository.ErrNotFound) {
			return brandID, 0, fmt.Errorf("%s: get venue id: %w", op, err)
		}
		venueID = 0
	}

	return brandID, venueID, nil
}
