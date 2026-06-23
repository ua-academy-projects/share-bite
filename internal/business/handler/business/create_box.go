package business

import (
	"errors"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	biserr "github.com/ua-academy-projects/share-bite/internal/business/error"
	repo "github.com/ua-academy-projects/share-bite/internal/business/repository/business"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
)

type createBoxRequest struct {
	VenueID       int                   `form:"venue_id" binding:"required"`
	CategoryID    *int                  `form:"category_id"`
	Image         *multipart.FileHeader `form:"image" binding:"required"`
	FullPrice     decimal.Decimal       `form:"price_full" binding:"required"`
	DiscountPrice decimal.Decimal       `form:"price_discount"`
	ExpiresAt     time.Time             `form:"expires_at" binding:"required"`
	Quantity      int                   `form:"quantity" binding:"required,min=1,max=1000"`
}

type CreateBoxResponse struct {
	ID      int64  `json:"id"`
	Message string `json:"message"`
}

// CreateBox creates a limited box for a venue.
//
//	@Summary		Create box
//	@Description	Creates a box and its limited box items for a specific venue if the user has permission
//	@Tags			boxes
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			venue_id		formData	int		true	"Venue ID"
//	@Param			category_id		formData	int		false	"Category ID"
//	@Param			image			formData	file	true	"Box image"
//	@Param			price_full		formData	number	true	"Full price"
//	@Param			price_discount	formData	number	false	"Discount price"
//	@Param			expires_at		formData	string	true	"Expires at"
//	@Param			quantity		formData	int		true	"Quantity"
//	@Success		201				{object}	CreateBoxResponse
//	@Failure		400				{object}	errorResponse
//	@Failure		401				{object}	errorResponse
//	@Failure		403				{object}	errorResponse
//	@Failure		404				{object}	errorResponse
//	@Failure		500				{object}	errorResponse
//	@Security		BearerAuth
//	@Router			/business/boxes [post]
func (h *handler) CreateBox(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req createBoxRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if req.VenueID <= 0 ||
		(req.CategoryID != nil && *req.CategoryID <= 0) ||
		req.FullPrice.LessThanOrEqual(decimal.Zero) ||
		req.DiscountPrice.LessThan(decimal.Zero) ||
		req.DiscountPrice.GreaterThan(req.FullPrice) ||
		!req.ExpiresAt.After(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	box, err := h.service.CreateBox(c.Request.Context(), userID, dto.CreateBoxRequest{
		VenueID:       req.VenueID,
		CategoryID:    req.CategoryID,
		FullPrice:     req.FullPrice,
		DiscountPrice: req.DiscountPrice,
		ExpiresAt:     req.ExpiresAt,
		Quantity:      req.Quantity,
	}, req.Image)

	if err != nil {
		switch {
		case errors.Is(err, biserr.WrongFileExtErr), errors.Is(err, biserr.FileToLargeErr):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
	h.metrics.RecordBoxCreated()

	c.JSON(http.StatusCreated, CreateBoxResponse{
		ID:      box.ID,
		Message: "box created",
	})
}
