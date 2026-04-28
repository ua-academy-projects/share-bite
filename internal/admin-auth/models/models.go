package models

import (
	"time"
)

type User struct {
	ID           string    `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	PasswordHash *string   `db:"password_hash" json:"-"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type Role struct {
	ID   int    `db:"id" json:"id"`
	Slug string `db:"slug" json:"slug"`
	Name string `db:"name" json:"name"`
}

type UserRole struct {
	UserID string `db:"user_id" json:"user_id"`
	RoleID int    `db:"role_id" json:"role_id"`
}

type SocialAccount struct {
	ID         string    `db:"id"`
	UserID     string    `db:"user_id"`
	Provider   string    `db:"provider"`
	ProviderID string    `db:"provider_id"`
	Email      string    `db:"email"`
	CreatedAt  time.Time `db:"created_at"`
}

type RefreshToken struct {
	ID        string     `db:"id"`
	TokenHash string     `db:"token_hash"`
	UserID    string     `db:"user_id"`
	CreatedAt time.Time  `db:"created_at"`
	ExpiresAt time.Time  `db:"expires_at"`
	RevokedAt *time.Time `db:"revoked_at"`
}
