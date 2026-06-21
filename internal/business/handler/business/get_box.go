package business

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	repo "github.com/ua-academy-projects/share-bite/internal/business/repository/business"
)

// GetBox retrieves a single food box by ID.
//
//	@Summary		Get food box by ID
//	@Description	Retrieves a single food box with all details (prices, expiration, images, etc.)
//	@Tags			boxes
//	@Produce		json
//	@Param			id	path		int	true	"Food Box ID"
//	@Success		200			{object}	dto.BoxResponse
//	@Failure		400			{object}	errorResponse
//	@Failure		404			{object}	errorResponse
//	@Failure		500			{object}	errorResponse
//	@Router			/business/food-boxes/{id} [get]
func (h *handler) GetBox(c *gin.Context) {
	boxID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || boxID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid box id"})
		return
	}

	ctx := c.Request.Context()

	box, err := h.service.GetBox(ctx, boxID)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "box not found"})
		default:
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, dto.BoxResponse{
		ID:            box.ID,
		VenueID:       box.VenueID,
		CategoryID:    box.CategoryID,
		Image:         box.Image,
		FullPrice:     box.FullPrice,
		DiscountPrice: box.DiscountPrice,
		CreatedAt:     box.CreatedAt,
		ExpiresAt:     box.ExpiresAt,
	})
}
