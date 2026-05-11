package admin

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/ua-academy-projects/share-bite/internal/admin-auth/dto"
	apperr "github.com/ua-academy-projects/share-bite/internal/admin-auth/error"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/repository/user"
	"github.com/ua-academy-projects/share-bite/pkg/database"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

type CustomerServiceClient interface {
	GetCustomerByUserID(ctx context.Context, userID string) (*dto.CustomerProfileData, error)
}

type BusinessServiceClient interface {
	GetBusinessByUserID(ctx context.Context, userID string) (*dto.BusinessProfileData, error)
}

type Service interface {
	GetUserDetails(ctx context.Context, userID string) (*dto.FullUserDetails, error)
	GetUsersList(ctx context.Context, filter dto.AdminUserFilter) (*dto.PaginatedAdminUsersResponse, error)
	ChangeUserRole(ctx context.Context, targetUserID string, newRoleSlug string) error
}

type service struct {
	adminRepo       user.AdminRepository
	authRepo        user.AuthRepository
	customerService CustomerServiceClient
	businessService BusinessServiceClient
	txManager       database.TxManager
}

func NewService(
	adminRepo user.AdminRepository,
	authRepo user.AuthRepository,
	customerSvc CustomerServiceClient,
	businessSvc BusinessServiceClient,
	txManager database.TxManager,
) Service {
	return &service{
		adminRepo:       adminRepo,
		authRepo:        authRepo,
		customerService: customerSvc,
		businessService: businessSvc,
		txManager:       txManager,
	}
}

func (s *service) GetUserDetails(ctx context.Context, userID string) (*dto.FullUserDetails, error) {
	u, err := s.adminRepo.GetAdminUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, apperr.ErrUserNotFound) {
			return nil, err
		}
		return nil, apperr.Wrap(http.StatusInternalServerError, "failed to fetch base user details", err)
	}

	switch u.RoleSlug {
	case "user":
		customerProfile, err := s.customerService.GetCustomerByUserID(ctx, userID)
		if err != nil {
			logger.ErrorKV(ctx, "failed to fetch customer profile", "user_id", userID, "error", err.Error())
		}
		u.CustomerProfile = customerProfile

	case "business":
		businessProfile, err := s.businessService.GetBusinessByUserID(ctx, userID)
		if err != nil {
			logger.ErrorKV(ctx, "failed to fetch business profile", "user_id", userID, "error", err.Error())
		}
		u.BusinessProfile = businessProfile
	}

	return u, nil
}

func (s *service) GetUsersList(ctx context.Context, filter dto.AdminUserFilter) (*dto.PaginatedAdminUsersResponse, error) {
	items, totalCount, err := s.adminRepo.GetAdminUsersList(ctx, filter)
	if err != nil {
		return nil, apperr.Wrap(http.StatusInternalServerError, "failed to get users list", err)
	}

	return &dto.PaginatedAdminUsersResponse{
		Items:      items,
		TotalCount: totalCount,
	}, nil
}

func (s *service) ChangeUserRole(ctx context.Context, targetUserID string, newRoleSlug string) error {
	targetUser, err := s.adminRepo.GetAdminUserByID(ctx, targetUserID)
	if err != nil {
		if errors.Is(err, apperr.ErrUserNotFound) {
			return err
		}
		return apperr.Wrap(http.StatusInternalServerError, "failed to fetch target user", err)
	}

	currentRole := targetUser.RoleSlug

	if currentRole == newRoleSlug {
		return apperr.New(http.StatusBadRequest, "user already has this role")
	}

	if currentRole == "admin" {
		return apperr.New(http.StatusForbidden, "cannot modify an admin account")
	}

	if currentRole == "business" {
		return apperr.New(
			http.StatusConflict,
			"cannot change role: business accounts are tied to a profile and cannot be converted to other roles",
		)
	}

	if newRoleSlug == "business" {
		return apperr.New(
			http.StatusConflict,
			"cannot change role: standard accounts cannot be converted to business accounts",
		)
	}

	isValidTransition := false
	switch currentRole {
	case "user":
		if newRoleSlug == "moderator" {
			isValidTransition = true
		}
	case "moderator":
		if newRoleSlug == "admin" || newRoleSlug == "user" {
			isValidTransition = true
		}
	}

	if !isValidTransition {
		msg := fmt.Sprintf("invalid role transition from '%s' to '%s'", currentRole, newRoleSlug)
		return apperr.New(http.StatusBadRequest, msg)
	}

	role, err := s.authRepo.FindRoleBySlug(ctx, newRoleSlug)
	if err != nil {
		return apperr.Wrap(http.StatusInternalServerError, "failed to query role", err)
	}

	if role == nil {
		return apperr.ErrRoleNotFound
	}

	txErr := s.txManager.ReadCommitted(ctx, func(txCtx context.Context) error {
		if err := s.adminRepo.UpdateUserRole(txCtx, targetUserID, role.ID); err != nil {
			return apperr.Wrap(http.StatusInternalServerError, "failed to update user role", err)
		}

		if err := s.authRepo.RevokeAllUserTokens(txCtx, targetUserID); err != nil {
			return apperr.Wrap(http.StatusInternalServerError, "failed to revoke existing sessions during role change", err)
		}

		return nil
	})

	return txErr
}
