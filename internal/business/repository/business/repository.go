package business

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"



	"github.com/ua-academy-projects/share-bite/pkg/database"
	"github.com/ua-academy-projects/share-bite/pkg/database/pg"

	biserr "github.com/ua-academy-projects/share-bite/internal/business/error"
)

var (
	ErrNotFound  = errors.New("not found")
	ErrForbidden = errors.New("forbidden")
)

var (
	ErrNotFound  = errors.New("not found")
	ErrForbidden = errors.New("forbidden")
)

type Repository struct {
	db database.Client
}

func New(db database.Client) *Repository {
	return &Repository{
		db: db,
	}
}