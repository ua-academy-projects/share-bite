package business

import (
	"context"
	"errors"

	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	repository "github.com/ua-academy-projects/share-bite/internal/business/repository/business"
	"github.com/ua-academy-projects/share-bite/pkg/database/pagination"
	"github.com/ua-academy-projects/share-bite/pkg/errwrap"
)

func (s *service) Get(ctx context.Context, id int) (*entity.OrgUnit, error) {
	orgUnit, err := s.businessRepo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperror.OrgUnitNotFoundID(id)
		}
		return nil, errwrap.Wrap("get org unit from business repository", err)
	}

	return orgUnit, nil
}

func (s *service) List(ctx context.Context, brandId, skip, limit int) (pagination.Result[entity.OrgUnit], error) {
	result, err := s.businessRepo.ListByParentID(ctx, brandId, skip, limit)
	if err != nil {
		return pagination.Result[entity.OrgUnit]{}, errwrap.Wrap("list locations from business repository", err)
	}

	return result, nil
}
