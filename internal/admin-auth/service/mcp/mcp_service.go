package mcp

import (
	"context"
	"errors"
	"fmt"

	apperr "github.com/ua-academy-projects/share-bite/internal/admin-auth/error"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/repository/user"
)

type PermissionService struct {
	adminRepo       user.AdminRepository
	rolePermissions map[string][]string
}

func NewMCPPermissionService(adminRepo user.AdminRepository) *PermissionService {
	return &PermissionService{
		adminRepo: adminRepo,
		rolePermissions: map[string][]string{
			"admin":     {"admin_auth_health_check", "get_current_admin_context", "validate_admin_permissions"},
			"moderator": {"admin_auth_health_check", "get_current_admin_context"},
			"user":      {},
			"business":  {},
		},
	}
}

func (s *PermissionService) HasPermission(ctx context.Context, userID string, permission string) (bool, error) {
	if userID == "" || permission == "" {
		return false, nil
	}

	targetUser, err := s.adminRepo.GetAdminUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, apperr.ErrUserNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("mcp permission check failed to fetch user: %w", err)
	}

	allowedPermissions, exists := s.rolePermissions[targetUser.RoleSlug]
	if !exists {
		return false, nil
	}

	for _, p := range allowedPermissions {
		if p == permission {
			return true, nil
		}
	}

	return false, nil
}
