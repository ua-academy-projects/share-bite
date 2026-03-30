package post

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
)

func (h *handler) get(c *gin.Context) {
	var req getRequest
	if err := request.BindUri(c, &req); err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	post, err := h.service.Get(ctx, req.ID)
	if err != nil {
		c.Error(err)
		return
	}

	resp := getResponse{Post: postToResponse(post)}
	c.JSON(http.StatusOK, resp)
}

type getRequest struct {
	ID string `uri:"id" binding:"required,numeric"`
}

type getResponse struct {
	Post postResponse `json:"post"`
}
