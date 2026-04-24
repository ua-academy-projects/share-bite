package follow

import (
	"context"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

func (s *service) Unfollow(
	ctx context.Context,
	userID, targetCustomerID string,
) error {

	currentCustomer, err := s.customerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if currentCustomer.ID == targetCustomerID {
		return apperror.ErrCannotUnfollowYourself
	}

	isFollowing, err := s.customerFollowRepo.IsFollowing(ctx, currentCustomer.ID, targetCustomerID)
	if err != nil {
		return err
	}

	if !isFollowing {
		return apperror.ErrFollowNotFound
	}

	if err := s.customerFollowRepo.Unfollow(ctx, currentCustomer.ID, targetCustomerID); err != nil {
		return err
	}

	return nil
}
