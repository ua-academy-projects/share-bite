package post

import (
	"github.com/gin-gonic/gin"
	_ "github.com/ua-academy-projects/share-bite/internal/guest/util/response"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
	"net/http"
)

// getAuthors returns all authors of a collaborative post.
//
//	@Summary		Get post authors
//	@Description	Returns all authors (owner and accepted collaborators) of the post.
//	@Tags			guest-posts
//	@Produce		json
//	@Param			id	path		int						true	"Post ID"
//	@Success		200	{object}	getAuthorsResponse
//	@Failure		400	{object}	response.ErrorResponse
//	@Failure		404	{object}	response.ErrorResponse
//	@Failure		500	{object}	response.ErrorResponse
//	@Router			/posts/{id}/authors [get]
func (h *handler) getAuthors(c *gin.Context) {
	var req getRequest
	if err := request.BindUri(c, &req); err != nil {
		c.Error(err)
		return
	}
	authors, err := h.service.GetPostAuthors(
		c.Request.Context(),
		req.ID,
	)
	if err != nil {
		c.Error(err)
		return
	}
	resp := getAuthorsResponse{
		Authors: authors,
		Count:   len(authors),
	}
	c.JSON(http.StatusOK, resp)
}

type getAuthorsResponse struct {
	Authors []string `json:"authors"`
	Count   int      `json:"count"`
}
