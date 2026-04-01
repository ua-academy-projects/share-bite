package business

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	repo "github.com/ua-academy-projects/share-bite/internal/business/repository/business"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
)

func (h *handler) DeletePost(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || postID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID, _ := middleware.GetUserID(c)

	err = h.service.DeletePost(c.Request.Context(), postID, userID)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
