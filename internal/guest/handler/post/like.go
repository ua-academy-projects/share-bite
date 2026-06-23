package post

import (
	"github.com/gin-gonic/gin"
	_ "github.com/ua-academy-projects/share-bite/internal/guest/util/response"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
	"net/http"
)

type likeUriRequest struct {
	ID string `uri:"id" binding:"required,numeric"`
}

// like adds a like to a post.
//
//	@Summary		Like post
//	@Description	Adds a like to the specified post.
//	@Tags			guest-posts
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path	int	true	"Post ID"
//	@Success 		200 {object} 	nil
//	@Failure		400	{object}	response.ErrorResponse
//	@Failure		401	{object}	response.ErrorResponse
//	@Failure		404	{object}	response.ErrorResponse
//	@Failure		500	{object}	response.ErrorResponse
//	@Router			/posts/{id}/like [post]
func (h *handler) like(c *gin.Context) {
	var uriReq likeUriRequest
	if err := request.BindUri(c, &uriReq); err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	userID, err := httpctx.GetUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	customer, err := h.customerService.GetByUserID(ctx, userID)
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.service.Like(ctx, uriReq.ID, customer.ID); err != nil {
		c.Error(err)
		return
	}
	if h.metrics != nil {
		h.metrics.RecordPostLike()
	}

	c.Status(http.StatusNoContent)
}

// unlike removes a like from a post.
//
//	@Summary		Unlike post
//	@Description	Removes a like from the specified post.
//	@Tags			guest-posts
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path	int	true	"Post ID"
//	@Success 		204 {object} 	nil
//	@Failure		400	{object}	response.ErrorResponse
//	@Failure		401	{object}	response.ErrorResponse
//	@Failure		404	{object}	response.ErrorResponse
//	@Failure		500	{object}	response.ErrorResponse
//	@Router			/posts/{id}/like [delete]
func (h *handler) unlike(c *gin.Context) {
	var uriReq likeUriRequest
	if err := request.BindUri(c, &uriReq); err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	userID, err := httpctx.GetUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	customer, err := h.customerService.GetByUserID(ctx, userID)
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.service.Unlike(ctx, uriReq.ID, customer.ID); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
