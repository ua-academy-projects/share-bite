package follow

import (
	"context"
	"fmt"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

func (s *service) ListFollowing(ctx context.Context, customerID string) ([]entity.Customer, error) {
	_, err := s.customerRepo.GetByID(ctx, customerID)
	if err != nil {
		return nil, err
	}
	follow, err := s.customerFollowRepo.ListFollowing(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("list following from repository: %w", err)
	}

	return follow, nil
}
