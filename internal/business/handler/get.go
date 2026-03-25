package business

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type getter interface {
	GetById(ctx context.Context, id in)
}

func (h *handler) getById(c *gin.Context) {
	id := c.Param("id")
	tokens, err := h.service.GetById(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
