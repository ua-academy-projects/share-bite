package business

import (
	"context"
	"errors"

	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	repository "github.com/ua-academy-projects/share-bite/internal/business/repository/business"
	"github.com/ua-academy-projects/share-bite/pkg/errwrap"
)

func (s *service) CreateLocation(ctx context.Context, brandID int, ownerUserID string, in dto.CreateLocationInput) (*entity.OrgUnit, error) {
	ownerBrandID, err := s.businessRepo.GetBrandIDByOwnerUserID(ctx, ownerUserID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperror.Forbidden("business profile not found")
		}
		return nil, errwrap.Wrap("get owner brand id", err)
	}

	if ownerBrandID != brandID {
		return nil, apperror.Forbidden("you can manage only your own brand locations")
	}

	brand, err := s.businessRepo.GetById(ctx, brandID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperror.OrgUnitNotFoundID(brandID)
		}
		return nil, errwrap.Wrap("get brand", err)
	}
	if brand.ProfileType != "BRAND" {
		return nil, apperror.BadRequest("target org unit is not a brand")
	}

	location, err := s.businessRepo.CreateLocation(ctx, brandID, ownerUserID, in)
	if err != nil {
		return nil, errwrap.Wrap("create location", err)
	}

	return location, nil
}

func (s *service) UpdateLocation(ctx context.Context, locationID int, ownerUserID string, in dto.UpdateLocationInput) (*entity.OrgUnit, error) {
	ownerBrandID, err := s.businessRepo.GetBrandIDByOwnerUserID(ctx, ownerUserID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperror.Forbidden("business profile not found")
		}
		return nil, errwrap.Wrap("get owner brand id", err)
	}

	location, err := s.businessRepo.GetById(ctx, locationID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperror.LocationNotFoundID(locationID)
		}
		return nil, errwrap.Wrap("get location", err)
	}

	if location.ProfileType != "VENUE" {
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
		return nil, errwrap.Wrap("update location", err)
	}

	return updated, nil
}

func (s *service) DeleteLocation(ctx context.Context, locationID int, ownerUserID string) error {
	ownerBrandID, err := s.businessRepo.GetBrandIDByOwnerUserID(ctx, ownerUserID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return apperror.Forbidden("business profile not found")
		}
		return errwrap.Wrap("get owner brand id", err)
	}

	location, err := s.businessRepo.GetById(ctx, locationID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return apperror.LocationNotFoundID(locationID)
		}
		return errwrap.Wrap("get location", err)
	}

	if location.ProfileType != "VENUE" {
		return apperror.BadRequest("target org unit is not a location")
	}

	if location.ParentId == nil || *location.ParentId != ownerBrandID {
		return apperror.Forbidden("you can manage only your own brand locations")
	}

	if err := s.businessRepo.DeleteLocation(ctx, locationID, ownerBrandID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return apperror.LocationNotFoundID(locationID)
		}
		return errwrap.Wrap("delete location", err)
	}

	return nil
}
