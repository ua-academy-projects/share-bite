package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

type recommendRequest struct {
	Skip  int     `form:"skip"`
	Limit int     `form:"limit"`
	Lat   float64 `form:"lat" binding:"required,latitude"`
	Lon   float64 `form:"lon" binding:"required,longitude"`
}

type recommendVenuesResponse struct {
	Items []venueResponse `json:"items"`
	Total int             `json:"total"`
}

// recommendVenues returns venues recommended to user based on likes and tags
// @Summary Get venues by user behavior
// @Description Returns recommended venues by user behavior using weighted tag quotas (5-3-2-1-1).
// @Tags locations
// @Produce json
// @Param          lat     query       float64 true    "User latitude"
// @Param          lon     query       float64 true    "User longitude"
// @Param			skip	query		int	false	"Number of items to skip (default: 0)"
// @Param			limit	query		int	false	"Items per page (default: 10, max: 100)"
// @Success 200 {object} recommendVenuesResponse
// @Security BearerAuth
// @Router /business/venues/recommend [get]
func (h *handler) recommendVenues(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	userID, ok := middleware.GetUserID(c)
	if !ok {
		log.Error(ctx, "unauthorized access attempt: user ID not found in gin context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	req := new(recommendRequest)
	if err := c.ShouldBindQuery(req); err != nil {
		c.Error(apperror.BadRequest("invalid query parameters"))
		return
	}

	if req.Skip < 0 {
		req.Skip = 0
	}
	if req.Limit == 0 {
		req.Limit = 24
	}
	if req.Limit < 1 {
		req.Limit = 1
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	log.Info("recommend venues", "user_id", userID, "skip", req.Skip, "limit", req.Limit)

	venues, err := h.service.RecommendVenues(ctx, userID, req.Lat, req.Lon, req.Skip, req.Limit)
	if err != nil {
		log.Error("failed to get recommended venues", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to recommend venues"})
		return
	}

	items := make([]venueResponse, 0, len(venues.Items))
	for _, v := range venues.Items {
		items = append(items, venueResponse{
			ID:          v.Id,
			Name:        v.Name,
			Description: v.Description,
			Avatar:      v.Avatar,
			Banner:      v.Banner,
		})
	}

	c.JSON(http.StatusOK, recommendVenuesResponse{
		Items: items,
		Total: venues.Total,
	})
}
