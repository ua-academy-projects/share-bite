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

func (h *handler) GetLikes(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
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
