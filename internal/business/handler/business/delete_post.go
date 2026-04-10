package business

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	repo "github.com/ua-academy-projects/share-bite/internal/business/repository/business"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
)

// DeletePost deletes a post by ID.
//
// @Summary      Delete post
// @Description  Deletes a post if the user has permission (owner of the org)
// @Tags         posts
// @Produce      json
// @Param        id   path      int  true  "Post ID"
// @Success      204  "No Content"
// @Failure      400  {object}  errorResponse  "invalid id"
// @Failure      401  {object}  errorResponse  "unauthorized"
// @Failure      403  {object}  errorResponse  "forbidden"
// @Failure      404  {object}  errorResponse  "post not found"
// @Failure      500  {object}  errorResponse  "internal error"
// @Security     BearerAuth
// @Router       /business/posts/{id} [delete]
func (h *handler) DeletePost(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || postID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err = h.service.DeletePost(c.Request.Context(), postID, userID)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		case errors.Is(err, repo.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
