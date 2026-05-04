package post

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) declineInvitation(c *gin.Context) {
	customer, err := h.getAuthenticatedCustomer(c)
	if err != nil {
		c.Error(err)
		return
	}

	id := c.Param("id")

	err = h.service.DeclineInvitation(c.Request.Context(), id, customer.ID)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
