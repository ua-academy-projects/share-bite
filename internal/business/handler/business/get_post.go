package business

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
)

type getPostsRequest struct {
	Skip   int    `form:"skip"`
	Limit  int    `form:"limit"`
	OrgIDs string `form:"org_id"` // Comma-separated list
}

type getPostsResponse struct {
	Items []dto.PostResponse `json:"items"`
	Total int                `json:"total"`
}

// GetPosts returns a paginated list of posts.
//
// @Summary      Get posts
// @Description  Returns paginated posts with images and organization info. Supports filtering by org_id (comma-separated).
// @Tags         posts
// @Produce      json
// @Param        skip    query     int     false  "Number of items to skip (default: 0)"
// @Param        limit   query     int     false  "Number of items to return (default: 10, max: 100)"
// @Param        org_id  query     string  false  "Comma-separated list of organization IDs to filter by"
// @Success      200     {object}  getPostsResponse
// @Failure      400     {object}  errorResponse
// @Failure      500     {object}  errorResponse
// @Router       /business/posts [get]
func (h *handler) GetPosts(c *gin.Context) {
	req := new(getPostsRequest)
	c.ShouldBindQuery(req)

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

	var orgIDs []int
	if req.OrgIDs != "" {
		parts := strings.Split(req.OrgIDs, ",")
		for _, p := range parts {
			if id, err := strconv.Atoi(strings.TrimSpace(p)); err == nil {
				orgIDs = append(orgIDs, id)
			}
		}
	}

	ctx := c.Request.Context()

	posts, err := h.service.GetPosts(ctx, req.Skip, req.Limit, orgIDs)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := make([]dto.PostResponse, 0, len(posts.Items))

	for _, post := range posts.Items {
		response = append(response, dto.PostResponse{
			ID:        post.ID,
			Content:   post.Content,
			CreatedAt: post.CreatedAt,
			Images:    post.Images,
			Org: struct {
				ID          int    `json:"id"`
				Name        string `json:"name"`
				ProfileType string `json:"profileType"`
			}{
				ID:          post.OrgID,
				Name:        post.OrgName,
				ProfileType: post.ProfileType,
			},
		})
	}

	c.JSON(http.StatusOK, getPostsResponse{
		Items: response,
		Total: posts.Total,
	})
}
