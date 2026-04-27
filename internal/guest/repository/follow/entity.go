package follow

import (
	"fmt"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"strconv"
	"time"
)

type CustomerFollow struct {
	ID int64 `db:"id"`

	FollowerCustomerID string `db:"follower_customer_id"`
	FollowedCustomerID string `db:"followed_customer_id"`

	CreatedAt time.Time `db:"created_at"`
}

type FollowCustomerRow struct {
	FollowID        int64     `db:"follow_id"`
	FollowCreatedAt time.Time `db:"follow_created_at"`

	ID                string  `db:"id"`
	AvatarObjectKey   *string `db:"avatar_object_key"`
	IsFollowersPublic bool    `db:"is_followers_public"`
	IsFollowingPublic bool    `db:"is_following_public"`

	IsFollowing  bool `db:"is_following"`
	IsFollowedBy bool `db:"is_followed_by"`
	IsMutual     bool `db:"is_mutual"`
}

func (r *FollowCustomerRow) ToEntity() entity.Follower {
	return entity.Follower{
		Customer: entity.Customer{
			ID:                r.ID,
			AvatarObjectKey:   r.AvatarObjectKey,
			IsFollowersPublic: r.IsFollowersPublic,
			IsFollowingPublic: r.IsFollowingPublic,
		},
		FollowCreatedAt: r.FollowCreatedAt,
		FollowID:        strconv.FormatInt(r.FollowID, 10),
		IsFollowing:     r.IsFollowing,
		IsFollowedBy:    r.IsFollowedBy,
		IsMutual:        r.IsMutual,
	}
}

func executeSQLError(err error) error {
	return fmt.Errorf("execute sql: %w", err)
}

func scanRowError(err error) error {
	return fmt.Errorf("scan row: %w", err)
}

func scanRowsError(err error) error {
	return fmt.Errorf("scan rows: %w", err)
}
