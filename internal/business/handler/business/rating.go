package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

type ratingRequest struct {
	ID int `uri:"id" binding:"required"`
}

type ratingResponse struct {
	Rating float32 `json:"rating" example:"3.45"`
}

// rating returns a rating of the venue
//
//	@Summary		Get rating by location ID
//	@Description	Returns a rating that belongs to a location.
//	@Tags			locations
//	@Produce		json
//	@Param			id	path		int	true	"Location ID"
//	@Success		200	{object}	ratingResponse
//	@Failure		400	{object}	errorResponse
//	@Failure		404	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/business/org-units/{id}/rating [get]
func (h *handler) rating(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	req := new(ratingRequest)
	if err := c.ShouldBindUri(req); err != nil || req.ID < 1 {
		c.Error(apperror.BadRequest("invalid location id"))
		return
	}

	log.Info("get location rating", "id", req.ID)

	rating, err := h.service.Rating(ctx, req.ID)
	if err != nil {
		log.Error("failed to get location rating", "id", req.ID, "error", err)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, ratingResponse{
		Rating: rating,
	})
}
