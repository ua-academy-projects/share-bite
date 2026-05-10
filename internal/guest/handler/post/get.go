package post

import (
	"context"
	"net/http"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"

	"github.com/gin-gonic/gin"
	_ "github.com/ua-academy-projects/share-bite/internal/guest/util/response"
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
//	@Failure		400	{object}	response.ErrorResponse	"Invalid post ID format"
//	@Failure		404	{object}	response.ErrorResponse	"Not found: post does not exist, is private, or is not published"
//	@Failure		500	{object}	response.ErrorResponse	"Internal server error"
//	@Router			/posts/{id} [get]
func (h *handler) get(c *gin.Context) {
	var req getRequest

	if err := request.BindUri(c, &req); err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()

	var customerID string

	optionalCustomerID, err := httpctx.GetOptionalCustomerID(c)
	if err != nil {
		c.Error(err)
		return
	}

	if optionalCustomerID != nil {
		customerID = *optionalCustomerID
	}

	post, err := h.service.Get(ctx, req.ID, customerID)
	if err != nil {
		c.Error(err)
		return
	}

	customer, err := h.customerService.GetByUserID(ctx, post.CustomerID)
	if err != nil {
		customer = entity.Customer{ID: post.CustomerID}
	}

	authors, err := h.buildAuthors(ctx, post.ID)
	if err != nil {
		c.Error(err)
		return
	}

	resp := getResponse{
		Post: postToResponse(post, h.storage, customer, authors),
	}

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

func (h *handler) buildAuthors(ctx context.Context, postID string) ([]authorResponse, error) {
	authorIDs, err := h.service.GetPostAuthors(ctx, postID)
	if err != nil {
		return nil, err
	}

	authorsCustomers, err := h.customerService.GetByIDs(ctx, authorIDs)
	if err != nil {
		return nil, err
	}

	authors := make([]authorResponse, 0, len(authorsCustomers))

	for _, author := range authorsCustomers {
		var avatarURL *string

		if author.AvatarObjectKey != nil && h.storage != nil {
			url := h.storage.BuildURL(*author.AvatarObjectKey)
			avatarURL = &url
		}

		authors = append(authors, authorResponse{
			ID:        author.ID,
			UserName:  author.UserName,
			AvatarURL: avatarURL,
		})
	}
	return authors, nil
}
