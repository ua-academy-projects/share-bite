package follow

import (
	"github.com/gin-gonic/gin"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"net/http"
)

func (h *handler) unfollow(c *gin.Context) {
	customerID, err := httpctx.GetCustomerID(c)
	if err != nil {
		c.Error(err)
		return
	}
	targetCustomerID := c.Param("id")
	if targetCustomerID == "" {
		c.Error(apperror.ErrInvalidParam)
		return
	}
	if err := h.service.Unfollow(c.Request.Context(), customerID, targetCustomerID); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
