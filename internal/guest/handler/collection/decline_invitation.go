package collection

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/util/request"

	_ "github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

// @Summary		Decline a collection invitation
// @Description	Declines a pending invitation to a collection.
// @Description	Expired invitations can also be declined to dismiss them.
//
// @Tags			invitations
// @Accept			json
// @Produce		json
// @Security		BearerAuth
//
// @Param			invitationId	path	string	true	"Invitation ID (UUID)"
//
// @Success		204				"Invitation successfully declined"
// @Failure		400				{object}	response.ErrorResponse		"Invalid invitation ID format"
// @Failure		401				{object}	response.AuthErrorResponse	"Unauthorized: Missing or invalid token"
// @Failure		403				{object}	response.ErrorResponse		"Forbidden: Customer profile not found or invitation belongs to another user"
// @Failure		404				{object}	response.ErrorResponse		"Not Found: Invitation does not exist"
// @Failure		409				{object}	response.ErrorResponse		"Conflict: Invitation has already been processed"
// @Failure		500				{object}	response.ErrorResponse		"Internal server error"
//
// @Router			/collections/invitations/{invitationId}/decline [post]
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
