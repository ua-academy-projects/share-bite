package business

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	repo "github.com/ua-academy-projects/share-bite/internal/business/repository/business"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

type reserveBoxRequest struct {
	BoxID int64 `uri:"boxID" binding:"required"`
}

type reserveBoxResponse struct {
	Image         string          `json:"image"`
	FullPrice     decimal.Decimal `json:"price_full"`
	DiscountPrice decimal.Decimal `json:"price_discount"`
	BoxCode       string          `json:"box_code"`
}

// reserveBox updates boxItem to be reserved by a user.
//
//	@Summary		Reserve an available box
//	@Description	Reserves an available box item for the authenticated user and returns brief box info.
//	@Tags			boxes
//	@Produce		json
//	@Param			boxID	path		int	true	"Box ID"
//	@Success		200		{object}	reserveBoxResponse
//	@Failure		400		{object}	errorResponse
//	@Failure		401		{object}	errorResponse
//	@Failure		404		{object}	errorResponse
//	@Failure		409		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Security		BearerAuth
//	@Router			/business/boxes/{boxID}/reserve [patch]
func (h *handler) reserveBox(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req reserveBoxRequest
	if err := c.ShouldBindUri(&req); err != nil || req.BoxID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid box id"})
		return
	}

	log.Info("reserve box", "boxID", req.BoxID, "userID", userID)

	resp, err := h.service.ReserveBox(ctx, userID, req.BoxID)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "box not found"})
		case errors.Is(err, repo.ErrNoAvailableItems):
			c.JSON(http.StatusConflict, gin.H{"error": "no available box items"})
		default:
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(http.StatusOK, reserveBoxResponse{
		Image:         resp.Image,
		FullPrice:     resp.FullPrice,
		DiscountPrice: resp.DiscountPrice,
		BoxCode:       resp.BoxCode,
	})
}
