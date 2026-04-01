package business

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/business/mapper"
	repo "github.com/ua-academy-projects/share-bite/internal/business/repository/business"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
)

type updatePostRequest struct {
	Content string `json:"content" binding:"required"`
}

func (h *handler) UpdatePost(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || postID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID, _ := middleware.GetUserID(c)

	var req updatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	post, err := h.service.UpdatePost(c.Request.Context(), postID, userID, req.Content)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found or access denied"})
		case errors.Is(err, repo.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(http.StatusOK, mapper.ToPostResponse(post))
}
