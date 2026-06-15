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

type PlatformStatisticsResponse struct {
	TotalUsers                   int64 `json:"total_users"`
	TotalAdminUsers              int64 `json:"total_admin_users"`
	TotalModeratorUsers          int64 `json:"total_moderator_users"`
	TotalRegularUsers            int64 `json:"total_regular_users"`
	TotalBusinessRoleUsers       int64 `json:"total_business_role_users"`
	TotalActiveUsers             int64 `json:"total_active_users"`
	TotalMutedUsers              int64 `json:"total_muted_users"`
	TotalSuspendedUsers          int64 `json:"total_suspended_users"`
	TotalCustomers               int64 `json:"total_customers"`
	TotalGuestPosts              int64 `json:"total_guest_posts"`
	TotalGuestComments           int64 `json:"total_guest_comments"`
	TotalGuestPostLikes          int64 `json:"total_guest_post_likes"`
	TotalCollections             int64   `json:"total_collections"`
	AvgPostsPerCustomer          float64 `json:"avg_posts_per_customer"`
	AvgCommentsPerCustomer       float64 `json:"avg_comments_per_customer"`
	AvgCommentsPerPost           float64 `json:"avg_comments_per_post"`
	CollectionsWithCollaborators int64   `json:"collections_with_collaborators"`
	PostsWithCollaborators       int64   `json:"posts_with_collaborators"`
	TotalBusinessOrgUnits        int64 `json:"total_business_org_units"`
	TotalBusinessPosts           int64 `json:"total_business_posts"`
	TotalBusinessComments        int64 `json:"total_business_comments"`
	TotalBusinessLikes           int64 `json:"total_business_likes"`
	TotalBusinessBoxes           int64   `json:"total_business_boxes"`
	TotalBusinessBoxItems        int64   `json:"total_business_box_items"`
	AvgPostsPerBusiness          float64 `json:"avg_posts_per_business"`
	AvgCommentsPerBusiness       float64 `json:"avg_comments_per_business"`
	AvgBusinessCommentsPerPost   float64 `json:"avg_business_comments_per_post"`
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
	Status      string   `json:"status"`
}

type GitHubUser struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type User struct {
	ID        string
	GitHubID  int64
	Login     string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PendingBusinessListItem struct {
	ID           int    `json:"id"`
	OrgAccountID string `json:"org_account_id"`
	Name         string `json:"name"`
	Avatar       string `json:"avatar"`
	Description  string `json:"description"`
	Status       string `json:"status"`
}

type PaginatedPendingBusinessesResponse struct {
	Items      []PendingBusinessListItem `json:"items"`
	TotalCount int                       `json:"total_count"`
}

type ReviewBusinessParams struct {
	OrgUnitID int     `json:"orgUnitId"`
	AdminID   string  `json:"adminId" binding:"required,uuid"`
	NewStatus string  `json:"newStatus"`
	Comment   *string `json:"comment"`
}
