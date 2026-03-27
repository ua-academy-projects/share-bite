package customer

import (
	"context"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/pkg/errwrap"
)

func (s *service) Update(ctx context.Context, in entity.UpdateCustomer) (entity.Customer, error) {
	customer, err := s.customerRepo.Update(ctx, in)
	if err != nil {
		return entity.Customer{}, errwrap.Wrap("update customer in repo", err)
	}

	return customer, nil
}
