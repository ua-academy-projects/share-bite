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

	return orgUnit, nil
}

func (s *service) List(ctx context.Context, brandId, skip, limit int) (pagination.Result[entity.OrgUnit], error) {
	const op = "service.business.List"

	result, err := s.businessRepo.ListByParentID(ctx, brandId, skip, limit)
	if err != nil {
		return pagination.Result[entity.OrgUnit]{}, fmt.Errorf("%s: %w", op, err)
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
