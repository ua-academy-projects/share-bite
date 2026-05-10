package business

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/internal/business/mapper"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
)

// ToggleLike adds or removes a like from a post.
//
//	@Summary		Toggle like
//	@Description	Adds a like if it doesn't exist, otherwise removes it.
//	@Tags			posts
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Post ID"
//	@Success		200	{object}	dto.ToggleLikeResponse
//	@Failure		400	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/business/posts/{id}/likes [post]
func (h *handler) ToggleLike(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || postID <= 0 {
		c.Error(apperror.BadRequest("invalid post id"))
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.Error(apperror.BadRequest("unauthorized"))
		return
	}

	liked, err := h.service.ToggleLike(c.Request.Context(), postID, userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.ToggleLikeResponse{Liked: liked})
}

// GetLikes returns a list of likes for a post.
//
//	@Summary		Get likes
//	@Description	Returns a paginated list of likes for the specified post.
//	@Tags			posts
//	@Produce		json
//	@Param			id		path		int	true	"Post ID"
//	@Param			limit	query		int	false	"Limit"
//	@Param			offset	query		int	false	"Offset"
//	@Success		200		{object}	dto.GetLikesResponse
//	@Failure		400		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/business/posts/{id}/likes [get]
func (h *handler) GetLikes(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || postID <= 0 {
		c.Error(apperror.BadRequest("invalid post id"))
		return
	}

	limit := 20
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	likes, err := h.service.GetLikes(c.Request.Context(), postID, limit, offset)
	if err != nil {
		c.Error(err)
		return
	}

	items := make([]dto.LikeItem, len(likes))
	for i, like := range likes {
		items[i] = mapper.ToLikeItem(like)
	}

	c.JSON(http.StatusOK, dto.GetLikesResponse{Likes: items})
}
