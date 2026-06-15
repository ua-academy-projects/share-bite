package post

import (
	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
	"net/http"

	_ "github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

// acceptInvitation accepts a collaborative post invitation.
//
//	@Summary		Accept post invitation
//	@Description	Accepts a pending collaborative post invitation for the authenticated customer.
//	@Tags			guest-posts
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Invitation ID"
//	@Success		204				"No Content"
//	@Failure		400	{object}	response.ErrorResponse
//	@Failure		401	{object}	response.ErrorResponse
//	@Failure		403	{object}	response.ErrorResponse
//	@Failure		404	{object}	response.ErrorResponse
//	@Failure		500	{object}	response.ErrorResponse
//	@Router			/posts/invitations/{id}/accept [post]
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
	InvitationID string `uri:"id" binding:"required"`
}
