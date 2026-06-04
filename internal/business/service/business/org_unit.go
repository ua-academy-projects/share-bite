package business

import (
	"context"
	"errors"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	repository "github.com/ua-academy-projects/share-bite/internal/business/repository/business"
	"github.com/ua-academy-projects/share-bite/pkg/database/pagination"
)

func (s *service) Get(ctx context.Context, id int) (*entity.OrgUnit, error) {
	const op = "service.business.Get"

	orgUnit, err := s.businessRepo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w", op, apperror.OrgUnitNotFoundID(id))
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	tags, err := s.businessRepo.GetOrgUnitTagSlugs(ctx, orgUnit.Id)
	if err != nil {
		return nil, fmt.Errorf("%s: get org unit tags: %w", op, err)
	}
	orgUnit.Tags = tags

	return orgUnit, nil
}

func (s *service) List(ctx context.Context, brandId, skip, limit int, tags []string) (pagination.Result[entity.OrgUnit], error) {
	const op = "service.business.List"

	result, err := s.businessRepo.ListByParentID(ctx, brandId, skip, limit, tags)
	if err != nil {
		return pagination.Result[entity.OrgUnit]{}, fmt.Errorf("%s: %w", op, err)
	}

	if len(result.Items) == 0 {
		return result, nil
	}

	ids := make([]int, 0, len(result.Items))
	for _, item := range result.Items {
		ids = append(ids, item.Id)
	}

	tagsByOrgUnitID, err := s.businessRepo.GetOrgUnitTagsByOrgUnitID(ctx, ids)
	if err != nil {
		return pagination.Result[entity.OrgUnit]{}, fmt.Errorf("%s: get org unit tags: %w", op, err)
	}

	for i := range result.Items {
		result.Items[i].Tags = tagsByOrgUnitID[result.Items[i].Id]
	}

	return result, nil
}

func (s *service) GetVenuesByIDs(ctx context.Context, ids []int) ([]entity.OrgUnit, error) {
	const op = "service.business.List"

	venues, err := s.businessRepo.GetVenuesByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return venues, nil
}

func (s *service) Rating(ctx context.Context, id int) (float32, error) {
	const op = "service.business.Rating"

	_, err := s.businessRepo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return 0, fmt.Errorf("%s: %w", op, apperror.OrgUnitNotFoundID(id))
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	rating, err := s.businessRepo.GetVenueRating(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return rating, nil
}

func (s *service) ListLocationTags(ctx context.Context) ([]entity.LocationTag, error) {
	const op = "service.business.ListLocationTags"

	tags, err := s.businessRepo.ListLocationTags(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return tags, nil
}

func (s *service) SearchVenues(ctx context.Context, query string, skip, limit int, tags []string) (pagination.Result[entity.OrgUnit], error) {
	const op = "service.business.SearchVenues"

	result, err := s.businessRepo.SearchVenues(ctx, query, skip, limit, tags)
	if err != nil {
		return pagination.Result[entity.OrgUnit]{}, fmt.Errorf("%s: %w", op, err)
	}

	if len(result.Items) == 0 {
		return result, nil
	}

	ids := make([]int, 0, len(result.Items))
	for _, item := range result.Items {
		ids = append(ids, item.Id)
	}

	tagsByOrgUnitID, err := s.businessRepo.GetOrgUnitTagsByOrgUnitID(ctx, ids)
	if err != nil {
		return pagination.Result[entity.OrgUnit]{}, fmt.Errorf("%s: get org unit tags: %w", op, err)
	}

	for i := range result.Items {
		result.Items[i].Tags = tagsByOrgUnitID[result.Items[i].Id]
	}

	return result, nil
}

func (s *service) ResubmitVerification(ctx context.Context, id int, userID string) error {
	const op = "service.business.ResubmitVerification"

	err := s.businessRepo.ResubmitVerification(ctx, id, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return fmt.Errorf("%s: %w", op, apperror.OrgUnitNotFoundID(id))
		}
		if errors.Is(err, repository.ErrInvalidStatus) {
			return fmt.Errorf("%s: %w", op, apperror.Conflict("organization unit is not in a resubmittable state (must be rejected)"))
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
