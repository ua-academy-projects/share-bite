package collection

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/util/request"

	_ "github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

// @Summary		Remove a collaborator from a collection
// @Description	Removes a specific collaborator from a collection.
// @Description	Only the collection owner can remove collaborators.
//
// @Tags			collections
// @Accept			json
// @Produce		json
// @Security		BearerAuth
//
// @Param			collectionId	path	string	true	"Collection ID (UUID)"
// @Param			customerId		path	string	true	"Target Customer ID (UUID)"
//
// @Success		204				"Collaborator successfully removed"
// @Failure		400				{object}	response.ErrorResponse		"Invalid path parameters"
// @Failure		401				{object}	response.AuthErrorResponse	"Unauthorized: Missing or invalid token"
// @Failure		403				{object}	response.ErrorResponse		"Forbidden: Customer profile not found or user does not own this collection"
// @Failure		404				{object}	response.ErrorResponse		"Not Found: Collection not found or customer is not a collaborator"
// @Failure		500				{object}	response.ErrorResponse		"Internal server error"
//
// @Router			/collections/{collectionId}/collaborators/{customerId} [delete]
func (h *handler) removeCollaborator(c *gin.Context) {
	var params removeCollaboratorParams
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
	in := removeCollaboratorRequestToRemoveCollaborator(params, customerID)

	if err := h.service.RemoveCollaborator(ctx, in); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

type removeCollaboratorParams struct {
	CollectionID     string `uri:"collectionId" binding:"required,uuid"`
	TargetCustomerID string `uri:"customerId" binding:"required,uuid"`
}
