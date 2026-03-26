package business

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/ua-academy-projects/share-bite/internal/middleware"

	"github.com/gin-gonic/gin"
	biserr "github.com/ua-academy-projects/share-bite/internal/business/error"
)

type handler struct {
	service businessService
}

type businessService interface {
	UpdatePost(ctx context.Context, postID int64, userID int64, content string) error
	DeletePost(ctx context.Context, postID int64, userID int64) error
}

type updatePostRequest struct {
	Content string `json:"content" binding:"required"`
}

func RegisterHandlers(
	r *gin.RouterGroup,
	service businessService,
) {
	h := &handler{
		service: service,
	}

	r.PUT("/posts/:id", h.UpdatePost)
	r.DELETE("/posts/:id", h.DeletePost)
}

func getUserID(c *gin.Context) (int64, bool) {
	val, exists := c.Get(middleware.CtxUserID)
	if !exists {
		return 0, false
	}

	userIDStr, ok := val.(string)
	if !ok {
		return 0, false
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0, false
	}

	return userID, true
}

func checkBusinessRole(c *gin.Context) bool {
	val, exists := c.Get(middleware.CtxUserRole)
	if !exists {
		return false
	}

	role, ok := val.(string)
	if !ok {
		return false
	}

	return role == "business"
}

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
		case errors.Is(err, biserr.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found or access denied"})
		case errors.Is(err, biserr.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post deleted"})
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

	err = h.service.UpdatePost(ctx, postID, userID, req.Content)
	if err != nil {
		switch {
		case errors.Is(err, biserr.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found or access denied"})
		case errors.Is(err, biserr.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post updated"})
}
