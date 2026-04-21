package business

import (
	"context"
	"errors"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	repository "github.com/ua-academy-projects/share-bite/internal/business/repository/business"
)

func (s *service) CreateLocation(ctx context.Context, brandID int, ownerUserID string, in dto.CreateLocationInput) (*entity.OrgUnit, error) {
	const op = "business.service.CreateLocation"

	ownerBrandID, err := s.businessRepo.GetBrandIDByOwnerUserID(ctx, ownerUserID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperror.Forbidden("business profile not found")
		}
		return nil, fmt.Errorf("%s: get owner brand id: %w", op, err)
	}

	if ownerBrandID != brandID {
		return nil, apperror.Forbidden("you can manage only your own brand locations")
	}

	brand, err := s.businessRepo.GetById(ctx, brandID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperror.OrgUnitNotFoundID(brandID)
		}
		return nil, fmt.Errorf("%s: get brand: %w", op, err)
	}
	if brand.ProfileType != entity.ProfileTypeBrand {
		return nil, apperror.BadRequest("target org unit is not a brand")
	}

	location, err := s.businessRepo.CreateLocation(ctx, brandID, ownerUserID, in)
	if err != nil {
		return nil, fmt.Errorf("%s: create location: %w", op, err)
	}

	return location, nil
}

func (s *service) UpdateLocation(ctx context.Context, locationID int, ownerUserID string, in dto.UpdateLocationInput) (*entity.OrgUnit, error) {
	const op = "business.service.UpdateLocation"

	ownerBrandID, err := s.businessRepo.GetBrandIDByOwnerUserID(ctx, ownerUserID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperror.Forbidden("business profile not found")
		}
		return nil, fmt.Errorf("%s: get owner brand id: %w", op, err)
	}

	location, err := s.businessRepo.GetById(ctx, locationID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperror.LocationNotFoundID(locationID)
		}
		return nil, fmt.Errorf("%s: get location: %w", op, err)
	}

	if location.ProfileType != entity.ProfileTypeVenue {
		return nil, apperror.BadRequest("target org unit is not a location")
	}

	if location.ParentId == nil || *location.ParentId != ownerBrandID {
		return nil, apperror.Forbidden("you can manage only your own brand locations")
	}

	updated, err := s.businessRepo.UpdateLocation(ctx, locationID, ownerBrandID, in)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperror.LocationNotFoundID(locationID)
		}
		return nil, fmt.Errorf("%s: update location: %w", op, err)
	}

	return updated, nil
}

func (s *service) DeleteLocation(ctx context.Context, locationID int, ownerUserID string) error {
	const op = "business.service.DeleteLocation"

	ownerBrandID, err := s.businessRepo.GetBrandIDByOwnerUserID(ctx, ownerUserID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return apperror.Forbidden("business profile not found")
		}
		return fmt.Errorf("%s: get owner brand id: %w", op, err)
	}

	location, err := s.businessRepo.GetById(ctx, locationID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return apperror.LocationNotFoundID(locationID)
		}
		return fmt.Errorf("%s: get location: %w", op, err)
	}

	if location.ProfileType != entity.ProfileTypeVenue {
		return apperror.BadRequest("target org unit is not a location")
	}

	if location.ParentId == nil || *location.ParentId != ownerBrandID {
		return apperror.Forbidden("you can manage only your own brand locations")
	}

	if err := s.businessRepo.DeleteLocation(ctx, locationID, ownerBrandID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return apperror.LocationNotFoundID(locationID)
		}
		return fmt.Errorf("%s: delete location: %w", op, err)
	}

	return nil
}
