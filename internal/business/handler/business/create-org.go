package business

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
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
		c.Error(apperror.BadRequest("invalid contract"))
		return
	}
	if req.Name == "" {
		c.Error(apperror.BadRequest("name is required"))
		return
	}

	userUUID, err := h.extractUserUUID(c)
	if err != nil {
		c.Error(err)
		return
	}

	orgEntity := entity.OrgUnit{
		OrgAccountId: userUUID,
		ProfileType:  entity.ProfileTypeBrand,
		Name:         req.Name,
		Description:  req.Description,
		Avatar:       req.Avatar,
		Banner:       req.Banner,
		Status:       entity.OrgStatusPending,
	}

	id, err := h.service.Create(c.Request.Context(), orgEntity)
	if err != nil {
		_ = c.Error(err)
		return
	}
	h.metrics.RecordBusinessRegistered()
	c.JSON(http.StatusCreated, gin.H{"id": id})
}
