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

type createCommentRequest struct {
	Content string `json:"content" binding:"required,max=1000"`
}

type updateCommentRequest struct {
	Content string `json:"content" binding:"required,max=1000"`
}

func (h *handler) CreateComment(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(apperror.BadRequest("invalid post id"))
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.Error(apperror.BadRequest("unauthorized"))
		return
	}

	var req createCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.BadRequest("invalid request"))
		return
	}

	comment, err := h.service.CreateComment(c.Request.Context(), postID, userID, req.Content)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, dto.CreateCommentResponse{
		Comment: mapper.ToCommentResponse(*comment),
	})
}

func (h *handler) GetComments(c *gin.Context) {
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

	comments, err := h.service.GetComments(c.Request.Context(), postID, limit, offset)
	if err != nil {
		c.Error(err)
		return
	}

	items := make([]dto.CommentWithAuthorResponse, len(comments))
	for i, comment := range comments {
		items[i] = mapper.ToCommentWithAuthorResponse(comment)
	}

	c.JSON(http.StatusOK, dto.GetCommentsResponse{Comments: items})
}

func (h *handler) UpdateComment(c *gin.Context) {
	commentID, err := strconv.ParseInt(c.Param("comment_id"), 10, 64)
	if err != nil {
		c.Error(apperror.BadRequest("invalid comment id"))
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.Error(apperror.BadRequest("unauthorized"))
		return
	}

	var req updateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.BadRequest("invalid request"))
		return
	}

	comment, err := h.service.UpdateComment(c.Request.Context(), commentID, userID, req.Content)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.UpdateCommentResponse{
		Comment: mapper.ToCommentResponse(*comment),
	})
}

func (h *handler) DeleteComment(c *gin.Context) {
	commentID, err := strconv.ParseInt(c.Param("comment_id"), 10, 64)
	if err != nil {
		c.Error(apperror.BadRequest("invalid comment id"))
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.Error(apperror.BadRequest("unauthorized"))
		return
	}

	err = h.service.DeleteComment(c.Request.Context(), commentID, userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
