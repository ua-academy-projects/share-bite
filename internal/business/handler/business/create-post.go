package business

import (
	"errors"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	biserr "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/internal/business/mapper"
	repo "github.com/ua-academy-projects/share-bite/internal/business/repository/business"
)

type CreatePostInput struct {
	TextData string                  `form:"content" binding:"required,min=5"`
	Images   []*multipart.FileHeader `form:"photos" binding:"required"`
}

type CreatePostResponse struct {
	ID        int64     `json:"id"`
	OrgId     int64     `json:"org_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Images    []string  `json:"images"`
}

func (h *handler) CreatePost(c *gin.Context) {
	var input CreatePostInput

	err := c.ShouldBind(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	unitIDStr := c.Param("id")
	unitID, err := strconv.Atoi(unitIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	err = h.service.CheckOwnership(ctx, userID, unitID)
	if err != nil {
		if errors.Is(err, biserr.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied to this unit"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "auth check failed"})
		return
	}

	post, err := h.service.CreatePost(ctx, userID, unitID, input.TextData, input.Images)
	if err != nil {
		switch{
		case errors.Is(err, biserr.WrongFileExtErr) || errors.Is(err, biserr.FileToLargeErr):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, repo.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, mapper.ToPostResponse(post))
}