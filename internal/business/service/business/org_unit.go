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
	orgUnit, err := s.businessRepo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperror.OrgUnitNotFoundID(id)
		}
		return nil, fmt.Errorf("%w", err)
	}

	return orgUnit, nil
}

func (s *service) List(ctx context.Context, brandId, skip, limit int) (pagination.Result[entity.OrgUnit], error) {
	result, err := s.businessRepo.ListByParentID(ctx, brandId, skip, limit)
	if err != nil {
		return pagination.Result[entity.OrgUnit]{}, fmt.Errorf("%w", err)
	}

	return result, nil
}
