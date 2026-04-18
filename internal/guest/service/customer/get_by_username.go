package customer

import (
	"context"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

func (s *service) GetByUserName(ctx context.Context, userName string) (entity.Customer, error) {
	customer, err := s.customerRepo.GetByUserName(ctx, userName)
	if err != nil {
		return entity.Customer{}, fmt.Errorf("failed to get customer by username: %w", err)
	}

	return customer, nil
}
