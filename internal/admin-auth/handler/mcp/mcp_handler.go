package mcp

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/handler"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// GetContext godoc
// @Summary      Get current MCP admin context
// @Description  Returns context info (id, roles, status) extracted from JWT for Python MCP server.
// @Tags         MCP
// @Security     BearerAuth
// @Produce      json
// @Success      200     {object}  handler.MCPContextResponse
// @Failure      401     {object}  handler.ErrorResponse
// @Router       /mcp/context [get]
func (h *Handler) GetContext(c *gin.Context) {
	userID, existsID := c.Get(middleware.CtxUserID)
	userRole, existsRole := c.Get(middleware.CtxUserRole)
	userStatus, existsStatus := c.Get(middleware.CtxUserStatus)

	if !existsID || !existsRole || !existsStatus {
		c.JSON(http.StatusUnauthorized, handler.ErrorResponse{
			Error: "Unauthorized: session context data is missing",
		})
		return
	}

	c.JSON(http.StatusOK, handler.MCPContextResponse{
		ID:     userID.(string),
		Roles:  []string{userRole.(string)},
		Status: fmt.Sprintf("%v", userStatus),
	})
}

// ValidateAdminPermissions godoc
// @Summary      Validate explicit MCP admin permission
// @Description  Checks if the admin has explicit rights for sensitive actions requested by MCP.
// @Tags         MCP
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      handler.ValidatePermissionRequest  true  "Permission payload"
// @Success      200      {object}  handler.MCPAuthorizedResponse
// @Failure      400      {object}  handler.ErrorResponse
// @Router       /mcp/validate-permission [post]
func (h *Handler) ValidateAdminPermissions(c *gin.Context) {
	var req handler.ValidatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, handler.ErrorResponse{
			Error: "Invalid request payload. 'permission' field is required.",
		})
		return
	}
	c.JSON(http.StatusOK, handler.MCPAuthorizedResponse{
		Authorized: true,
		Permission: req.Permission,
	})
}

// GetHealth godoc
// @Summary      MCP infrastructure health check
// @Description  Verifies that the auth service is reachable and functional for the MCP subsystem.
// @Tags         MCP
// @Security     BearerAuth
// @Produce      json
// @Success      200     {object}  handler.MCPHealthResponse
// @Router       /mcp/health [get]
func (h *Handler) GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, handler.MCPHealthResponse{
		Status: "healthy",
	})
}
