package post

import (
	"net/http"

	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
)

// get returns a single published post by ID.
//
//	@Summary		Get post by ID
//	@Description	Returns a published post by its numeric ID.
//	@Tags			guest-posts
//	@Produce		json
//	@Param			id	path		int				true	"Post ID"
//	@Success		200	{object}	getResponse		"Successfully retrieved the post"
//	@Failure		400	{object}	errorResponse	"Invalid post ID format"
//	@Failure		404	{object}	errorResponse	"Not found: post does not exist, is private, or is not published"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/posts/{id} [get]
func (h *handler) get(c *gin.Context) {
	var req getRequest
	if err := request.BindUri(c, &req); err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	customerID := getOptionalCustomerID(c, h.customerService)
	post, err := h.service.Get(ctx, req.ID, customerID)
	if err != nil {
		c.Error(err)
		return
	}

	resp := getResponse{Post: postToResponse(post, h.storage)}
	c.JSON(http.StatusOK, resp)
}

type getRequest struct {
	ID string `uri:"id" binding:"required,numeric"`
}

type getResponse struct {
	Post postResponse `json:"post"`
}

func getOptionalCustomerID(c *gin.Context, custSvc customerService) string {
	userID, err := httpctx.GetUserID(c)
	if err == nil && userID != "" {
		customer, err := custSvc.GetByUserID(c.Request.Context(), userID)
		if err == nil {
			return customer.ID
		}
	}
	return ""
}
