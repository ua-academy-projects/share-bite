package post

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) get(c *gin.Context) {
	req := new(getRequest)
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
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
	ID string `uri:"id" binding:"required,uuid"`
}

type getResponse struct {
	Post postResponse `json:"post"`
}
