package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
)

type listBoxesByBusinessRequest struct {
	Skip  int `form:"skip"`
	Limit int `form:"limit"`
}

// ListBoxesByBusiness returns a paginated list of food boxes created by the authenticated business.
//
//	@Summary		List food boxes for authenticated business
//	@Description	Returns a paginated list of food boxes created by the authenticated business account
//	@Tags			boxes
//	@Produce		json
//	@Param			skip	query		int	false	"Number of items to skip (default: 0)"
//	@Param			limit	query		int	false	"Items per page (default: 10, max: 100)"
//	@Success		200		{object}	dto.BusinessBoxesListResponse
//	@Failure		400		{object}	errorResponse
//	@Failure		401		{object}	errorResponse
//	@Failure		403		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Security		BearerAuth
//	@Router			/business/boxes [get]
func (h *handler) ListBoxesByBusiness(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req listBoxesByBusinessRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

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

	ctx := c.Request.Context()

	result, err := h.service.ListBoxesByBusiness(ctx, userID, req.Skip, req.Limit)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	items := make([]dto.BoxResponse, 0, len(result.Items))
	for _, box := range result.Items {
		items = append(items, dto.BoxResponse{
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

	c.JSON(http.StatusOK, dto.BusinessBoxesListResponse{
		Items: items,
		Total: result.Total,
	})
}
