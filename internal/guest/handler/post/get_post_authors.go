package post

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) getAuthors(c *gin.Context) {
	postID := c.Param("id")

	authors, err := h.service.GetPostAuthors(c.Request.Context(), postID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"authors": authors,
	})
}
