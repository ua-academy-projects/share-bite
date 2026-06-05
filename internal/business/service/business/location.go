package business

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ua-academy-projects/share-bite/pkg/database/pagination"

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
		if in.Latitude != nil && in.Longitude != nil {
			hash := s.h3Service.GetH3Index(float64(*in.Latitude), float64(*in.Longitude), s.h3Config.Resolution)
			in.H3Hash = &hash
		}

		location, err = s.businessRepo.CreateLocation(ctxTx, brandID, ownerUserID, in)
		if err != nil {
			return fmt.Errorf("%s: create location: %w", "business.service.CreateLocation", err)
		}

		if err := s.businessRepo.SetOrgUnitTagsByIDs(ctxTx, location.Id, in.TagIDs); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return apperror.BadRequest("one or more tags are invalid")
			}
			return fmt.Errorf("%s: set location tags: %w", op, err)
		}

		tags, err := s.businessRepo.GetOrgUnitTagSlugs(ctxTx, location.Id)
		if err != nil {
			return fmt.Errorf("%s: get location tags: %w", op, err)
		}
		location.Tags = tags

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

		var h3Hash *string
		if in.Latitude != nil && in.Longitude != nil {
			hash := s.h3Service.GetH3Index(float64(*in.Latitude), float64(*in.Longitude), s.h3Config.Resolution)
			h3Hash = &hash
		}

		updated, err = s.businessRepo.UpdateLocation(ctxTx, locationID, ownerBrandID, in, h3Hash)
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

func (s *service) ListNearbyVenues(ctx context.Context, lat, lon float64, skip, limit int) (pagination.Result[entity.OrgUnitWithDistance], error) {
	return s.businessRepo.ListNearbyVenues(ctx, lat, lon, skip, limit)
}

func (s *service) UpdateVenueHours(
	ctx context.Context,
	locationID int,
	ownerUserID string,
	in dto.UpdateVenueHoursInput,
) (*dto.UpdateVenueHoursOutput, error) {
	const op = "business.service.UpdateVenueHours"

	if len(in.Days) == 0 {
		return nil, apperror.BadRequest("days is required")
	}

	seen := make(map[int]struct{}, 7)
	for _, d := range in.Days {
		if _, ok := seen[d.Weekday]; ok {
			return nil, apperror.BadRequest("duplicate weekday")
		}
		seen[d.Weekday] = struct{}{}

		if d.OpenTime == nil && d.CloseTime == nil {
			continue
		}

		if d.OpenTime == nil || d.CloseTime == nil {
			return nil, apperror.BadRequest("both openTime and closeTime must be provided together")
		}

		openT, err := time.Parse("15:04", *d.OpenTime)
		if err != nil {
			return nil, apperror.BadRequest("openTime must be HH:MM")
		}

		closeT, err := time.Parse("15:04", *d.CloseTime)
		if err != nil {
			return nil, apperror.BadRequest("closeTime must be HH:MM")
		}

		if !openT.Before(closeT) {
			return nil, apperror.BadRequest("openTime must be before closeTime")
		}
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

	var days []dto.VenueHoursDayInput

	err = s.txManager.ReadCommitted(ctx, func(ctxTx context.Context) error {
		if err := s.businessRepo.ReplaceLocationHours(ctxTx, locationID, in.Days); err != nil {
			return fmt.Errorf("%s: replace location hours: %w", op, err)
		}

		var err error
		days, err = s.businessRepo.GetLocationHours(ctxTx, locationID)
		if err != nil {
			return fmt.Errorf("%s: get location hours: %w", op, err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &dto.UpdateVenueHoursOutput{
		VenueID: locationID,
		Days:    days,
	}, nil
}
