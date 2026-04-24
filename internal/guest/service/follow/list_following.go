package follow

import (
	"context"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

func (s *service) ListFollowing(ctx context.Context, in entity.ListFollowingInput) (entity.ListFollowingOutput, error) {
	cursorTime, cursorID, err := s.parsePageToken(in.PageToken)
	if err != nil {
		return entity.ListFollowingOutput{}, apperror.ErrInvalidPageToken
	}

	limit := normalizeLimit(in.PageSize)

	targetCustomer, err := s.customerRepo.GetByID(ctx, in.TargetCustomerID)
	if err != nil {
		return entity.ListFollowingOutput{}, err
	}

	isOwner := in.RequesterCustomerID != nil && *in.RequesterCustomerID == targetCustomer.ID
	if !isOwner && !targetCustomer.IsFollowersPublic {
		return entity.ListFollowingOutput{}, apperror.ErrFollowersListPrivate
	}

	rows, err := s.customerFollowRepo.ListFollowing(
		ctx,
		in.TargetCustomerID,
		cursorTime,
		cursorID,
		limit+1,
	)
	if err != nil {
		return entity.ListFollowingOutput{}, err
	}

	var nextPageToken string

	if len(rows) > limit {
		rows = rows[:limit]

		last := rows[len(rows)-1]
		nextPageToken = s.generatePageToken(
			last.FollowCreatedAt,
			last.FollowID,
		)
	}

	return entity.ListFollowingOutput{
		Customers:     followersToCustomers(rows),
		NextPageToken: nextPageToken,
	}, nil
}
