package business

import (
	"net/http"
	"time"

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

type postResponse struct {
	ID        int64     `json:"id"`
	OrgID     int       `json:"org_id"`
	Content   string    `json:"content"`
	PostType  string    `json:"post_type"`
	CreatedAt time.Time `json:"created_at"`
}

type recommendPostsResponse struct {
	Items []postResponse `json:"items"`
	Total int            `json:"total"`
}

// recommendPosts returns posts recommended to user based on likes and tags
// @Summary Get posts by user behavior
// @Description Returns recommended posts by user behavior using weighted tag quotas (5-3-2-1-1).
// @Tags feed
// @Produce json
// @Param          lat     query       float64 true    "User latitude"
// @Param          lon     query       float64 true    "User longitude"
// @Param          skip    query       int     false   "Number of items to skip (default: 0)"
// @Param          limit   query       int     false   "Items per page (default: 10, max: 100)"
// @Success 200 {object} recommendPostsResponse
// @Security BearerAuth
// @Router /business/posts/recommend [get]
func (h *handler) recommendPosts(c *gin.Context) {
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

	log.Info("recommend posts", "user_id", userID, "skip", req.Skip, "limit", req.Limit)

	postsResult, err := h.service.RecommendPosts(ctx, userID, req.Lat, req.Lon, req.Skip, req.Limit)
	if err != nil {
		log.Error("failed to get recommended posts", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to recommend posts"})
		return
	}

	items := make([]postResponse, 0, len(postsResult.Items))
	for _, p := range postsResult.Items {
		items = append(items, postResponse{
			ID:        p.ID,
			OrgID:     p.OrgID,
			Content:   p.Content,
			PostType:  p.PostType,
			CreatedAt: p.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, recommendPostsResponse{
		Items: items,
		Total: postsResult.Total,
	})
}
