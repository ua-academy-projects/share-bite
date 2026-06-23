package business

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	repo "github.com/ua-academy-projects/share-bite/internal/business/repository/business"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
)

type updateBoxRequest struct {
	CategoryID    *int    `json:"category_id"`
	FullPrice     *string `json:"price_full"`
	DiscountPrice *string `json:"price_discount"`
	ExpiresAt     *string `json:"expires_at"`
}

// UpdateBox updates a food box with validation.
//
//	@Summary		Update food box
//	@Description	Updates a food box with validation (price constraints, expires_at must be in future)
//	@Tags			boxes
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"Box ID"
//	@Param			input	body		updateBoxRequest		true	"Updated box data"
//	@Success		200		{object}	dto.BoxResponse
//	@Failure		400		{object}	errorResponse
//	@Failure		401		{object}	errorResponse
//	@Failure		403		{object}	errorResponse
//	@Failure		404		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Security		BearerAuth
//	@Router			/business/boxes/{id} [put]
func (h *handler) UpdateBox(c *gin.Context) {
	boxID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || boxID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid box id"})
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req updateBoxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	input := entity.BoxUpdateInput{
		CategoryID: req.CategoryID,
	}

	if req.FullPrice != nil {
		price, err := decimal.NewFromString(*req.FullPrice)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid price_full format"})
			return
		}
		input.FullPrice = &price
	}

	if req.DiscountPrice != nil {
		price, err := decimal.NewFromString(*req.DiscountPrice)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid price_discount format"})
			return
		}
		input.DiscountPrice = &price
	}

	if req.ExpiresAt != nil {
		expiresAt, err := time.Parse(time.RFC3339, *req.ExpiresAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expires_at format (use RFC3339)"})
			return
		}
		input.ExpiresAt = &expiresAt
	}

	updatedBox, err := h.service.UpdateBox(c.Request.Context(), boxID, userID, input)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "box not found"})
		case errors.Is(err, repo.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		default:
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, dto.BoxResponse{
		ID:            updatedBox.ID,
		VenueID:       updatedBox.VenueID,
		CategoryID:    updatedBox.CategoryID,
		Image:         updatedBox.Image,
		FullPrice:     updatedBox.FullPrice,
		DiscountPrice: updatedBox.DiscountPrice,
		CreatedAt:     updatedBox.CreatedAt,
		ExpiresAt:     updatedBox.ExpiresAt,
	})
}
