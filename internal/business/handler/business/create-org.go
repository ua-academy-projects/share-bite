package business

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
)

// createOrgUnit godoc
// @Summary      Create new business organization
// @Description  Creates a BRAND (top-level)
// @Tags         business
// @Accept       json
// @Produce      json
// @Param        request body      dto.CreateOrgRequest  true  "Brand initialization data"
// @Success      201     {object}  map[string]int    "Returns created organization ID"
// @Failure      400     {object}  errorResponse     "Validation error"
// @Failure      401     {object}  errorResponse     "Unauthorized (missing token)"
// @Failure      403     {object}  errorResponse     "Forbidden"
// @Failure      409     {object}  errorResponse     "Conflict"
// @Failure      500     {object}  errorResponse     "Internal server error"
// @Router       /business/ [post]
// @Security     BearerAuth
func (h *handler) createOrgUnit(c *gin.Context) {
	var req dto.CreateOrgRequest

	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid contract: " + err.Error()})
		return
	}

	userUUID, ok := middleware.GetUserUUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_uuid type in context"})
		return
	}

	orgEntity := entity.OrgUnit{
		OrgAccountId: userUUID,
		ProfileType:  entity.ProfileTypeBrand,
		Name:         req.Name,
		Description:  req.Description,
		Avatar:       req.Avatar,
		Banner:       req.Banner,
	}

	id, err := h.service.Create(c.Request.Context(), orgEntity)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}
