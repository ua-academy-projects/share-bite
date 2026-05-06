package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

type getVenuesByIDsRequest struct {
	IDs []int `json:"ids" binding:"required,min=1" example:"3"`
}

type venueResponse struct {
	ID          int     `json:"id" example:"3"`
	Name        string  `json:"name" example:"ShareBite Downtown"`
	Description *string `json:"description" example:"A cozy place in the city center."`
	Avatar      *string `json:"avatar" example:"https://cdn.example.com/avatar.png"`
	Banner      *string `json:"banner" example:"https://cdn.example.com/banner.png"`
}

// getVenuesByIDs returns venues by their IDs.
//
//	@Summary		Get venues by IDs
//	@Description	Returns a list of venues (org units with a parent) matching the provided IDs.
//	@Tags			venues
//	@Accept			json
//	@Produce		json
//	@Param			request	body		getVenuesByIDsRequest	true	"Array of venue IDs"
//	@Success		200		{array}		venueResponse
//	@Failure		400		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/business/org-units/venues [post]
func (h *handler) getVenuesByIDs(c *gin.Context) {
	var req getVenuesByIDsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.BadRequest("invalid request body"))
		return
	}

	ctx := c.Request.Context()
	logger.InfoKV(ctx, "get venues by ids", "count", len(req.IDs))

	venues, err := h.service.GetVenuesByIDs(ctx, req.IDs)
	if err != nil {
		logger.ErrorKV(ctx, "failed to get venues by ids", "error", err)
		c.Error(err)
		return
	}

	items := make([]venueResponse, 0, len(venues))
	for _, v := range venues {
		items = append(items, venueResponse{
			ID:          v.Id,
			Name:        v.Name,
			Description: v.Description,
			Avatar:      v.Avatar,
			Banner:      v.Banner,
		})
	}

	c.JSON(http.StatusOK, items)
}
