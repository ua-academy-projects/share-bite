package follow

import (
	"github.com/gin-gonic/gin"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"net/http"
)

func (h *handler) listFollowers(c *gin.Context) {
	ctx := c.Request.Context()

	targetCustomerID := c.Param("id")
	if targetCustomerID == "" {
		c.Error(apperror.ErrInvalidParam)
		return
	}

	requesterUserID, _ := httpctx.GetUserID(c)

	customers, err := h.service.ListFollowers(ctx, targetCustomerID, requesterUserID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, listCustomersResponse{
		Customers: customersToResponse(customers),
	})
}
