package collection

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/util/request"

	_ "github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

// @Summary		Accept a collection invitation
// @Description	Accepts a pending invitation to become a collaborator in a collection.
// @Description	Fails if the invitation has expired or has already been processed.
//
// @Tags			invitations
// @Accept			json
// @Produce		json
// @Security		BearerAuth
//
// @Param			invitationId	path	string	true	"Invitation ID (UUID)"
//
// @Success		204				"Invitation successfully accepted"
// @Failure		400				{object}	response.ErrorResponse		"Invalid invitation ID format or invitation has expired"
// @Failure		401				{object}	response.AuthErrorResponse	"Unauthorized: Missing or invalid token"
// @Failure		403				{object}	response.ErrorResponse		"Forbidden: Customer profile not found or invitation belongs to another user"
// @Failure		404				{object}	response.ErrorResponse		"Not Found: Invitation does not exist"
// @Failure		409				{object}	response.ErrorResponse		"Conflict: Invitation has already been processed"
// @Failure		500				{object}	response.ErrorResponse		"Internal server error"
//
// @Router			/collections/invitations/{invitationId}/accept [post]
func (h *handler) acceptInvitation(c *gin.Context) {
	var params acceptInvitationParams
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
	if err := h.service.AcceptInvitation(ctx, params.InvitationID, customerID); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

type acceptInvitationParams struct {
	InvitationID string `uri:"invitationId" binding:"required,uuid"`
}
