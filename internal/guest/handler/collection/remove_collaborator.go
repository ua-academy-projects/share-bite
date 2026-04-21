package collection

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
)

func (h *handler) removeCollaborator(c *gin.Context) {
	var uri removeCollaboratorUri
	if err := request.BindUri(c, &uri); err != nil {
		c.Error(err)
		return
	}

	customerID, err := httpctx.GetCustomerID(c)
	if err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	in := removeCollaboratorRequestToRemoveCollaborator(uri, customerID)

	if err := h.service.RemoveCollaborator(ctx, in); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

type removeCollaboratorUri struct {
	CollectionID     string `uri:"collectionId" binding:"required,uuid"`
	TargetCustomerID string `uri:"customerId" binding:"required,uuid"`
}
