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

// CreateComment adds a comment to a post.
//
//	@Summary		Create comment
//	@Description	Adds a new comment to the specified post.
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"Post ID"
//	@Param			request	body		createCommentRequest	true	"Comment content"
//	@Success		201		{object}	dto.CreateCommentResponse
//	@Failure		400		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/business/posts/{id}/comments [post]
func (h *handler) CreateComment(c *gin.Context) {
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

// GetComments returns a list of comments for a post.
//
//	@Summary		Get comments
//	@Description	Returns a paginated list of comments for the specified post.
//	@Tags			posts
//	@Produce		json
//	@Param			id		path		int	true	"Post ID"
//	@Param			limit	query		int	false	"Limit"
//	@Param			offset	query		int	false	"Offset"
//	@Success		200		{object}	dto.GetCommentsResponse
//	@Failure		400		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/business/posts/{id}/comments [get]
func (h *handler) GetComments(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || postID <= 0 {
		c.Error(apperror.BadRequest("invalid post id"))
		return
	}

	limit := 20
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit >= 0 {
			limit = parsedLimit
			if limit > 100 {
				limit = 100
			}
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

// UpdateComment modifies an existing comment.
//
//	@Summary		Update comment
//	@Description	Updates the content of an existing comment.
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id			path		int						true	"Post ID"
//	@Param			comment_id	path		int						true	"Comment ID"
//	@Param			request		body		updateCommentRequest	true	"New comment content"
//	@Success		200			{object}	dto.UpdateCommentResponse
//	@Failure		400			{object}	errorResponse
//	@Failure		500			{object}	errorResponse
//	@Router			/business/posts/{id}/comments/{comment_id} [patch]
func (h *handler) UpdateComment(c *gin.Context) {
	commentID, err := strconv.ParseInt(c.Param("comment_id"), 10, 64)
	if err != nil || commentID <= 0 {
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

// DeleteComment removes a comment from a post.
//
//	@Summary		Delete comment
//	@Description	Deletes an existing comment. Authors can delete their own comments, and business owners can delete any comment on their posts.
//	@Tags			posts
//	@Security		BearerAuth
//	@Param			id			path	int	true	"Post ID"
//	@Param			comment_id	path	int	true	"Comment ID"
//	@Success		204			"No Content"
//	@Failure		400			{object}	errorResponse
//	@Failure		403			{object}	errorResponse
//	@Failure		404			{object}	errorResponse
//	@Failure		500			{object}	errorResponse
//	@Router			/business/posts/{id}/comments/{comment_id} [delete]
func (h *handler) DeleteComment(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || postID <= 0 {
		c.Error(apperror.BadRequest("invalid post id"))
		return
	}

	commentID, err := strconv.ParseInt(c.Param("comment_id"), 10, 64)
	if err != nil || commentID <= 0 {
		c.Error(apperror.BadRequest("invalid comment id"))
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.Error(apperror.BadRequest("unauthorized"))
		return
	}

	err = h.service.DeleteComment(c.Request.Context(), postID, commentID, userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
