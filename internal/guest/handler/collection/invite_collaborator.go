package collection

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
)

func (h *handler) inviteCollaborator(c *gin.Context) {
	var params inviteCollaboratorParams
	if err := request.BindUri(c, &params); err != nil {
		c.Error(err)
		return
	}

	var req inviteCollaboratorRequest
	if err := request.BindJSON(c, &req); err != nil {
		c.Error(err)
		return
	}

	customerID, err := httpctx.GetCustomerID(c)
	if err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	in := inviteCollaboratorRequestToInput(params.CollectionID, req.CustomerID, customerID)

	if err := h.service.InviteCollaborator(ctx, in); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

type inviteCollaboratorParams struct {
	CollectionID string `uri:"collectionId" binding:"required,uuid"`
}

type inviteCollaboratorRequest struct {
	CustomerID string `json:"customerId" binding:"required,uuid"`
}
