package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
)

// ListNearbyBoxes returns a paginated list of nearby boxes sorted by distance.
//
//		@Summary        List nearby boxes
//		@Description    Returns a paginated list of available boxes sorted by distance from the provided coordinates (lat/lon).
//		@Tags           boxes
//		@Produce        json
//		@Param          lat         query       float64 true    "User latitude"
//		@Param          lon         query       float64 true    "User longitude"
//		@Param          skip        query       int     false   "Number of items to skip (default: 0)"
//		@Param          limit       query       int     false  	"Items per page (default: 10, max: 100)"
// 		@Param			org_id		query		int 	false	"Optional Organisation ID to filter by"
//		@Param          category_id query       int     false   "Optional Category ID to filter by"
//		@Success        200         {object}    dto.ListResponse
//		@Failure        400         {object}    errorResponse
//		@Failure        500         {object}    errorResponse
//	 	@Router 		/business/nearby-boxes  [get]
func (h *handler) ListNearbyBoxes(c *gin.Context) {
	var req dto.GetNearbyBoxesReq

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	ctx := c.Request.Context()

	res, err := h.service.ListNearbyBoxes(ctx, req.Skip, req.Limit, req.Lat, req.Lon, req.CategoryID, req.OrgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	items := make([]dto.NearbyBoxesResp, 0, len(res.Items))
	for _, u := range res.Items {
		items = append(items, dto.NearbyBoxesResp{
			ID:                 u.Box.ID,
			VenueID:            u.Box.VenueID,
			CategoryID:         u.Box.CategoryID,
			Image:              u.Box.Image,
			FullPrice:          u.Box.FullPrice,
			DiscountPrice:      u.Box.DiscountPrice,
			CreatedAt:          u.Box.CreatedAt,
			ExpiresAt:          u.Box.ExpiresAt,
			AvailabilityStatus: string(u.AvailabilityStatus()),
			Distance:           u.Distance,
		})
	}

	c.JSON(http.StatusOK, dto.ListResponse{
		Items: items,
		Total: res.Total,
	})
}
