package business

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
)

// createOrgUnit godoc
// @Summary      Create new business organization
// @Description  Creates a BRAND (top-level) or VENUE (location under a brand). VENUEs must specify parent_id.
// @Tags         business
// @Accept       json
// @Produce      json
// @Param        request body      dto.CreateOrgRequest  true  "Business data (parent_id required for VENUE)"
// @Success      201     {object}  map[string]int    "Returns created organization ID"
// @Failure      400     {object}  errorResponse     "Validation error"
// @Failure      401     {object}  errorResponse     "Unauthorized (missing token)"
// @Router       /business/ [post]
// @Security     BearerAuth
func (h *handler) createOrgUnit(c *gin.Context) {
	var req dto.CreateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profileType := strings.ToUpper(strings.TrimSpace(req.ProfileType))
	if profileType != entity.ProfileTypeBrand && profileType != entity.ProfileTypeVenue {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid profile type"})
		return
	}

	val, exists := c.Get(middleware.CtxUserID)

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found"})
		return
	}


	userID, ok := val.(string)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user_id type in context"})
		return
	}
	parsedUUID, err := uuid.Parse(userID)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id format"})
		return
	}

	orgEntity := entity.OrgUnit{
		OrgAccountId: parsedUUID,
		ProfileType:  profileType,
		ParentId:     req.ParentID,
		Name:         req.Name,
		Description:  req.Description,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
	}

	id, err := h.service.Create(c.Request.Context(), orgEntity)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}
