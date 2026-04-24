package follow

import (
	"context"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

func (s *service) Follow(ctx context.Context, userID, targetCustomerID string) (entity.CustomerFollow, error) {
	currentCustomer, err := s.customerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return entity.CustomerFollow{}, err
	}

	if currentCustomer.ID == targetCustomerID {
		return entity.CustomerFollow{}, apperror.ErrCannotFollowYourself
	}

	isFollowing, err := s.customerFollowRepo.IsFollowing(ctx, currentCustomer.ID, targetCustomerID)
	if err != nil {
		return entity.CustomerFollow{}, err
	}

	if isFollowing {
		return entity.CustomerFollow{}, apperror.AlreadyFollowing(currentCustomer.ID, targetCustomerID)
	}

	follow, err := s.customerFollowRepo.Follow(ctx, currentCustomer.ID, targetCustomerID)
	if err != nil {
		return entity.CustomerFollow{}, err
	}

	return follow, nil
}
