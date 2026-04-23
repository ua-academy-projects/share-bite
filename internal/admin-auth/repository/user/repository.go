package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/dto"
	apperr "github.com/ua-academy-projects/share-bite/internal/admin-auth/error"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/models"
	"github.com/ua-academy-projects/share-bite/pkg/database"
)

type Repository interface {
	FindByEmail(ctx context.Context, email string) (*dto.UserWithRole, error)
	FindRoleBySlug(ctx context.Context, slug string) (*models.Role, error)
	CreateUser(ctx context.Context, user CreateUser) (*CreatedUser, error)
	AssignRole(ctx context.Context, userID string, roleID int) error
	FindBySocialProvider(ctx context.Context, provider, providerID string) (*dto.UserWithRole, error)
	CreateWithSocial(ctx context.Context, params dto.CreateUserWithSocialParams) (*dto.CreatedUser, error)
	LinkSocialAccount(ctx context.Context, params dto.CreateSocialAccountParams) error
	CreatePasswordResetToken(ctx context.Context, params dto.CreatePasswordResetTokenParams) error
	ResetPassword(ctx context.Context, tokenHash, passwordHash string) (string, bool, error)
	StoreRefreshToken(ctx context.Context, params dto.StoreRefreshTokenParams) error
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
	GetUserIDByRefreshToken(ctx context.Context, tokenHash string) (string, error)
	RevokeAllUserTokens(ctx context.Context, userID string) error
	CountActiveSessions(ctx context.Context, userID string) (int, error)
	DeleteOldestSession(ctx context.Context, userID string) error
	DeleteExpiredTokens(ctx context.Context) error
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

func (r *repository) CreatePasswordResetToken(ctx context.Context, params dto.CreatePasswordResetTokenParams) error {
	invalidateQuery := database.Query{
		Name: "user.InvalidatePasswordResetTokens",
		Sql: `
			UPDATE auth.password_reset_tokens
			SET used_at = NOW()
			WHERE user_id = $1
			  AND used_at IS NULL
		`,
	}

	insertQuery := database.Query{
		Name: "user.InsertPasswordResetToken",
		Sql: `
			INSERT INTO auth.password_reset_tokens (user_id, token_hash, expires_at)
			VALUES ($1, $2, $3)
		`,
	}

	if _, err := r.client.DB().ExecContext(ctx, invalidateQuery, params.UserID); err != nil {
		return fmt.Errorf("invalidate previous password reset tokens: %w", err)
	}

	if _, err := r.client.DB().ExecContext(
		ctx,
		insertQuery,
		params.UserID,
		params.TokenHash,
		params.ExpiresAt,
	); err != nil {
		return fmt.Errorf("insert password reset token: %w", err)
	}

	return nil
}

func (r *repository) ResetPassword(ctx context.Context, tokenHash, passwordHash string) (string, bool, error) {
	var userID string

	consumeTokenQuery := database.Query{
		Name: "user.ConsumePasswordResetToken",
		Sql: `
			UPDATE auth.password_reset_tokens
			SET used_at = NOW()
			WHERE token_hash = $1
			  AND used_at IS NULL
			  AND expires_at > NOW()
			RETURNING user_id
		`,
	}

	updatePasswordQuery := database.Query{
		Name: "user.UpdatePasswordByUserID",
		Sql: `
			UPDATE auth.users
			SET password_hash = $2,
			    updated_at = NOW()
			WHERE id = $1
		`,
	}

	if err := r.client.DB().QueryRowContext(ctx, consumeTokenQuery, tokenHash).Scan(&userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", false, apperr.ErrInvalidResetToken
		}

		return "", false, fmt.Errorf("consume password reset token: %w", err)
	}

	if _, err := r.client.DB().ExecContext(ctx, updatePasswordQuery, userID, passwordHash); err != nil {
		return "", false, fmt.Errorf("update password: %w", err)
	}

	return userID, true, nil
}

