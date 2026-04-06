package business

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
)

func (h *handler) GetPosts(c *gin.Context) {
	ctx := c.Request.Context()

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	posts, err := h.service.GetPosts(ctx, page, limit)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := make([]dto.PostResponse, 0, len(posts))

	for _, post := range posts {
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

	c.JSON(http.StatusOK, response)
}
