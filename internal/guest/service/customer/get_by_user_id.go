package customer

import (
	"context"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/pkg/errwrap"
)

func (s *service) GetByUserID(ctx context.Context, userID string) (entity.Customer, error) {
	customer, err := s.customerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return entity.Customer{}, errwrap.Wrap("get customer by user id from repository", err)
	}

	return customer, nil
}
