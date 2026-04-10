package follow

import (
	"github.com/gin-gonic/gin"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"net/http"
)

func (h *handler) follow(c *gin.Context) {
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

	follow, err := h.service.Follow(c.Request.Context(), userID, followedID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, followResponse{
		Follow: follow,
	})
}
