package admin

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ua-academy-projects/share-bite/internal/admin-auth/dto"
	apperr "github.com/ua-academy-projects/share-bite/internal/admin-auth/error"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/repository/user"
	"github.com/ua-academy-projects/share-bite/pkg/database"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	"github.com/ua-academy-projects/share-bite/pkg/outbox"
)

type CustomerServiceClient interface {
	GetCustomerByUserID(ctx context.Context, userID string) (*dto.CustomerProfileData, error)
}

type BusinessServiceClient interface {
	GetBusinessByUserID(ctx context.Context, userID string) (*dto.BusinessProfileData, error)
	GetPendingBusinesses(ctx context.Context, limit, offset int) ([]dto.PendingBusinessListItem, int, error)
	GetBusinessStatusAndOwner(ctx context.Context, orgUnitID int) (string, string, error)
	ReviewBusiness(ctx context.Context, params dto.ReviewBusinessParams) error
}

type Service interface {
	GetUserDetails(ctx context.Context, userID string) (*dto.FullUserDetails, error)
	GetUsersList(ctx context.Context, filter dto.AdminUserFilter) (*dto.PaginatedAdminUsersResponse, error)
	GetPlatformStatistics(ctx context.Context) (*dto.PlatformStatisticsResponse, error)
	ChangeUserRole(ctx context.Context, targetUserID string, newRoleSlug string) error
	GetPendingBusinessesList(ctx context.Context, limit, offset int) (*dto.PaginatedPendingBusinessesResponse, error)
	ReviewBusinessStatus(ctx context.Context, params dto.ReviewBusinessParams) error
}

type service struct {
	adminRepo       user.AdminRepository
	authRepo        user.AuthRepository
	customerService CustomerServiceClient
	businessService BusinessServiceClient
	outboxWriter    outbox.Writer
	txManager       database.TxManager
}

func NewService(
	adminRepo user.AdminRepository,
	authRepo user.AuthRepository,
	customerSvc CustomerServiceClient,
	businessSvc BusinessServiceClient,
	outboxWriter outbox.Writer,
	txManager database.TxManager,
) Service {
	return &service{
		adminRepo:       adminRepo,
		authRepo:        authRepo,
		customerService: customerSvc,
		businessService: businessSvc,
		outboxWriter:    outboxWriter,
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

func (s *service) GetPlatformStatistics(ctx context.Context) (*dto.PlatformStatisticsResponse, error) {
	stats, err := s.adminRepo.GetPlatformStatistics(ctx)
	if err != nil {
		return nil, apperr.Wrap(http.StatusInternalServerError, "failed to get platform statistics", err)
	}

	return stats, nil
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

func (s *service) GetPendingBusinessesList(ctx context.Context, limit, offset int) (*dto.PaginatedPendingBusinessesResponse, error) {
	items, totalCount, err := s.businessService.GetPendingBusinesses(ctx, limit, offset)
	if err != nil {
		return nil, apperr.Wrap(http.StatusInternalServerError, "failed to get pending businesses list", err)
	}

	return &dto.PaginatedPendingBusinessesResponse{
		Items:      items,
		TotalCount: totalCount,
	}, nil
}

func (s *service) ReviewBusinessStatus(ctx context.Context, params dto.ReviewBusinessParams) error {
	if params.NewStatus != "verified" && params.NewStatus != "rejected" {
		return apperr.New(http.StatusBadRequest, "invalid verification status, must be 'verified' or 'rejected'")
	}

	if params.NewStatus == "rejected" && (params.Comment == nil || *params.Comment == "") {
		return apperr.New(http.StatusBadRequest, "feedback comment is required when rejecting a business")
	}

	currentStatus, ownerID, err := s.businessService.GetBusinessStatusAndOwner(ctx, params.OrgUnitID)
	if err != nil {
		if errors.Is(err, apperr.ErrBusinessNotFound) {
			return apperr.Wrap(http.StatusNotFound, "business establishment not found", err)
		}
		return apperr.Wrap(http.StatusInternalServerError, "failed to check current business status", err)
	}

	if currentStatus != "pending" {
		return apperr.New(
			http.StatusConflict,
			fmt.Sprintf("cannot review business: it has already been reviewed and has status '%s'", currentStatus),
		)
	}

	err = s.businessService.ReviewBusiness(ctx, params)
	if err != nil {
		return apperr.Wrap(http.StatusInternalServerError, "failed to complete business verification review", err)
	}

	var eventType string
	metadata := make(map[string]any)

	if params.NewStatus == "verified" {
		eventType = outbox.EventTypeBusinessVerified
		metadata["message"] = "Congratulations! Your business verification request has been successfully approved."
	} else {
		eventType = outbox.EventTypeBusinessRejected
		metadata["message"] = fmt.Sprintf("We regret to inform you that your business verification request has been rejected. Reason: %s", *params.Comment)
	}

	// A business can be reviewed multiple times over its lifecycle (e.g. rejected,
	// resubmitted, then reviewed again), and each review is a distinct notification.
	// Mix the review timestamp into the event id so repeated reviews don't collide on
	// the notification_id UNIQUE constraint, while a single review stays idempotent
	// across relay/consumer delivery retries (the id is stored once with the row).
	now := time.Now().UTC()
	entityID := fmt.Sprintf("%d", params.OrgUnitID)
	event := outbox.Event{
		EventType: eventType,
		Payload: outbox.Message{
			EventID:     outbox.NewEventID(eventType, ownerID, params.AdminID, "business", entityID, fmt.Sprintf("%d", now.UnixNano())),
			EventType:   eventType,
			RecipientID: ownerID,
			ActorID:     params.AdminID,
			EntityType:  "business",
			EntityID:    entityID,
			Metadata:    metadata,
			CreatedAt:   now,
		},
		SourceService: outbox.DefaultSourceService,
	}

	if err := s.outboxWriter.Enqueue(ctx, event); err != nil {
		logger.ErrorKV(ctx, "failed to enqueue business review notification", "org_unit_id", params.OrgUnitID, "error", err.Error())
	}
	return nil
}
