package business

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	biserr "github.com/ua-academy-projects/share-bite/internal/business/error"
)

type handler struct {
	service businessService
}

type businessService interface {
	UpdatePost(ctx context.Context, postID int64, orgID int, content string) error
	DeletePost(ctx context.Context, postID int64, orgID int) error
}

type updatePostRequest struct {
	Content string `json:"content" binding:"required"`
}

func getOrgID(c *gin.Context) (int, bool) {
	val, exists := c.Get("orgID")
	if !exists {
		return 0, false
	}

	orgID, ok := val.(int)
	return orgID, ok
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

func (h *handler) DeletePost(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	orgID, ok := getOrgID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ctx := c.Request.Context()
	err = h.service.DeletePost(ctx, id, orgID)
	if err != nil {
		switch err {
		case biserr.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		case biserr.ErrForbidden:
			c.JSON(http.StatusForbidden, gin.H{"error": "you are not the author"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post deleted"})
}

func (h *handler) UpdatePost(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	orgID, ok := getOrgID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req updatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	ctx := c.Request.Context()
	err = h.service.UpdatePost(ctx, id, orgID, req.Content)
	if err != nil {
		switch err {
		case biserr.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		case biserr.ErrForbidden:
			c.JSON(http.StatusForbidden, gin.H{"error": "you are not the author"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post updated"})
}
