package follow

import (
	"github.com/gin-gonic/gin"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"net/http"
)

func (h *handler) listFollowing(c *gin.Context) {
	ctx := c.Request.Context()

	targetCustomerID := c.Param("id")
	if targetCustomerID == "" {
		c.Error(apperror.ErrInvalidParam)
		return
	}

	requesterUserID, _ := httpctx.GetUserID(c)

	customers, err := h.service.ListFollowing(ctx, targetCustomerID, requesterUserID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, listCustomersResponse{
		Customers: customersToResponse(customers),
	})
}

func (h *handler) listMyFollowing(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := httpctx.GetUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	currentCustomer, err := h.customerService.GetByUserID(ctx, userID)
	if err != nil {
		c.Error(err)
		return
	}

	customers, err := h.service.ListFollowing(ctx, currentCustomer.ID, userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, listCustomersResponse{
		Customers: customersToResponse(customers),
	})
}
