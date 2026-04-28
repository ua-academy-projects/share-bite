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

const maxLocationTags = 5

func (s *service) CreateLocation(ctx context.Context, brandID int, ownerUserID string, in dto.CreateLocationInput) (*entity.OrgUnit, error) {
	const op = "business.service.CreateLocation"

	if len(in.TagIDs) > maxLocationTags {
		return nil, apperror.BadRequest("location can have at most 5 tags")
	}

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

	var location *entity.OrgUnit

	err = s.txManager.ReadCommitted(ctx, func(ctxTx context.Context) error {
		created, err := s.businessRepo.CreateLocation(ctxTx, brandID, ownerUserID, in)
		if err != nil {
			return fmt.Errorf("%s: create location: %w", op, err)
		}

		if err := s.businessRepo.SetOrgUnitTagsByIDs(ctxTx, created.Id, in.TagIDs); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return apperror.BadRequest("one or more tags are invalid")
			}
			return fmt.Errorf("%s: set location tags: %w", op, err)
		}

		tags, err := s.businessRepo.GetOrgUnitTagSlugs(ctxTx, created.Id)
		if err != nil {
			return fmt.Errorf("%s: get location tags: %w", op, err)
		}
		created.Tags = tags

		location = created
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return location, nil
}

func (s *service) UpdateLocation(ctx context.Context, locationID int, ownerUserID string, in dto.UpdateLocationInput) (*entity.OrgUnit, error) {
	const op = "business.service.UpdateLocation"

	if in.TagIDs != nil && len(*in.TagIDs) > maxLocationTags {
		return nil, apperror.BadRequest("location can have at most 5 tags")
	}

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

	var updated *entity.OrgUnit

	err = s.txManager.ReadCommitted(ctx, func(ctxTx context.Context) error {
		var err error

		updated, err = s.businessRepo.UpdateLocation(ctxTx, locationID, ownerBrandID, in)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return apperror.LocationNotFoundID(locationID)
			}
			return fmt.Errorf("%s: update location: %w", op, err)
		}

		if in.TagIDs != nil {
			if err := s.businessRepo.SetOrgUnitTagsByIDs(ctxTx, locationID, *in.TagIDs); err != nil {
				if errors.Is(err, repository.ErrNotFound) {
					return apperror.BadRequest("one or more tags are invalid")
				}
				return fmt.Errorf("%s: set location tags: %w", op, err)
			}
		}

		tags, err := s.businessRepo.GetOrgUnitTagSlugs(ctxTx, locationID)
		if err != nil {
			return fmt.Errorf("%s: get location tags: %w", op, err)
		}
		updated.Tags = tags

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
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
