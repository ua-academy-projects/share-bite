package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

type updateVenueHoursURI struct {
	ID int `uri:"id" binding:"required,min=1"`
}

// updateVenueHours updates business venue working hours.
//
//	@Summary		Update venue hours
//	@Description	Updates weekly working hours for a business venue.
//	@Tags			locations
//	@Accept			json
//	@Produce		json
//	@Param			id				path		int						true	"Location ID"
//	@Param			Authorization	header		string					true	"Bearer access token"
//	@Param			request			body		dto.UpdateVenueHoursInput	true	"Venue hours payload"
//	@Success		200				{object}	dto.UpdateVenueHoursOutput
//	@Failure		400				{object}	errorResponse
//	@Failure		401				{object}	errorResponse
//	@Failure		403				{object}	errorResponse
//	@Failure		404				{object}	errorResponse
//	@Failure		500				{object}	errorResponse
//	@Router			/business/locations/{id}/hours [patch]
func (h *handler) updateVenueHours(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.Error(apperror.Unauthorized("unauthorized"))
		return
	}

	var uri updateVenueHoursURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.Error(apperror.BadRequest("invalid location id"))
		return
	}

	var req dto.UpdateVenueHoursInput
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorKV(ctx, "failed to bind venue hours request", "locationId", uri.ID, "error", err)
		c.Error(apperror.BadRequest("invalid request"))
		return
	}

	out, err := h.service.UpdateVenueHours(ctx, uri.ID, userID, req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, out)
}
