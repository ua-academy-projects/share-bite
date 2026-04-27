package business

import (
	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	"net/http"
)

// ListNearbyVenues returns nearby locations (venues).
//
//	@Summary		List nearby locations
//	@Description	Returns a list of nearby venues sorted by distance using coordinates.
//	@Tags			locations
//	@Accept			json
//	@Produce		json
//	@Param			lat		query		float64	true	"Latitude"
//	@Param			lon		query		float64	true	"Longitude"
//	@Param			skip	query		int		false	"Number of items to skip (default: 0)"
//	@Param			limit	query		int		false	"Items per page (default: 10, max: 100)"
//	@Success		200		{object}	dto.NearbyVenuesListResponse
//	@Failure		400		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/business/locations/nearby [get]
func (h *handler) ListNearbyVenues(c *gin.Context) {
	var req dto.ListNearbyVenuesInput

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Skip < 0 {
		req.Skip = 0
	}

	ctx := c.Request.Context()

	res, err := h.service.ListNearbyVenues(ctx, req.Lat, req.Lon, req.Skip, req.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	items := make([]dto.NearbyVenueItem, 0, len(res.Items))
	for _, u := range res.Items {
		items = append(items, dto.NearbyVenueItem{
			ID:       u.OrgUnit.Id,
			Name:     u.OrgUnit.Name,
			Avatar:   u.OrgUnit.Avatar,
			Distance: u.Distance * 1.60934,
		})
	}

	c.JSON(http.StatusOK, dto.ListNearbyVenuesOutput{
		Items: items,
		Total: res.Total,
	})
}
