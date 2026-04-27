package collection

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/util/request"

	_ "github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

// @Summary		Invite a collaborator to a collection
// @Description	Sends an invitation to a customer to become a collaborator in a collection.
// @Description	Only the collection owner can send invitations.
// @Description	If an invitation already exists, it will be refreshed with a new 7-day TTL.
// @Description	A 1-hour cooldown applies between resends.
//
// @Tags			invitations
// @Accept			json
// @Produce		json
// @Security		BearerAuth
//
// @Param			collectionId	path	string						true	"Collection ID (UUID)"
// @Param			request			body	inviteCollaboratorRequest	true	"Invitee details"
//
// @Success		204				"Invitation successfully sent or refreshed"
// @Failure		400				{object}	response.ErrorResponse		"Invalid path parameters or request body"
// @Failure		401				{object}	response.AuthErrorResponse	"Unauthorized: Missing or invalid token"
// @Failure		403				{object}	response.ErrorResponse		"Forbidden: Customer profile not found or user does not own this collection"
// @Failure		404				{object}	response.ErrorResponse		"Not Found: Collection or invitee not found"
// @Failure		409				{object}	response.ErrorResponse		"Conflict: Customer is already a collaborator"
// @Failure		429				{object}	response.ErrorResponse		"Too Many Requests: Please wait before resending the invitation"
// @Failure		500				{object}	response.ErrorResponse		"Internal server error"
//
// @Router			/collections/{collectionId}/invitations [post]
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
