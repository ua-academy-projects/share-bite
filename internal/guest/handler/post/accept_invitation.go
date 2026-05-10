package post

import (
	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
	"net/http"
)

func (h *handler) acceptInvitation(c *gin.Context) {
	customer, err := h.getAuthenticatedCustomer(c)
	if err != nil {
		c.Error(err)
		return
	}

	var params invitationParams

	if err := request.BindUri(c, &params); err != nil {
		c.Error(err)
		return
	}

	err = h.service.AcceptInvitation(c.Request.Context(), params.InvitationID, customer.ID)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

type invitationParams struct {
	InvitationID string `uri:"invitationId" binding:"required"`
}
