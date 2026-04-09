package dto

import (
	"time"

	"github.com/ua-academy-projects/share-bite/internal/admin-auth/entity"
)

type CreateWithRoleParams struct {
	Email        string
	PasswordHash string
	RoleID       int
}

type CreatedUser struct {
	ID    string
	Email string
}

type CreatePasswordResetTokenParams struct {
	UserID    string
	TokenHash string
	ExpiresAt time.Time
}

type UserWithRole struct {
	entity.User
	RoleSlug string
}
