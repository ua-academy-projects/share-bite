package customer

import (
	"context"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/pkg/errwrap"
)

func (s *service) GetByUserName(ctx context.Context, userName string) (entity.Customer, error) {
	customer, err := s.customerRepo.GetByUserName(ctx, userName)
	if err != nil {
		return entity.Customer{}, errwrap.Wrap("get customer by user name from repo", err)
	}

	return customer, nil
}
