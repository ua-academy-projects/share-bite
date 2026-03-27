package business

import (
	"context"

	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/pkg/errwrap"
)

type businessRepository interface {
	GetById(ctx context.Context, id int) (*entity.OrgUnit, error)
	List(ctx context.Context, offset, limit int) ([]entity.OrgUnit, error)
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
		return nil, errwrap.Wrap("get org unit from business repository", err)
	}

	return orgUnit, nil
}

func (s *service) List(ctx context.Context, page, limit int) ([]entity.OrgUnit, error) {
	offset := (page - 1) * limit

	orgUnits, err := s.businessRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, errwrap.Wrap("list org units from business repository", err)
	}

	return orgUnits, nil
}
