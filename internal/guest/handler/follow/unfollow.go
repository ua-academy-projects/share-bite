package follow

import (
	"github.com/gin-gonic/gin"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"net/http"
)

func (h *handler) unfollow(c *gin.Context) {
	userID, err := httpctx.GetUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	followedID := c.Param("id")
	if followedID == "" {
		c.Error(apperror.ErrInvalidParam)
		return
	}

	if err := h.service.Unfollow(c.Request.Context(), userID, followedID); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
