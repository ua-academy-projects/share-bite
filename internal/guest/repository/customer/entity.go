package customer

import (
	"time"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/pkg/errwrap"
)

type Customer struct {
	ID     string `db:"id"`
	UserID string `db:"user_id"`

	UserName  string `db:"username"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`

	AvatarObjectKey *string `db:"avatar_object_key"`
	Bio             *string `db:"bio"`

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

		CreatedAt: e.CreatedAt,
	}
}

func executeSQLError(err error) error {
	return errwrap.Wrap("execute sql", err)
}

func scanRowError(err error) error {
	return errwrap.Wrap("scan row", err)
}

func scanRowsError(err error) error {
	return errwrap.Wrap("scan rows", err)
}
