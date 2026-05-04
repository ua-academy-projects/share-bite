package business

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/internal/business/error/code"
)

// getOrgUnit godoc
// @Summary      Get business organization
// @Description  Returns a business BRAND profile by id (public access)
// @Tags         business
// @Produce      json
// @Param        id  path      int                 true  "Organization ID"
// @Success      200 {object}  dto.UpdateOrgResponse
// @Failure      400 {object}  errorResponse  "Validation error"
// @Failure      401 {object}  errorResponse  "Unauthorized"
// @Failure      403 {object}  errorResponse  "Forbidden"
// @Failure      404 {object}  errorResponse  "Not found"
// @Failure      500 {object}  errorResponse  "Internal server error"
// @Router       /business/{id} [get]
func (h *handler) getOrgUnit(c *gin.Context) {
	reqURI := new(getRequest)
	if err := c.ShouldBindUri(reqURI); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	org, err := h.service.Get(c.Request.Context(), reqURI.ID)
	if err != nil {
		var appErr *apperror.Error
		if errors.As(err, &appErr) {
			switch appErr.Code {
			case code.BadRequest:
				c.JSON(http.StatusBadRequest, gin.H{"error": appErr.Error()})
				return
			case code.NotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": appErr.Error()})
				return
			case code.Forbidden:
				c.JSON(http.StatusForbidden, gin.H{"error": appErr.Error()})
				return
			case code.Unauthorized:
				c.JSON(http.StatusUnauthorized, gin.H{"error": appErr.Error()})
				return
			}
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if org.ProfileType != entity.ProfileTypeBrand {
		c.JSON(http.StatusBadRequest, gin.H{"error": "target org unit is not a brand"})
		return
	}

	c.JSON(http.StatusOK, dto.UpdateOrgResponse{
		Id:          org.Id,
		Name:        org.Name,
		Avatar:      org.Avatar,
		Banner:      org.Banner,
		Description: org.Description,
	})
}
