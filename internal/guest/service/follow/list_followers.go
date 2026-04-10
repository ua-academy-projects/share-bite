package follow

import (
	"context"
	"fmt"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

func (s *service) ListFollowers(ctx context.Context, customerID string) ([]entity.Customer, error) {
	_, err := s.customerRepo.GetByID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	follow, err := s.customerFollowRepo.ListFollowers(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("list followers from repository: %w", err)
	}

	return follow, nil
}
