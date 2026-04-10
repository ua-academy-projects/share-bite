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

type OAuthUserInfo struct {
	Provider      string
	ProviderID    string
	Email         string
	EmailVerified bool
}

// CreateUserWithSocialParams — реєстрація нового юзера через OAuth
type CreateUserWithSocialParams struct {
	Email      string
	Provider   string
	ProviderID string
	RoleID     int
}

// CreateSocialAccountParams — прив'язка провайдера до вже існуючого юзера
type CreateSocialAccountParams struct {
	UserID     string
	Provider   string
	ProviderID string
	Email      string
}
