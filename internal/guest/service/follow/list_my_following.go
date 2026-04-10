package follow

import (
	"context"
	"fmt"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

func (s *service) ListMyFollowing(ctx context.Context, userID string) ([]entity.Customer, error) {
	currentCustomer, err := s.customerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	follow, err := s.customerFollowRepo.ListFollowing(ctx, currentCustomer.ID)
	if err != nil {
		return nil, fmt.Errorf("list following from repository: %w", err)
	}

	return follow, nil
}
