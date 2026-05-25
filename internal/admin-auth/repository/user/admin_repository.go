package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/dto"
	apperr "github.com/ua-academy-projects/share-bite/internal/admin-auth/error"
	"github.com/ua-academy-projects/share-bite/pkg/database"
)

type AdminRepository interface {
	GetAdminUsersList(ctx context.Context, filter dto.AdminUserFilter) ([]dto.AdminUserListItem, int, error)
	UpdateUserRole(ctx context.Context, userID string, roleID int) error
	GetAdminUserByID(ctx context.Context, userID string) (*dto.FullUserDetails, error)
	GetPlatformStatistics(ctx context.Context) (*dto.PlatformStatisticsResponse, error)
}

type adminRepository struct {
	client database.Client
}

func NewAdmin(client database.Client) AdminRepository {
	return &adminRepository{client: client}
}

func (r *adminRepository) GetAdminUsersList(ctx context.Context, filter dto.AdminUserFilter) ([]dto.AdminUserListItem, int, error) {
	sortDir := "DESC"
	if strings.ToUpper(filter.SortOrder) == "ASC" {
		sortDir = "ASC"
	}

	queryText := fmt.Sprintf(`
       SELECT 
          u.id, 
          u.email, 
          r.slug as role_slug, 
          u.status, 
          u.created_at,
          COUNT(*) OVER() AS total_count
       FROM auth.users u
       JOIN auth.user_roles ur ON u.id = ur.user_id
       JOIN auth.roles r ON ur.role_id = r.id
       WHERE ($1 = '' OR u.email ILIKE '%%' || $1 || '%%')
         AND ($2 = '' OR r.slug = $2)
         AND ($3 = '' OR u.status = $3)
       ORDER BY u.created_at %s
       LIMIT $4 OFFSET $5
    `, sortDir)

	q := database.Query{
		Name: "admin.GetAdminUsersList",
		Sql:  queryText,
	}

	rows, err := r.client.DB().QueryContext(ctx, q, filter.SearchEmail, filter.RoleSlug, filter.Status, filter.Limit, filter.Offset)
	if err != nil {
		return nil, 0, apperr.Wrap(http.StatusInternalServerError, "failed to get admin users list", err)
	}
	defer rows.Close()

	var users []dto.AdminUserListItem
	var totalCount int

	for rows.Next() {
		var u dto.AdminUserListItem
		if err := rows.Scan(&u.ID, &u.Email, &u.RoleSlug, &u.Status, &u.CreatedAt, &totalCount); err != nil {
			return nil, 0, apperr.Wrap(http.StatusInternalServerError, "failed to scan admin user list item", err)
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, apperr.Wrap(http.StatusInternalServerError, "rows iteration error", err)
	}

	if users == nil {
		users = make([]dto.AdminUserListItem, 0)
	}

	return users, totalCount, nil
}

func (r *adminRepository) UpdateUserRole(ctx context.Context, userID string, roleID int) error {
	q := database.Query{
		Name: "admin.UpdateUserRole",
		Sql: `
          WITH deleted AS (
             DELETE FROM auth.user_roles WHERE user_id = $1
          )
          INSERT INTO auth.user_roles (user_id, role_id) 
          VALUES ($1, $2)
       `,
	}

	_, err := r.client.DB().ExecContext(ctx, q, userID, roleID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				if pgErr.ConstraintName == "fk_user_roles_role_id" {
					return apperr.ErrRoleNotFound
				}
				if pgErr.ConstraintName == "fk_user_roles_user_id" {
					return apperr.ErrUserNotFound
				}
			}
		}
		return apperr.Wrap(http.StatusInternalServerError, "failed to update user role", err)
	}

	return nil
}

func (r *adminRepository) GetAdminUserByID(ctx context.Context, userID string) (*dto.FullUserDetails, error) {
	q := database.Query{
		Name: "admin.GetAdminUserByID",
		Sql: `
          SELECT u.id, u.email, r.slug, u.status, u.created_at
          FROM auth.users u
          JOIN auth.user_roles ur ON u.id = ur.user_id
          JOIN auth.roles r ON ur.role_id = r.id
          WHERE u.id = $1
       `,
	}

	var user dto.FullUserDetails
	err := r.client.DB().QueryRowContext(ctx, q, userID).Scan(
		&user.ID, &user.Email, &user.RoleSlug, &user.Status, &user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.ErrUserNotFound
		}
		return nil, apperr.Wrap(http.StatusInternalServerError, "failed to get admin user by id", err)
	}

	return &user, nil
}

func (r *adminRepository) GetPlatformStatistics(ctx context.Context) (*dto.PlatformStatisticsResponse, error) {
	q := database.Query{
		Name: "admin.GetPlatformStatistics",
		Sql: `
          SELECT
              total_users,
              total_admin_users,
              total_moderator_users,
              total_regular_users,
              total_business_role_users,
              total_active_users,
              total_muted_users,
              total_suspended_users,
              total_customers,
              total_guest_posts,
              total_guest_comments,
              total_guest_post_likes,
              total_collections,
              total_collection_venues,
              total_collection_collaborators,
              total_collection_invitations,
              total_customer_follows,
              total_business_org_units,
              total_business_posts,
              total_business_comments,
              total_business_likes,
              total_business_boxes,
              total_business_box_items
          FROM analytics.platform_statistics
		  
       `,
	}

	var stats dto.PlatformStatisticsResponse
	err := r.client.DB().QueryRowContext(ctx, q).Scan(
		&stats.TotalUsers,
		&stats.TotalAdminUsers,
		&stats.TotalModeratorUsers,
		&stats.TotalRegularUsers,
		&stats.TotalBusinessRoleUsers,
		&stats.TotalActiveUsers,
		&stats.TotalMutedUsers,
		&stats.TotalSuspendedUsers,
		&stats.TotalCustomers,
		&stats.TotalGuestPosts,
		&stats.TotalGuestComments,
		&stats.TotalGuestPostLikes,
		&stats.TotalCollections,
		&stats.TotalCollectionVenues,
		&stats.TotalCollectionCollaborators,
		&stats.TotalCollectionInvitations,
		&stats.TotalCustomerFollows,
		&stats.TotalBusinessOrgUnits,
		&stats.TotalBusinessPosts,
		&stats.TotalBusinessComments,
		&stats.TotalBusinessLikes,
		&stats.TotalBusinessBoxes,
		&stats.TotalBusinessBoxItems,
	)
	if err != nil {
		return nil, apperr.Wrap(http.StatusInternalServerError, "failed to get platform statistics", err)
	}

	return &stats, nil
}
