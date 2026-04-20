package business

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	repo "github.com/ua-academy-projects/share-bite/internal/business/repository/business"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
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
		c.JSON(http.StatusForbidden, gin.H{"error": "only business accounts can update organizations"})
		return
	}

	reqURI := new(getRequest)
	if err := c.ShouldBindUri(reqURI); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	val, exists := c.Get(middleware.CtxUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := val.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user_id type in context"})
		return
	}

	orgAccountID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id format"})
		return
	}

	var req dto.UpdateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	updated, err := h.service.UpdateOrg(c.Request.Context(), reqURI.ID, orgAccountID, entity.UpdateOrgUnitInput{
		Name:        req.Name,
		Avatar:      req.Avatar,
		Banner:      req.Banner,
		Description: req.Description,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
	})
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		default:
			_ = c.Error(err)
		}
		return
	}

	c.JSON(http.StatusOK, dto.UpdateOrgResponse{
		Id:          updated.Id,
		Name:        updated.Name,
		Avatar:      updated.Avatar,
		Banner:      updated.Banner,
		Description: updated.Description,
		Latitude:    updated.Latitude,
		Longitude:   updated.Longitude,
	})
}
