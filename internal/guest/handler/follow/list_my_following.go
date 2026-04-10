package follow

import (
	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"net/http"
)

func (h *handler) listMyFollowing(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := httpctx.GetUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	customers, err := h.service.ListMyFollowing(ctx, userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, listCustomersResponse{
		Customers: customersToResponse(customers),
	})
}
