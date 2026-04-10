package follow

import (
	"context"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/internal/guest/service/customer"
)

type CustomerFollowRepository interface {
	Follow(ctx context.Context, followerID, followedID string) (entity.CustomerFollow, error)
	Unfollow(ctx context.Context, followerID, followedID string) error
	ListFollowing(ctx context.Context, customerID string) ([]entity.Customer, error)
	ListFollowers(ctx context.Context, customerID string) ([]entity.Customer, error)
	IsFollowing(ctx context.Context, followerID, followedID string) (bool, error)
}

type service struct {
	customerFollowRepo CustomerFollowRepository
	customerRepo       customer.CustomerRepository
}

func New(
	customerFollowRepo CustomerFollowRepository,
	customerRepo customer.CustomerRepository,
) *service {
	return &service{
		customerFollowRepo: customerFollowRepo,
		customerRepo:       customerRepo,
	}
}
