package business

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	repo "github.com/ua-academy-projects/share-bite/internal/business/repository/business"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
)

type createBoxRequest struct {
	VenueID       int             `json:"venue_id" binding:"required"`
	CategoryID    *int            `json:"category_id"`
	Image         string          `json:"image" binding:"required"`
	FullPrice     decimal.Decimal `json:"price_full" binding:"required"`
	DiscountPrice decimal.Decimal `json:"price_discount"`
	ExpiresAt     time.Time       `json:"expires_at" binding:"required"`
	Quantity      int             `json:"quantity" binding:"required,min=1,max=1000"`
}

// CreateBox creates a limited box for a venue.
//
//	@Summary		Create box
//	@Description	Creates a box and its limited box items for a specific venue if the user has permission
//	@Tags			boxes
//	@Accept			json
//	@Produce		json
//	@Param			input	body		createBoxRequest	true	"Box data"
//	@Success		201		{object}	CreateBoxResponse
//	@Failure		400		{object}	errorResponse
//	@Failure		401		{object}	errorResponse
//	@Failure		403		{object}	errorResponse
//	@Failure		404		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Security		BearerAuth
//	@Router			/business/boxes [post]
func (h *handler) CreateBox(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req createBoxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if req.VenueID <= 0 ||
		(req.CategoryID != nil && *req.CategoryID <= 0) ||
		req.FullPrice.LessThanOrEqual(decimal.Zero) ||
		req.DiscountPrice.LessThan(decimal.Zero) ||
		req.DiscountPrice.GreaterThan(req.FullPrice) ||
		len(req.Image) > 256 ||
		!req.ExpiresAt.After(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	box, err := h.service.CreateBox(c.Request.Context(), userID, dto.CreateBoxRequest{
		VenueID:       req.VenueID,
		CategoryID:    req.CategoryID,
		Image:         req.Image,
		FullPrice:     req.FullPrice,
		DiscountPrice: req.DiscountPrice,
		ExpiresAt:     req.ExpiresAt,
		Quantity:      req.Quantity,
	})
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		case errors.Is(err, repo.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "venue not found"})
		default:
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(http.StatusCreated, CreateBoxResponse{
		ID:      box.ID,
		Message: "box created",
	})
}
