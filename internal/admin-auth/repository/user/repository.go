package user

import (
	"context"
	"fmt"

	"github.com/ua-academy-projects/share-bite/pkg/database"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
}

type Repository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
}

type repository struct {
	client database.Client
}

func New(client database.Client) Repository {
	return &repository{client: client}
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	q := database.Query{
		Name: "user.FindByEmail",
		Sql:  `SELECT id, email, password_hash FROM auth.users WHERE email = $1`,
	}

	row := r.client.DB().QueryRowContext(ctx, q, email)

	u := new(User)
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash); err != nil {
		return nil, fmt.Errorf("scan user: %w", err)
	}

	return u, nil
}
