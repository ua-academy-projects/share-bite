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

func (c *CustomerFollow) ToEntity() entity.CustomerFollow {
	return entity.CustomerFollow{
		ID:                 strconv.FormatInt(c.ID, 10),
		FollowerCustomerID: c.FollowerCustomerID,
		FollowedCustomerID: c.FollowedCustomerID,
		CreatedAt:          c.CreatedAt,
	}
}

type CustomerFollows []CustomerFollow

func (fs CustomerFollows) ToEntities() []entity.CustomerFollow {
	res := make([]entity.CustomerFollow, 0, len(fs))
	for i := range fs {
		res = append(res, fs[i].ToEntity())
	}
	return res
}

type FollowCustomerRow struct {
	FollowID        int64     `db:"follow_id"`
	FollowCreatedAt time.Time `db:"follow_created_at"`

	ID                string    `db:"id"`
	UserID            string    `db:"user_id"`
	UserName          string    `db:"username"`
	FirstName         string    `db:"first_name"`
	LastName          string    `db:"last_name"`
	AvatarObjectKey   *string   `db:"avatar_object_key"`
	Bio               *string   `db:"bio"`
	IsFollowersPublic bool      `db:"is_followers_public"`
	IsFollowingPublic bool      `db:"is_following_public"`
	CreatedAt         time.Time `db:"created_at"`
}

func (r *FollowCustomerRow) ToEntity() entity.Follower {
	return entity.Follower{
		Customer: entity.Customer{
			ID:                r.ID,
			UserID:            r.UserID,
			UserName:          r.UserName,
			FirstName:         r.FirstName,
			LastName:          r.LastName,
			AvatarObjectKey:   r.AvatarObjectKey,
			Bio:               r.Bio,
			IsFollowersPublic: r.IsFollowersPublic,
			IsFollowingPublic: r.IsFollowingPublic,
			CreatedAt:         r.CreatedAt,
		},
		FollowCreatedAt: r.FollowCreatedAt,
		FollowID:        strconv.FormatInt(r.FollowID, 10),
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
