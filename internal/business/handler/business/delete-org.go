package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
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
	role, err := httpctx.GetUserRole(c)
	if err != nil {
		c.Error(apperror.Unauthorized("unauthorized"))
		return
	}

	if role != "business" {
		c.Error(apperror.Forbidden("only business accounts can delete organizations"))
		return
	}

	reqURI := new(getRequest)
	if err := c.ShouldBindUri(reqURI); err != nil {
		c.Error(apperror.BadRequest("invalid id"))
		return
	}

	orgAccountID, err := h.extractUserUUID(c)
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.service.DeleteOrg(c.Request.Context(), reqURI.ID, orgAccountID); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
