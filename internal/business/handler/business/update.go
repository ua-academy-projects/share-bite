package business

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/business/mapper"
	repo "github.com/ua-academy-projects/share-bite/internal/business/repository/business"
)

type updatePostRequest struct {
	Content string `json:"content" binding:"required"`
}

func (h *handler) UpdatePost(c *gin.Context) {
	idStr := c.Param("id")

	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || postID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if !checkBusinessRole(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "only business accounts can update posts"})
		return
	}

	var req updatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	ctx := c.Request.Context()

	post, err := h.service.UpdatePost(ctx, postID, userID, req.Content)
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