func (r *repository) FindRoleBySlug(ctx context.Context, slug string) (*models.Role, error) {
	q := database.Query{
		Name: "user.FindRoleBySlug",
		Sql:  `SELECT id, slug, name FROM auth.roles WHERE slug = $1`,
	}

	row := r.client.DB().QueryRowContext(ctx, q, slug)

	role := new(models.Role)
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

func (r *repository) StoreRefreshToken(ctx context.Context, params dto.StoreRefreshTokenParams) error {
	q := database.Query{
		Name: "user.StoreRefreshToken",
		Sql: `INSERT INTO auth.refresh_tokens (user_id, token_hash, expires_at)
			  VALUES ($1, $2, $3)`,
	}
	_, err := r.client.DB().ExecContext(ctx, q, params.UserID, params.TokenHash, params.ExpiresAt)
	if err != nil {
		return apperr.Wrap(http.StatusInternalServerError, "failed to store refresh token", err)
	}
	return nil
}

func (r *repository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	q := database.Query{
		Name: "user.RevokeRefreshToken",
		Sql:  `UPDATE auth.refresh_tokens SET revoked_at = NOW() WHERE token_hash = $1 AND revoked_at IS NULL RETURNING token_hash`,
	}
	var returnedHash string
	err := r.client.DB().QueryRowContext(ctx, q, tokenHash).Scan(&returnedHash)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperr.ErrInvalidToken
		}
		return apperr.Wrap(http.StatusInternalServerError, "failed to revoke refresh token", err)
	}

	return nil
}

func (r *repository) GetUserIDByRefreshToken(ctx context.Context, tokenHash string) (string, error) {
	q := database.Query{
		Name: "user.GetUserIDByRefreshToken",
		Sql:  `SELECT user_id FROM auth.refresh_tokens WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > NOW()`,
	}
	var userID string
	row := r.client.DB().QueryRowContext(ctx, q, tokenHash)
	if err := row.Scan(&userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", apperr.ErrInvalidToken
		}
		return "", apperr.Wrap(http.StatusInternalServerError, "failed to fetch user id", err)
	}
	return userID, nil
}

func (r *repository) RevokeAllUserTokens(ctx context.Context, userID string) error {
	q := database.Query{
		Name: "user.RevokeAllUserTokens",
		Sql:  `UPDATE auth.refresh_tokens SET revoked_at = NOW() WHERE user_id = $1 AND revoked_at IS NULL`,
	}
	_, err := r.client.DB().ExecContext(ctx, q, userID)
	if err != nil {
		return apperr.Wrap(http.StatusInternalServerError, "failed to revoke all user tokens", err)
	}
	return nil
}

func (r *repository) CountActiveSessions(ctx context.Context, userID string) (int, error) {
	q := database.Query{
		Name: "user.CountActiveSessions",
		Sql: `
          SELECT COUNT(*) 
          FROM auth.refresh_tokens 
          WHERE user_id = $1 
            AND revoked_at IS NULL 
            AND expires_at > NOW()
       `,
	}

	var count int
	err := r.client.DB().QueryRowContext(ctx, q, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count active sessions: %w", err)
	}

	return count, nil
}

func (r *repository) DeleteOldestSession(ctx context.Context, userID string) error {
	q := database.Query{
		Name: "user.DeleteOldestSession",
		Sql: `
         DELETE FROM auth.refresh_tokens 
			WHERE id = (
   			 SELECT id FROM auth.refresh_tokens 
    		 WHERE user_id = $1 
      		 AND revoked_at IS NULL 
      		 AND expires_at > NOW()
    		 ORDER BY created_at ASC 
    		 LIMIT 1
			)
       `,
	}

	_, err := r.client.DB().ExecContext(ctx, q, userID)
	if err != nil {
		return fmt.Errorf("delete oldest session: %w", err)
	}

	return nil
}

func (r *repository) DeleteExpiredTokens(ctx context.Context) error {
	q := database.Query{
		Name: "user.DeleteExpiredTokens",
		Sql: `
          DELETE FROM auth.refresh_tokens 
          WHERE expires_at < NOW() - INTERVAL '3 days' 
             OR (revoked_at IS NOT NULL AND revoked_at < NOW() - INTERVAL '3 days')
       `,
	}

	_, err := r.client.DB().ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("cleanup expired tokens: %w", err)
	}

	return nil
}
