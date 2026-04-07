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
	return fmt.Errorf("execute sql: %w", err)
}

func scanRowError(err error) error {
	return fmt.Errorf("scan row: %w", err)
}

func scanRowsError(err error) error {
	return fmt.Errorf("scan rows: %w", err)
}
