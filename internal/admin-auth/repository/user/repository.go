package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/dto"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/entity"
	apperr "github.com/ua-academy-projects/share-bite/internal/admin-auth/error"
	"github.com/ua-academy-projects/share-bite/pkg/database"
)

type Repository interface {
	FindByEmail(ctx context.Context, email string) (*dto.UserWithRole, error)
	FindRoleBySlug(ctx context.Context, slug string) (*entity.Role, error)
	CreateWithRole(ctx context.Context, params dto.CreateWithRoleParams) (*dto.CreatedUser, error)
	FindBySocialProvider(ctx context.Context, provider, providerID string) (*dto.UserWithRole, error)
	CreateWithSocial(ctx context.Context, params dto.CreateUserWithSocialParams) (*dto.CreatedUser, error)
	LinkSocialAccount(ctx context.Context, params dto.CreateSocialAccountParams) error
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("find by email: %w", err)
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

func (r *repository) FindBySocialProvider(ctx context.Context, provider, providerID string) (*dto.UserWithRole, error) {
	q := database.Query{
		Name: "user.FindBySocialProvider",
		Sql: `
			SELECT u.id, u.email, u.password_hash, r.slug
			FROM auth.users u
			JOIN auth.social_accounts sa
				on sa.user_id = u.id
			LEFT JOIN auth.user_roles ur 
			    ON ur.user_id = u.id
			LEFT JOIN auth.roles r 
			    ON r.id = ur.role_id
			WHERE sa.provider = $1
				AND sa.provider_id = $2
		`,
	}
	row := r.client.DB().QueryRowContext(ctx, q, provider, providerID)
	u := new(dto.UserWithRole)
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.RoleSlug); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("find user with social provider: %w", err)
	}
	return u, nil
}

func (r *repository) CreateWithSocial(ctx context.Context, params dto.CreateUserWithSocialParams) (*dto.CreatedUser, error) {
	q := database.Query{
		Name: "user.CreateWithSocial",
		Sql: `
			WITH created_user AS (
				INSERT INTO auth.users (email, password_hash)
				VALUES ($1, NULL)
				RETURNING id, email
			),
			assigned_role AS (
				INSERT INTO auth.user_roles (user_id, role_id)
				SELECT id, $2
				FROM created_user
			),
			linked_social AS (
				INSERT INTO auth.social_accounts (user_id, provider, provider_id, email)
				SELECT id, $3, $4, $1
				FROM created_user
			)
			SELECT id, email FROM created_user
		`}
	row := r.client.DB().QueryRowContext(ctx, q, params.Email, params.RoleID, params.Provider, params.ProviderID)
	u := new(dto.CreatedUser)
	if err := row.Scan(&u.ID, &u.Email); err != nil {

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {

			if pgErr.ConstraintName == "uq_provider_account" {
				return nil, apperr.ErrProviderAlreadyLinked
			}

			if pgErr.ConstraintName == "users_email_key" {
				return nil, apperr.ErrUserAlreadyExists
			}
		}

		return nil, apperr.Wrap(http.StatusInternalServerError, "failed to create user with social", err)
	}

	return u, nil
}

func (r *repository) LinkSocialAccount(ctx context.Context, params dto.CreateSocialAccountParams) error {
	q := database.Query{
		Name: "user.LinkSocialAccount",
		Sql: `
		INSERT INTO auth.social_accounts (user_id, provider, provider_id, email)
		VALUES ($1, $2, $3, $4)`,
	}
	_, err := r.client.DB().ExecContext(ctx, q,
		params.UserID,
		params.Provider,
		params.ProviderID,
		params.Email,
	)
	if err != nil {

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "uq_provider_account" {
				return apperr.ErrProviderAlreadyLinked
			}
		}

		return apperr.Wrap(http.StatusInternalServerError, "failed to link social account", err)
	}

	return nil
}
