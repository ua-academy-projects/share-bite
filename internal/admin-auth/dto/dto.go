package dto

import (
	"time"

	"github.com/ua-academy-projects/share-bite/internal/admin-auth/models"
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
	models.User
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

type StoreRefreshTokenParams struct {
	TokenHash string
	UserID    string
	ExpiresAt time.Time
}

type AdminUserFilter struct {
	SearchEmail string
	RoleSlug    string
	Status      string
	Limit       int
	Offset      int
	SortOrder   string
}

type AdminUserListItem struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	RoleSlug  string    `json:"role_slug"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
type PaginatedAdminUsersResponse struct {
	Items      []AdminUserListItem `json:"items"`
	TotalCount int                 `json:"total_count"`
}

type ChangeRoleRequest struct {
	RoleSlug string `json:"role_slug" validate:"required"`
}

type FullUserDetails struct {
	ID              string               `json:"id"`
	Email           string               `json:"email"`
	RoleSlug        string               `json:"role_slug"`
	Status          string               `json:"status"`
	CreatedAt       time.Time            `json:"created_at"`
	CustomerProfile *CustomerProfileData `json:"customer_profile,omitempty"`
	BusinessProfile *BusinessProfileData `json:"business_profile,omitempty"`
}

type CustomerProfileData struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	AvatarKey string `json:"avatar_object_key"`
	Bio       string `json:"bio"`
}

type BusinessProfileData struct {
	ProfileType string   `json:"profile_type"`
	Name        string   `json:"name"`
	Avatar      string   `json:"avatar"`
	Banner      string   `json:"banner"`
	Description string   `json:"description"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
}

type GitHubUser struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type User struct {
	ID        int64
	GitHubID  int64
	Login     string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
