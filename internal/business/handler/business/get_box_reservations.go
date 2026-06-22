package business

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	repo "github.com/ua-academy-projects/share-bite/internal/business/repository/business"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
)

type getBoxReservationsRequest struct {
	Skip  int `form:"skip"`
	Limit int `form:"limit"`
}

// GetBoxReservations returns a paginated list of reservations (box items) for a specific food box.
//
//	@Summary		Get food box reservations
//	@Description	Returns a paginated list of all box items (reservations) for a specific food box
//	@Tags			boxes
//	@Produce		json
//	@Param			id		path		int							true	"Box ID"
//	@Param			skip	query		int							false	"Number of items to skip (default: 0)"
//	@Param			limit	query		int							false	"Items per page (default: 10, max: 100)"
//	@Success		200		{object}	dto.BoxReservationsResponse
//	@Failure		400		{object}	errorResponse
//	@Failure		401		{object}	errorResponse
//	@Failure		403		{object}	errorResponse
//	@Failure		404		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Security		BearerAuth
//	@Router			/business/boxes/{id}/reservations [get]
func (h *handler) GetBoxReservations(c *gin.Context) {
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

	var req getBoxReservationsRequest
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

	result, err := h.service.GetBoxReservations(ctx, boxID, userID, req.Skip, req.Limit)
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

	items := make([]dto.BoxReservationItem, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, dto.BoxReservationItem{
			BoxCode:          item.BoxCode,
			ReservedByUserID: item.ReservedByUserID,
		})
	}

	c.JSON(http.StatusOK, dto.BoxReservationsResponse{
		Items: items,
		Total: result.Total,
	})
}
