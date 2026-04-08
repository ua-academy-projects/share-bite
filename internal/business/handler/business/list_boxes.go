package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
)

func (h *handler) ListNearbyBoxes (c *gin.Context) {
	var req dto.GetNearbyBoxesReq

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	res, err := h.service.ListNearbyBoxes(ctx, req.Skip, req.Limit, req.Lat, req.Lon, req.CategoryID)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items := make([]dto.NearbyBoxesResp, 0, len(res.Items))
	for _, u := range res.Items{
		items = append(items, dto.NearbyBoxesResp{
			Id: u.Box.Id,
			VenueId: u.Box.VenueId,
			CategoryID: u.Box.CategoryID,
			Image: u.Box.Image,
			FullPrice: u.Box.FullPrice,
			DiscountPrice: u.Box.DiscountPrice,
			CreatedAt: u.Box.CreatedAt,
			ExpiresAt: u.Box.ExpiresAt,
			Distance: u.Distance,
		})
	}

	c.JSON(http.StatusOK, dto.ListResponse{
		Items: items,
		Total: res.Total,
	})
}