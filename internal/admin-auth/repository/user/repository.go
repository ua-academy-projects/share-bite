package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/admin-auth/dto"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/admin-auth/error"
	"github.com/ua-academy-projects/share-bite/pkg/database"
)

type Repository interface {
	FindByEmail(ctx context.Context, email string) (*dto.UserWithRole, error)
	FindRoleBySlug(ctx context.Context, slug string) (*entity.Role, error)
	CreateWithRole(ctx context.Context, params dto.CreateWithRoleParams) (*dto.CreatedUser, error)
	CreatePasswordResetToken(ctx context.Context, params dto.CreatePasswordResetTokenParams) error
	ResetPassword(ctx context.Context, tokenHash, passwordHash string) (bool, error)
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
			return nil, apperror.UserNotFoundEmail(email)
		}
		return nil, err
	}

	return u, nil
}

func (r *repository) CreatePasswordResetToken(ctx context.Context, params dto.CreatePasswordResetTokenParams) error {
	q := database.Query{
		Name: "user.CreatePasswordResetToken",
		Sql: `
			WITH invalidated AS (
				UPDATE auth.password_reset_tokens
				SET used_at = NOW()
				WHERE user_id = $1
				  AND used_at IS NULL
			),
			inserted AS (
				INSERT INTO auth.password_reset_tokens (user_id, token_hash, expires_at)
				VALUES ($1, $2, $3)
				RETURNING id
			)
			SELECT id
			FROM inserted
		`,
	}

	var tokenID string
	if err := r.client.DB().QueryRowContext(
		ctx,
		q,
		params.UserID,
		params.TokenHash,
		params.ExpiresAt,
	).Scan(&tokenID); err != nil {
		return fmt.Errorf("create password reset token: %w", err)
	}

	return nil
}

func (r *repository) ResetPassword(ctx context.Context, tokenHash, passwordHash string) (bool, error) {
	q := database.Query{
		Name: "user.ResetPassword",
		Sql: `
			WITH consumed AS (
				UPDATE auth.password_reset_tokens
				SET used_at = NOW()
				WHERE token_hash = $1
				  AND used_at IS NULL
				  AND expires_at > NOW()
				RETURNING user_id
			),
			updated_user AS (
				UPDATE auth.users AS u
				SET password_hash = $2,
				    updated_at = NOW()
				FROM consumed
				WHERE u.id = consumed.user_id
				RETURNING u.id
			)
			SELECT id
			FROM updated_user
		`,
	}

	var userID string
	if err := r.client.DB().QueryRowContext(ctx, q, tokenHash, passwordHash).Scan(&userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, apperror.ErrInvalidResetToken
		}

		return false, fmt.Errorf("reset password by token hash: %w", err)
	}

	return true, nil
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
