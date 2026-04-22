package post

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
)

type likeUriRequest struct {
	ID string `uri:"id" binding:"required,numeric"`
}

// like likes a post on behalf of authenticated customer.
//
//	@Summary		Like post
//	@Description	Adds authenticated customer like to the post.
//	@Tags			guest-posts
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path	int	true	"Post ID"
//	@Success		204	"Successfully liked the post"
//	@Failure		400	{object}	errorResponse	"Invalid post ID format"
//	@Failure		401	{object}	errorResponse	"Unauthorized: token is missing, invalid, or expired"
//	@Failure		403	{object}	errorResponse	"Forbidden: customer profile was not found"
//	@Failure		404	{object}	errorResponse	"Not found: post not found or not accessible"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/posts/{id}/like [post]
func (h *handler) like(c *gin.Context) {
	var uriReq likeUriRequest
	if err := request.BindUri(c, &uriReq); err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	customer, err := h.getAuthenticatedCustomer(c)
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.service.Like(ctx, uriReq.ID, customer.ID); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

// unlike removes authenticated customer like from a post.
//
//	@Summary		Unlike post
//	@Description	Removes authenticated customer like from the post.
//	@Tags			guest-posts
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path	int	true	"Post ID"
//	@Success		200	"Successfully unliked the post"
//	@Failure		400	{object}	errorResponse	"Invalid post ID format"
//	@Failure		401	{object}	errorResponse	"Unauthorized: token is missing, invalid, or expired"
//	@Failure		403	{object}	errorResponse	"Forbidden: customer profile was not found"
//	@Failure		404	{object}	errorResponse	"Not found: post not found or not accessible"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/posts/{id}/like [delete]
func (h *handler) unlike(c *gin.Context) {
	var uriReq likeUriRequest
	if err := request.BindUri(c, &uriReq); err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	customer, err := h.getAuthenticatedCustomer(c)
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.service.Unlike(ctx, uriReq.ID, customer.ID); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}
