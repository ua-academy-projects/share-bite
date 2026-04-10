package follow

import (
	"context"
	"fmt"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

func (s *service) Follow(ctx context.Context, userID, targetCustomerID string) (entity.CustomerFollow, error) {
	currentCustomer, err := s.customerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return entity.CustomerFollow{}, fmt.Errorf("get current customer by user id: %w", err)
	}

	targetCustomer, err := s.customerRepo.GetByUserID(ctx, targetCustomerID)
	if err != nil {
		return entity.CustomerFollow{}, fmt.Errorf("get target customer by user id: %w", err)
	}

	if currentCustomer.ID == targetCustomer.ID {
		return entity.CustomerFollow{}, apperror.ErrCannotFollowYourself
	}

	isFollowing, err := s.customerFollowRepo.IsFollowing(ctx, currentCustomer.ID, targetCustomer.ID)
	if err != nil {
		return entity.CustomerFollow{}, fmt.Errorf("check follow relation: %w", err)
	}
	if isFollowing {
		return entity.CustomerFollow{}, apperror.AlreadyFollowing(currentCustomer.ID, targetCustomer.ID)
	}

	follow, err := s.customerFollowRepo.Follow(ctx, currentCustomer.ID, targetCustomer.ID)
	if err != nil {
		return entity.CustomerFollow{}, fmt.Errorf("create follow relation: %w", err)
	}

	return follow, nil
}
