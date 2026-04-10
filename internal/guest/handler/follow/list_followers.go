package follow

import (
	"github.com/gin-gonic/gin"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"net/http"
)

func (h *handler) listFollowers(c *gin.Context) {
	ctx := c.Request.Context()

	customerID := c.Param("id")
	if customerID == "" {
		c.Error(apperror.ErrInvalidParam)
		return
	}

	customers, err := h.service.ListFollowers(ctx, customerID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, listCustomersResponse{
		Customers: customersToResponse(customers),
	})
}
