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

func (s *service) List(ctx context.Context, brandId, skip, limit int) (pagination.Result[entity.OrgUnit], error) {
	const op = "service.business.List"

	result, err := s.businessRepo.ListByParentID(ctx, brandId, skip, limit)
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
