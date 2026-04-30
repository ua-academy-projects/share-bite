package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
)

// deleteOrgUnit godoc
// @Summary      Delete business organization
// @Description  Deletes a business organization owned by the authenticated business account.
// @Tags         business
// @Produce      json
// @Param        id  path      int            true  "Organization ID"
// @Success      204 {object}  nil
// @Failure      400 {object}  errorResponse  "Validation error"
// @Failure      401 {object}  errorResponse  "Unauthorized"
// @Failure      403 {object}  errorResponse  "Forbidden"
// @Failure      404 {object}  errorResponse  "Not found"
// @Failure      500 {object}  errorResponse  "Internal server error"
// @Router       /business/{id} [delete]
// @Security     BearerAuth
func (h *handler) deleteOrgUnit(c *gin.Context) {
	roleVal, exists := c.Get(middleware.CtxUserRole)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	role, ok := roleVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user_role type in context"})
		return
	}

	if role != "business" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only business accounts can delete organizations"})
		return
	}

	reqURI := new(getRequest)
	if err := c.ShouldBindUri(reqURI); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	orgAccountID, ok := middleware.GetUserUUID(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user_uuid type in context"})
		return
	}

	if err := h.service.DeleteOrg(c.Request.Context(), reqURI.ID, orgAccountID); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
