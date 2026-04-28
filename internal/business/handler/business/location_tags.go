package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

type locationTagResponse struct {
	ID   int    `json:"id" example:"1"`
	Name string `json:"name" example:"Breakfast"`
	Slug string `json:"slug" example:"breakfast"`
}

// listLocationTags returns all available location tags.
//
//	@Summary		List location tags
//	@Description	Returns reference list of available location tags.
//	@Tags			locations
//	@Produce		json
//	@Success		200	{array}		locationTagResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/business/location-tags [get]
func (h *handler) listLocationTags(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	tags, err := h.service.ListLocationTags(ctx)
	if err != nil {
		log.Error("failed to list location tags", "error", err)
		c.Error(err)
		return
	}

	resp := make([]locationTagResponse, 0, len(tags))
	for _, t := range tags {
		resp = append(resp, locationTagResponse{
			ID:   t.ID,
			Name: t.Name,
			Slug: t.Slug,
		})
	}

	c.JSON(http.StatusOK, resp)
}
