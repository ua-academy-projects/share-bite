package user

import (
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/admin-auth/dto"
)

type PlatformStatistics struct {
	TotalUsers                   int64   `db:"total_users"`
	TotalAdminUsers              int64   `db:"total_admin_users"`
	TotalModeratorUsers          int64   `db:"total_moderator_users"`
	TotalRegularUsers            int64   `db:"total_regular_users"`
	TotalBusinessRoleUsers       int64   `db:"total_business_role_users"`
	TotalActiveUsers             int64   `db:"total_active_users"`
	TotalMutedUsers              int64   `db:"total_muted_users"`
	TotalSuspendedUsers          int64   `db:"total_suspended_users"`
	TotalCustomers               int64   `db:"total_customers"`
	TotalGuestPosts              int64   `db:"total_guest_posts"`
	TotalGuestComments           int64   `db:"total_guest_comments"`
	TotalGuestPostLikes          int64   `db:"total_guest_post_likes"`
	TotalCollections             int64   `db:"total_collections"`
	AvgPostsPerCustomer          float64 `db:"avg_posts_per_customer"`
	AvgCommentsPerCustomer       float64 `db:"avg_comments_per_customer"`
	AvgCommentsPerPost           float64 `db:"avg_comments_per_post"`
	CollectionsWithCollaborators int64   `db:"collections_with_collaborators"`
	PostsWithCollaborators       int64   `db:"posts_with_collaborators"`
	TotalBusinessOrgUnits        int64   `db:"total_business_org_units"`
	TotalBusinessPosts           int64   `db:"total_business_posts"`
	TotalBusinessComments        int64   `db:"total_business_comments"`
	TotalBusinessLikes           int64   `db:"total_business_likes"`
	TotalBusinessBoxes           int64   `db:"total_business_boxes"`
	TotalBusinessBoxItems        int64   `db:"total_business_box_items"`
	AvgPostsPerBusiness          float64 `db:"avg_posts_per_business"`
	AvgCommentsPerBusiness       float64 `db:"avg_comments_per_business"`
	AvgBusinessCommentsPerPost   float64 `db:"avg_business_comments_per_post"`
}

func (s PlatformStatistics) ToDTO() dto.PlatformStatisticsResponse {
	return dto.PlatformStatisticsResponse{
		TotalUsers:                   s.TotalUsers,
		TotalAdminUsers:              s.TotalAdminUsers,
		TotalModeratorUsers:          s.TotalModeratorUsers,
		TotalRegularUsers:            s.TotalRegularUsers,
		TotalBusinessRoleUsers:       s.TotalBusinessRoleUsers,
		TotalActiveUsers:             s.TotalActiveUsers,
		TotalMutedUsers:              s.TotalMutedUsers,
		TotalSuspendedUsers:          s.TotalSuspendedUsers,
		TotalCustomers:               s.TotalCustomers,
		TotalGuestPosts:              s.TotalGuestPosts,
		TotalGuestComments:           s.TotalGuestComments,
		TotalGuestPostLikes:          s.TotalGuestPostLikes,
		TotalCollections:             s.TotalCollections,
		AvgPostsPerCustomer:          s.AvgPostsPerCustomer,
		AvgCommentsPerCustomer:       s.AvgCommentsPerCustomer,
		AvgCommentsPerPost:           s.AvgCommentsPerPost,
		CollectionsWithCollaborators: s.CollectionsWithCollaborators,
		PostsWithCollaborators:       s.PostsWithCollaborators,
		TotalBusinessOrgUnits:        s.TotalBusinessOrgUnits,
		TotalBusinessPosts:           s.TotalBusinessPosts,
		TotalBusinessComments:        s.TotalBusinessComments,
		TotalBusinessLikes:           s.TotalBusinessLikes,
		TotalBusinessBoxes:           s.TotalBusinessBoxes,
		TotalBusinessBoxItems:        s.TotalBusinessBoxItems,
		AvgPostsPerBusiness:          s.AvgPostsPerBusiness,
		AvgCommentsPerBusiness:       s.AvgCommentsPerBusiness,
		AvgBusinessCommentsPerPost:   s.AvgBusinessCommentsPerPost,
	}
}

func executeSQLError(err error) error {
	return fmt.Errorf("execute sql: %w", err)
}

func scanRowError(err error) error {
	return fmt.Errorf("scan row: %w", err)
}
