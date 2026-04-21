package collection

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
)

func (h *handler) addCollaborator(c *gin.Context) {
	var uri addCollaboratorUri
	if err := request.BindUri(c, &uri); err != nil {
		c.Error(err)
		return
	}

	var req addCollaboratorBody
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
	in := addCollaboratorRequestToAddCollaborator(req, uri.CollectionID, customerID)

	if err := h.service.AddCollaborator(ctx, in); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

type addCollaboratorUri struct {
	CollectionID string `uri:"collectionId" binding:"required,uuid"`
}

type addCollaboratorBody struct {
	TargetCustomerID string `json:"customerId" binding:"required,uuid"`
}
