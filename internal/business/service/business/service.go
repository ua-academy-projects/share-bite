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

type businessRepository interface {
	GetById(ctx context.Context, id int) (*entity.OrgUnit, error)
	ListByParentID(ctx context.Context, parentID, offset, limit int) (pagination.Result[entity.OrgUnit], error)
	GetVenuesByIDs(ctx context.Context, ids []int) ([]entity.OrgUnit, error)
}

type service struct {
	businessRepo businessRepository
}

func New(businessRepo businessRepository) *service {
	return &service{
		businessRepo: businessRepo,
	}
}

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

func (s *service) GetVenuesByIDs(ctx context.Context, ids []int) ([]entity.OrgUnit, error) {
	venues, err := s.businessRepo.GetVenuesByIDs(ctx, ids)
	if err != nil {
		return nil, errwrap.Wrap("get venues by ids from business repository", err)
	}

	return venues, nil
}
