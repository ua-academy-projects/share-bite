package customer

import (
	"context"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

func (s *service) GetByUserID(ctx context.Context, userID string) (entity.Customer, error) {
	customer, err := s.customerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return entity.Customer{}, fmt.Errorf("get customer by user id from repository: %w", err)
	}

	return customer, nil
}
