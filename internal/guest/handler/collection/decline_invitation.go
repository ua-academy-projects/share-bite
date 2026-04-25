package collection

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
)

func (h *handler) declineInvitation(c *gin.Context) {
	var params declineInvitationParams
	if err := request.BindUri(c, &params); err != nil {
		c.Error(err)
		return
	}

	customerID, err := httpctx.GetCustomerID(c)
	if err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	if err := h.service.DeclineInvitation(ctx, params.InvitationID, customerID); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

type declineInvitationParams struct {
	InvitationID string `uri:"invitationId" binding:"required,uuid"`
}
