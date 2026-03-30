package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

type listRequest struct {
	BrandId int `uri:"id" binding:"required"`
	Page    int `form:"page"`
	Limit   int `form:"limit"`
}

type listItem struct {
	Id          int      `json:"id" example:"42"`
	Name        string   `json:"name" example:"ShareBite Downtown"`
	Avatar      *string  `json:"avatar" example:"https://cdn.example.com/avatar.png"`
	Description *string  `json:"description" example:"A cozy place in the city center."`
	Latitude    *float32 `json:"latitude" example:"50.4501"`
	Longitude   *float32 `json:"longitude" example:"30.5234"`
}

type listResponse struct {
	Items []listItem `json:"items"`
}

// list returns locations (venues) belonging to a brand.
//
//	@Summary		List brand locations
//	@Description	Returns a paginated list of venues/locations that belong to the given brand.
//	@Tags			locations
//	@Produce		json
//	@Param			id		path		int	true	"Brand ID"
//	@Param			page	query		int	false	"Page number (default: 1)"
//	@Param			limit	query		int	false	"Items per page (default: 10, max: 100)"
//	@Success		200		{object}	listResponse
//	@Failure		400		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/business/{id}/locations [get]
func (h *handler) list(c *gin.Context) {
	req := new(listRequest)
	if err := c.ShouldBindUri(req); err != nil {
		c.Error(apperror.BadRequest("invalid brand id"))
		return
	}
	c.ShouldBindQuery(req)

	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 10
	}

	ctx := c.Request.Context()
	logger.InfoKV(ctx, "list locations", "brandId", req.BrandId, "page", req.Page, "limit", req.Limit)

	orgUnits, err := h.service.List(ctx, req.BrandId, req.Page, req.Limit)
	if err != nil {
		logger.ErrorKV(ctx, "failed to list locations", "brandId", req.BrandId, "error", err)
		c.Error(err)
		return
	}

	items := make([]listItem, 0, len(orgUnits))
	for _, u := range orgUnits {
		items = append(items, listItem{
			Id:          u.Id,
			Name:        u.Name,
			Avatar:      u.Avatar,
			Description: u.Description,
			Latitude:    u.Latitude,
			Longitude:   u.Longitude,
		})
	}

	c.JSON(http.StatusOK, listResponse{
		Items: items,
	})
}
