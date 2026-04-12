package customer

import (
	"fmt"
	"time"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

type Customer struct {
	ID     string `db:"id"`
	UserID string `db:"user_id"`

	UserName  string `db:"username"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`

	AvatarObjectKey *string `db:"avatar_object_key"`
	Bio             *string `db:"bio"`

	IsFollowersPublic bool `db:"is_followers_public"`
	IsFollowingPublic bool `db:"is_following_public"`

	CreatedAt time.Time `db:"created_at"`
}

func (e Customer) ToEntity() entity.Customer {
	return entity.Customer{
		ID:     e.ID,
		UserID: e.UserID,

		UserName:  e.UserName,
		FirstName: e.FirstName,
		LastName:  e.LastName,

		AvatarObjectKey: e.AvatarObjectKey,
		Bio:             e.Bio,

		IsFollowersPublic: e.IsFollowersPublic,
		IsFollowingPublic: e.IsFollowingPublic,

		CreatedAt: e.CreatedAt,
	}
}

type Customers []Customer

func (es Customers) ToEntities() []entity.Customer {
	res := make([]entity.Customer, 0, len(es))
	for i := range es {
		res = append(res, es[i].ToEntity())
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
