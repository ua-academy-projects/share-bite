package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
)

// updateOrgUnit godoc
// @Summary      Update business organization
// @Description  Updates mutable business profile fields for organization owned by the authenticated business account.
// @Tags         business
// @Accept       json
// @Produce      json
// @Param        id      path      int               true  "Organization ID"
// @Param        request body      dto.UpdateOrgRequest  true  "Fields to update (partial allowed)"
// @Success      200     {object}  dto.UpdateOrgResponse
// @Failure      400     {object}  errorResponse     "Validation error"
// @Failure      401     {object}  errorResponse     "Unauthorized"
// @Failure      403     {object}  errorResponse     "Forbidden"
// @Failure      404     {object}  errorResponse     "Not found"
// @Failure      500     {object}  errorResponse     "Internal server error"
// @Router       /business/{id} [put]
// @Router       /business/{id} [patch]
// @Security     BearerAuth
func (h *handler) updateOrgUnit(c *gin.Context) {
	role, err := httpctx.GetUserRole(c)
	if err != nil {
		c.Error(apperror.Unauthorized("unauthorized"))
		return
	}

	if role != RoleBusiness {
		c.Error(apperror.Forbidden("only business accounts can update organizations"))
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

	var req dto.UpdateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.BadRequest("invalid request"))
		return
	}

	updated, err := h.service.UpdateOrg(c.Request.Context(), reqURI.ID, orgAccountID, entity.UpdateOrgUnitInput{
		Name:        req.Name,
		Avatar:      req.Avatar,
		Banner:      req.Banner,
		Description: req.Description,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.UpdateOrgResponse{
		Id:          updated.Id,
		Name:        updated.Name,
		Avatar:      updated.Avatar,
		Banner:      updated.Banner,
		Description: updated.Description,
		Status:      string(updated.Status),
	})
}
