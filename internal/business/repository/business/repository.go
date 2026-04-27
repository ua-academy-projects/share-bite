package business

import (
	"errors"

	"github.com/ua-academy-projects/share-bite/pkg/database"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrForbidden        = errors.New("forbidden")
	ErrNoAvailableItems = errors.New("no available box items")
)

type Repository struct {
	db database.Client
}

func New(db database.Client) *Repository {
	return &Repository{
		db: db,
	}
}
