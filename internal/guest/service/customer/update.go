package customer

import (
	"context"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

func (s *service) Update(ctx context.Context, in entity.UpdateCustomer) (entity.Customer, error) {
	_, err := s.customerRepo.GetByUserID(ctx, in.UserID)
	if err != nil {
		return entity.Customer{}, fmt.Errorf("get customer by user id from repository: %w", err)
	}

	customer, err := s.customerRepo.Update(ctx, in)
	if err != nil {
		return entity.Customer{}, fmt.Errorf("update customer in repo: %w", err)
	}

	return customer, nil
}