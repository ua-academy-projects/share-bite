package follow

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/internal/guest/service/customer"
	"time"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

type CustomerFollowRepository interface {
	Follow(ctx context.Context, followerID, followedID string) (entity.CustomerFollow, error)
	Unfollow(ctx context.Context, followerID, followedID string) error
	ListFollowing(ctx context.Context, customerID string, cursorTime time.Time, cursorID string, limit int) ([]entity.Follower, error)
	ListFollowers(ctx context.Context, customerID string, cursorTime time.Time, cursorID string, limit int) ([]entity.Follower, error)
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

type pageToken struct {
	CreatedAt time.Time `json:"created_at"`
	ID        string    `json:"id"`
}

func (s *service) isOwner(ctx context.Context, requesterUserID string, targetCustomerID string) bool {
	if requesterUserID == "" {
		return false
	}

	reqCustomer, err := s.customerRepo.GetByUserID(ctx, requesterUserID)
	if err != nil {
		return false
	}

	return reqCustomer.ID == targetCustomerID
}

func normalizeLimit(limit int) int {
	switch {
	case limit <= 0:
		return defaultLimit
	case limit > maxLimit:
		return maxLimit
	default:
		return limit
	}
}

func (s *service) parsePageToken(token string) (time.Time, string, error) {
	if token == "" {
		return time.Time{}, "", nil
	}
	data, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return time.Time{}, "", err
	}
	var pt pageToken
	if err := json.Unmarshal(data, &pt); err != nil {
		return time.Time{}, "", err
	}
	return pt.CreatedAt, pt.ID, nil
}

func (s *service) generatePageToken(createdAt time.Time, id string) string {
	pt := pageToken{
		CreatedAt: createdAt,
		ID:        id,
	}
	data, _ := json.Marshal(pt)
	return base64.StdEncoding.EncodeToString(data)
}

func followersToCustomers(rows []entity.Follower) []entity.Customer {
	customers := make([]entity.Customer, 0, len(rows))
	for _, f := range rows {
		customers = append(customers, f.Customer)
	}
	return customers
}
