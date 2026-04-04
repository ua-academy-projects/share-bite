package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ua-academy-projects/share-bite/pkg/database"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
}

type Role struct {
	ID   int
	Slug string
	Name string
}

type CreateUser struct {
	Email        string
	PasswordHash string
}

type CreatedUser struct {
	ID    string
	Email string
}



type Repository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindRoleBySlug(ctx context.Context, slug string) (*Role, error)
	CreateUser(ctx context.Context, user CreateUser) (*CreatedUser, error)
	AssignRole(ctx context.Context, userID string, roleID int) error
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("scan user: %w", err)
	}

	return u, nil
}


func (r *repository) FindRoleBySlug(ctx context.Context, slug string) (*Role, error) {
	q := database.Query{
		Name: "user.FindRoleBySlug",
		Sql:  `SELECT id, slug, name FROM auth.roles WHERE slug = $1`,
	}

	row := r.client.DB().QueryRowContext(ctx, q, slug)

	role := new(Role)
	if err := row.Scan(&role.ID, &role.Slug, &role.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("scan role: %w", err)
	}

	return role, nil
}

func (r *repository) CreateUser(ctx context.Context, user CreateUser) (*CreatedUser, error) {
	q := database.Query{
		Name: "user.CreateUser",
		Sql: `INSERT INTO auth.users (email, password_hash)
			  VALUES ($1, $2)
			  RETURNING id, email`,
	}

	row := r.client.DB().QueryRowContext(
		ctx,
		q,
		user.Email,
		user.PasswordHash,
	)

	u := new(CreatedUser)
	if err := row.Scan(&u.ID, &u.Email); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return u, nil
}

func (r *repository) AssignRole(ctx context.Context, userID string, roleID int) error {
	q := database.Query{
		Name: "user.AssignRole",
		Sql: `INSERT INTO auth.user_roles (user_id, role_id)
			  VALUES ($1, $2)`,
	}

	_, err := r.client.DB().ExecContext(
		ctx,
		q,
		userID,
		roleID,
	)

	if err != nil {
		return fmt.Errorf("assign role to user: %w", err)
	}

	return nil
}