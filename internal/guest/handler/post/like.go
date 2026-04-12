package post

import (
	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
	"net/http"
)

type likeUriRequest struct {
	ID string `uri:"id" binding:"required,numeric"`
}

// like likes a post on behalf of authenticated customer.
//
// @Summary      Like post
// @Description  Adds authenticated customer like to the post.
// @Tags         guest-posts
// @Produce      json
// @Security     BearerAuth
// @Param        id   path  int  true  "Post ID"
// @Success      204
// @Failure      400  {object}  errorResponse
// @Failure      401  {object}  errorResponse
// @Failure      404  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /posts/{id}/like [post]
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

	c.Status(http.StatusNoContent)
}

// unlike removes authenticated customer like from a post.
//
// @Summary      Unlike post
// @Description  Removes authenticated customer like from the post.
// @Tags         guest-posts
// @Produce      json
// @Security     BearerAuth
// @Param        id   path  int  true  "Post ID"
// @Success      200
// @Failure      400  {object}  errorResponse
// @Failure      401  {object}  errorResponse
// @Failure      404  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /posts/{id}/like [delete]
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

	c.Status(http.StatusOK)
}
