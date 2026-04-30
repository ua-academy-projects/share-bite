package follow

import (
	"context"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

func (s *service) Follow(ctx context.Context, userID, targetCustomerID string) error {
	currentCustomer, err := s.customerRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if currentCustomer.ID == targetCustomerID {
		return apperror.ErrCannotFollowYourself
	}

	return s.customerFollowRepo.Follow(ctx, currentCustomer.ID, targetCustomerID)
}
