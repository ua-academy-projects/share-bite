package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
)

type getPostsRequest struct {
	Skip  int `form:"skip"`
	Limit int `form:"limit"`
}

type getPostsResponse struct {
	Items []dto.PostResponse `json:"items"`
	Total int                `json:"total"`
}

// GetPosts returns a paginated list of posts.
//
//	@Summary		Get posts
//	@Description	Returns paginated posts with images and organization info
//	@Tags			posts
//	@Produce		json
//	@Param			skip	query		int	false	"Number of items to skip (default: 0)"
//	@Param			limit	query		int	false	"Number of items to return (default: 10, max: 100)"
//	@Success		200		{object}	getPostsResponse
//	@Failure		400		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/business/posts [get]
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

	ctx := c.Request.Context()

	posts, err := h.service.GetPosts(ctx, req.Skip, req.Limit)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal error",
		})
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
