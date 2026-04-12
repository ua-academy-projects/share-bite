package follow

import (
	"context"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

func (s *service) ListFollowing(ctx context.Context, targetCustomerID, requesterUserID string) ([]entity.Customer, error) {
	targetCustomer, err := s.customerRepo.GetByUserID(ctx, targetCustomerID)
	if err != nil {
		return nil, err
	}

	isOwner := false
	if requesterUserID != "" {
		requesterCustomer, err := s.customerRepo.GetByUserID(ctx, requesterUserID)
		if err == nil && requesterCustomer.ID == targetCustomer.ID {
			isOwner = true
		}
	}

	if !isOwner && !targetCustomer.IsFollowingPublic {
		return nil, apperror.ErrFollowingListPrivate
	}

	return s.customerFollowRepo.ListFollowing(ctx, targetCustomer.ID)
}
