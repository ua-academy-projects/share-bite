package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

type listRequest struct {
	BrandID int `uri:"id" binding:"required"`
	Skip    int `form:"skip"`
	Limit   int `form:"limit"`
}

type listItem struct {
	ID          int      `json:"id" example:"42"`
	Name        string   `json:"name" example:"ShareBite Downtown"`
	Avatar      *string  `json:"avatar" example:"https://cdn.example.com/avatar.png"`
	Description *string  `json:"description" example:"A cozy place in the city center."`
	Latitude    *float32 `json:"latitude" example:"50.4501"`
	Longitude   *float32 `json:"longitude" example:"30.5234"`
	Tags        []string `json:"tags"`
}

type listResponse struct {
	Items []listItem `json:"items"`
	Total int        `json:"total" example:"42"`
}

// list returns locations (venues) belonging to a brand.
//
//	@Summary		List brand locations
//	@Description	Returns a paginated list of venues/locations that belong to the given brand.
//	@Tags			locations
//	@Produce		json
//	@Param			id		path		int	true	"Brand ID"
//	@Param			skip	query		int	false	"Number of items to skip (default: 0)"
//	@Param			limit	query		int	false	"Items per page (default: 10, max: 100)"
//	@Success		200		{object}	listResponse
//	@Failure		400		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/business/org-units/{id}/locations [get]
func (h *handler) list(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	req := new(listRequest)
	if err := c.ShouldBindUri(req); err != nil {
		c.Error(apperror.BadRequest("invalid brand id"))
		return
	}
	c.ShouldBindQuery(req)

	if req.Skip < 0 {
		req.Skip = 0
	}
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.Limit < 0 {
		req.Limit = 1
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	log.Info("list locations", "brandId", req.BrandID, "skip", req.Skip, "limit", req.Limit)

	result, err := h.service.List(ctx, req.BrandID, req.Skip, req.Limit)
	if err != nil {
		log.Error("failed to list locations", "brandId", req.BrandID, "error", err)
		c.Error(err)
		return
	}

	items := make([]listItem, 0, len(result.Items))
	for _, u := range result.Items {
		items = append(items, listItem{
			ID:          u.Id,
			Name:        u.Name,
			Avatar:      u.Avatar,
			Description: u.Description,
			Latitude:    u.Latitude,
			Longitude:   u.Longitude,
			Tags:        normalizeTags(u.Tags),
		})
	}

	c.JSON(http.StatusOK, listResponse{
		Items: items,
		Total: result.Total,
	})
}
