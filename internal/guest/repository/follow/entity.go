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

func executeSQLError(err error) error {
	return fmt.Errorf("execute sql: %w", err)
}

func scanRowError(err error) error {
	return fmt.Errorf("scan row: %w", err)
}

func scanRowsError(err error) error {
	return fmt.Errorf("scan rows: %w", err)
}
