package follow

import (
	"context"
	"fmt"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

func (s *service) Unfollow(ctx context.Context, userID, targetCustomerID string) error {
	currentCustomer, err := s.customerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get current customer by user id: %w", err)
	}

	targetCustomer, err := s.customerRepo.GetByUserID(ctx, targetCustomerID)
	if err != nil {
		return fmt.Errorf("get target customer by user id: %w", err)
	}

	if currentCustomer.ID == targetCustomer.ID {
		return apperror.ErrCannotUnfollowYourself
	}

	if err := s.customerFollowRepo.Unfollow(ctx, currentCustomer.ID, targetCustomer.ID); err != nil {
		return fmt.Errorf("delete follow relation: %w", err)
	}

	return nil
}
