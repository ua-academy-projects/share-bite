package customer

import (
	"context"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/pkg/errwrap"
)

func (s *service) Create(ctx context.Context, in entity.CreateCustomer) (string, error) {
	customerID, err := s.customerRepo.Create(ctx, in)
	if err != nil {
		return "", errwrap.Wrap("create customer in repo", err)
	}

	return customerID, nil
}
