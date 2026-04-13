package customer

import (
	"context"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

func (s *service) Create(ctx context.Context, in entity.CreateCustomer) (string, error) {
	customerID, err := s.customerRepo.Create(ctx, in)
	if err != nil {
		return "", fmt.Errorf("create customer in repo: %w", err)
	}

	return customerID, nil
}
