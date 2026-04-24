package follow

import (
	"github.com/gin-gonic/gin"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"net/http"
)

func (h *handler) follow(c *gin.Context) {
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
	follow, err := h.service.Follow(c.Request.Context(), customerID, targetCustomerID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, follow)
}
