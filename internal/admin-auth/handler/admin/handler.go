package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/dto"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/handler"
	adminsvc "github.com/ua-academy-projects/share-bite/internal/admin-auth/service/admin"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

type Handler struct {
	service adminsvc.Service
}

func NewHandler(service adminsvc.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetUsersList godoc
// @Summary      Get list of users
// @Description  Retrieves a paginated list of users for the admin panel with optional filtering.
// @Tags         Admin
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        limit   query     int     false  "Number of items to return (default 10)"
// @Param        offset  query     int     false  "Number of items to skip (default 0)"
// @Param        search_email  query     string  false  "Search by email (partial match)"
// @Param        role    query     string  false  "Filter by role slug (e.g., user, business, moderator)"
// @Param        status  query     string  false  "Filter by status (e.g., active, blocked)"
// @Param        sort_order query  string  false "Sort by created_at (asc/desc)"
// @Success      200     {object}  dto.PaginatedAdminUsersResponse  "Success. Returns paginated list of users."
// @Failure      400     {object}  handler.ErrorResponse            "Invalid query parameters."
// @Failure      401     {object}  handler.ErrorResponse            "Unauthorized access."
// @Failure      403     {object}  handler.ErrorResponse            "Forbidden. Admin or moderator role required."
// @Failure      500     {object}  handler.ErrorResponse            "Internal server error."
// @Router       /admin/users [get]
func (h *Handler) GetUsersList(c *gin.Context) {
	var query handler.UsersFilterQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, handler.ErrorResponse{Error: "Invalid query parameters."})
		return
	}
	limit := 10
	offset := 0

	if query.Limit != nil {
		limit = *query.Limit
	}

	if query.Offset != nil {
		offset = *query.Offset
	}

	filter := dto.AdminUserFilter{
		Limit:       limit,
		Offset:      offset,
		SearchEmail: query.Search,
		RoleSlug:    query.Role,
		Status:      query.Status,
		SortOrder:   query.SortOrder,
	}

	resp, err := h.service.GetUsersList(c.Request.Context(), filter)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetUserDetails godoc
// @Summary      Get user details
// @Description  Retrieves detailed profile information for a specific user, including business or guest profile data.
// @Tags         Admin
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "User ID"
// @Success      200     {object}  dto.FullUserDetails    "Success. Returns detailed user information."
// @Failure      401     {object}  handler.ErrorResponse  "Unauthorized access."
// @Failure      403     {object}  handler.ErrorResponse  "Forbidden. Admin or moderator role required."
// @Failure      404     {object}  handler.ErrorResponse  "User not found."
// @Failure      500     {object}  handler.ErrorResponse  "Internal server error."
// @Router       /admin/users/{id} [get]
func (h *Handler) GetUserDetails(c *gin.Context) {
	userID := c.Param("id")

	if err := uuid.Validate(userID); err != nil {
		c.JSON(http.StatusBadRequest, handler.ErrorResponse{Error: "invalid user id format"})
		return
	}

	user, err := h.service.GetUserDetails(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetPlatformStatistics godoc
// @Summary      Get platform statistics
// @Description  Retrieves all-time aggregated platform metrics across auth, guest, and business domains.
// @Tags         Admin
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200     {object}  dto.PlatformStatisticsResponse  "Success. Returns aggregated platform statistics."
// @Failure      401     {object}  handler.ErrorResponse           "Unauthorized access."
// @Failure      403     {object}  handler.ErrorResponse           "Forbidden. Admin or moderator role required."
// @Failure      500     {object}  handler.ErrorResponse           "Internal server error."
// @Router       /admin/statistics [get]
func (h *Handler) GetPlatformStatistics(c *gin.Context) {
	stats, err := h.service.GetPlatformStatistics(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, stats)
}

// ChangeUserRole godoc
// @Summary      Change user role
// @Description  Changes the role of a user and invalidates all their active sessions.
// @Tags         Admin
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path      string                     true  "User ID"
// @Param        request  body      handler.ChangeRoleRequest  true  "New role payload"
// @Success      200      {object}  handler.MessageResponse    "Success message."
// @Failure      400      {object}  handler.ErrorResponse      "Validation error or invalid role transition."
// @Failure      401      {object}  handler.ErrorResponse      "Unauthorized access."
// @Failure      403      {object}  handler.ErrorResponse      "Forbidden. Super admin role required."
// @Failure      404      {object}  handler.ErrorResponse      "User or role not found."
// @Failure      409      {object}  handler.ErrorResponse      "Business logic conflict (e.g., mixing business and personal accounts)."
// @Failure      500      {object}  handler.ErrorResponse      "Internal server error."
// @Router       /admin/users/{id}/role [patch]
func (h *Handler) ChangeUserRole(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.Param("id")

	if err := uuid.Validate(userID); err != nil {
		c.JSON(http.StatusBadRequest, handler.ErrorResponse{Error: "invalid user id format"})
		return
	}

	var req handler.ChangeRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, handler.ErrorResponse{Error: "Invalid request payload. 'role_slug' is required."})
		return
	}

	err := h.service.ChangeUserRole(ctx, userID, req.RoleSlug)
	if err != nil {
		_ = c.Error(err)
		return
	}
	logger.InfoKV(ctx, "user role successfully changed", "target_user_id", userID, "new_role", req.RoleSlug)

	c.JSON(http.StatusOK, handler.MessageResponse{Message: "User role has been successfully updated."})
}
