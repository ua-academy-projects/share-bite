package business

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/ua-academy-projects/share-bite/internal/business/dto"
	"github.com/ua-academy-projects/share-bite/internal/business/mapper"
	repo "github.com/ua-academy-projects/share-bite/internal/business/repository/business"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
)

type updatePostRequest struct {
	Content string `json:"content" binding:"required"`
}

// UpdatePost updates post content by ID.
//
// @Summary      Update post
// @Description  Updates the content of a post if the user has permission
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        id    path      int                 true  "Post ID"
// @Param        input body      updatePostRequest   true  "Updated post content"
// @Success      200   {object}  dto.PostResponse
// @Failure      400   {object}  errorResponse
// @Failure      401   {object}  errorResponse
// @Failure      403   {object}  errorResponse
// @Failure      404   {object}  errorResponse
// @Failure      500   {object}  errorResponse
// @Security     BearerAuth
// @Router       /business/posts/{id} [put]
func (h *handler) UpdatePost(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if postID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be positive"})
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req updatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	post, err := h.service.UpdatePost(c.Request.Context(), postID, userID, req.Content)
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

	c.JSON(http.StatusOK, mapper.ToPostResponse(post))
}
