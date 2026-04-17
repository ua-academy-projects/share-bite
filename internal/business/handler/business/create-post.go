package business

import (
	"errors"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	biserr "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/internal/business/mapper"
	repo "github.com/ua-academy-projects/share-bite/internal/business/repository/business"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

type CreatePostInput struct {
	TextData string                  `form:"content" binding:"required,min=5"`
	Images   []*multipart.FileHeader `form:"photos" binding:"required"`
}

// CreatePost creates a new business post with images.
//
// @Summary      Create business post
// @Description  Creates a new post for a specific organizational unit with multiple images
// @Tags         posts
// @Accept       multipart/form-data
// @Produce      json
// @Param        id      path      int                  true  "Unit ID"
// @Param        content formData  string               true  "Post content"
// @Param        photos  formData  file                 true  "Post images"
// @Success      201     {object}  dto.PostResponse
// @Failure      400     {object}  errorResponse
// @Failure      403     {object}  errorResponse
// @Failure      500     {object}  errorResponse
// @Security     BearerAuth
// @Router       /business/posts/{id} [post]
func (h *handler) CreatePost(c *gin.Context) {
	var input CreatePostInput
	ctx := c.Request.Context()

	err := c.ShouldBind(&input)
	if err != nil {
		logger.ErrorKV(ctx, "failed to bind input data", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := middleware.getUserID(c)
	if !ok {
		logger.ErrorKV(ctx, "unauthorized access attempt: user ID not found in gin context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	unitIDStr := c.Param("id")
	unitID, err := strconv.Atoi(unitIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.InfoKV(ctx, "attempting to create post", "unit_id", unitID, "user_id", userID)

	err = h.service.CheckOwnership(ctx, userID, unitID)
	if err != nil {
		logger.ErrorKV(ctx, "ownership check failed", "unit_id", unitID, "user_id", userID, "error", err.Error())
		if errors.Is(err, biserr.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied to this unit"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "auth check failed"})
		return
	}

	post, err := h.service.CreatePost(ctx, userID, unitID, input.TextData, input.Images)
	if err != nil {
		logger.ErrorKV(ctx, "failed to create post in service layer", "unit_id", unitID, "error", err.Error())
		switch {
		case errors.Is(err, biserr.WrongFileExtErr) || errors.Is(err, biserr.FileToLargeErr):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, repo.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	logger.InfoKV(ctx, "post successfully created", "post_id", post.ID, "unit_id", unitID)

	c.JSON(http.StatusCreated, mapper.ToPostResponse(post))
}
