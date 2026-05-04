package post

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) getMyInvitations(c *gin.Context) {
	customer, err := h.getAuthenticatedCustomer(c)
	if err != nil {
		c.Error(err)
		return
	}

	res, err := h.service.GetMyPostInvitations(c.Request.Context(), customer.ID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}
