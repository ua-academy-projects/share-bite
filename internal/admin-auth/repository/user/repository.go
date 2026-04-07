package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/admin-auth/dto"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database"
)

type Repository interface {
	FindByEmail(ctx context.Context, email string) (*dto.UserWithRole, error)
	FindRoleBySlug(ctx context.Context, slug string) (*entity.Role, error)
	CreateWithRole(ctx context.Context, params dto.CreateWithRoleParams) (*dto.CreatedUser, error)
}

type repository struct {
	client database.Client
}

func New(client database.Client) Repository {
	return &repository{client: client}
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*dto.UserWithRole, error) {
	q := database.Query{
		Name: "user.FindByEmail",
		Sql: `
			SELECT u.id, u.email, u.password_hash, r.slug
			FROM auth.users u
			LEFT JOIN auth.user_roles ur ON u.id = ur.user_id
			LEFT JOIN auth.roles r ON ur.role_id = r.id
			WHERE u.email = $1
		`,
	}

	row := r.client.DB().QueryRowContext(ctx, q, email)
	u := new(dto.UserWithRole)
	if err := row.Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.RoleSlug,
	); err != nil {
		return nil, err
	}

	return u, nil
}

func (r *repository) FindRoleBySlug(ctx context.Context, slug string) (*entity.Role, error) {
	q := database.Query{
		Name: "user.FindRoleBySlug",
		Sql:  `SELECT id, slug, name FROM auth.roles WHERE slug = $1`,
	}

	row := r.client.DB().QueryRowContext(ctx, q, slug)

	role := new(entity.Role)
	if err := row.Scan(&role.ID, &role.Slug, &role.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("scan role: %w", err)
	}

	return role, nil
}

func (r *repository) CreateWithRole(ctx context.Context, params dto.CreateWithRoleParams) (*dto.CreatedUser, error) {
	q := database.Query{
		Name: "user.CreateWithRole",
		Sql: `
			WITH created_user AS (
				INSERT INTO auth.users (email, password_hash)
				VALUES ($1, $2)
				RETURNING id, email
			),
			assigned_role AS (
				INSERT INTO auth.user_roles (user_id, role_id)
				SELECT id, $3
				FROM created_user
			)
			SELECT id, email
			FROM created_user
		`,
	}

	row := r.client.DB().QueryRowContext(
		ctx,
		q,
		params.Email,
		params.PasswordHash,
		params.RoleID,
	)

	u := new(dto.CreatedUser)
	if err := row.Scan(&u.ID, &u.Email); err != nil {
		return nil, fmt.Errorf("create user with role: %w", err)
	}

	return u, nil
}
