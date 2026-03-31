package business

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	repo "github.com/ua-academy-projects/share-bite/internal/business/repository/business"
)

func (h *handler) DeletePost(c *gin.Context) {
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
		c.JSON(http.StatusForbidden, gin.H{"error": "only business accounts can delete posts"})
		return
	}

	ctx := c.Request.Context()

	err = h.service.DeletePost(ctx, postID, userID)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrNotFound):
			c.Status(http.StatusNoContent)
		case errors.Is(err, repo.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
