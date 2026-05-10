package post

import (
	"github.com/gin-gonic/gin"
	_ "github.com/ua-academy-projects/share-bite/internal/guest/util/response"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
	"net/http"
)

// declineInvitation declines a collaborative post invitation.
//
//	@Summary		Decline post invitation
//	@Description	Declines a pending collaborative post invitation for the authenticated customer.
//	@Tags			guest-posts
//	@Produce		json
//	@Security		BearerAuth
//	@Param			invitationId	path		string	true	"Invitation ID"
//	@Success		204				"No Content"
//	@Failure		400	{object}	response.ErrorResponse
//	@Failure		401	{object}	response.ErrorResponse
//	@Failure		403	{object}	response.ErrorResponse
//	@Failure		404	{object}	response.ErrorResponse
//	@Failure		500	{object}	response.ErrorResponse
//	@Router			/posts/invitations/{invitationId}/decline [post]
func (h *handler) declineInvitation(c *gin.Context) {
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

	err = h.service.DeclineInvitation(c.Request.Context(), params.InvitationID, customer.ID)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
