package collection

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/util/request"

	_ "github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

const (
	defaultInvitationsLimit = 20
	maxInvitationsLimit     = 100
)

// @Summary		List invitations
// @Description	Retrieves a paginated list of invitations.
// @Description	You must provide exactly ONE filter parameter: collectionId, inviterId, or inviteeId.
// @Description	Users can only view their own inbound (inviteeId) and outbound (inviterId) invitations.
// @Description	Only the collection owner can view invitations filtered by collectionId.
//
// @Tags			invitations
// @Accept			json
// @Produce		json
// @Security		BearerAuth
//
// @Param			collectionId	query		string						false	"Filter by collection ID (UUID)"
// @Param			inviterId		query		string						false	"Filter by inviter customer ID (UUID). Must be your own ID."
// @Param			inviteeId		query		string						false	"Filter by invitee customer ID (UUID). Must be your own ID."
// @Param			status			query		string						false	"Filter by status"	Enums(pending, accepted, declined)
// @Param			pageSize		query		int							false	"Number of items to return (default is 20, max is 100)"
// @Param			pageToken		query		string						false	"Pagination token returned from a previous request"
//
// @Success		200				{object}	listInvitationsResponse		"Successfully retrieved the list of invitations"
// @Failure		400				{object}	response.ErrorResponse		"Invalid query parameters or missing required filters"
// @Failure		401				{object}	response.AuthErrorResponse	"Unauthorized: Missing or invalid token"
// @Failure		403				{object}	response.ErrorResponse		"Forbidden: Customer profile not found or attempting to view someone else's invitations"
// @Failure		404				{object}	response.ErrorResponse		"Not Found: Collection not found"
// @Failure		500				{object}	response.ErrorResponse		"Internal server error"
//
// @Router			/collections/invitations [get]
func (h *handler) listInvitations(c *gin.Context) {
	var params listInvitationsParams
	if err := request.BindQuery(c, &params); err != nil {
		c.Error(err)
		return
	}

	customerID, err := httpctx.GetCustomerID(c)
	if err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	in, err := listInvitationsRequestToInput(params, customerID)
	if err != nil {
		c.Error(err)
		return
	}

	out, err := h.service.ListInvitations(ctx, in)
	if err != nil {
		c.Error(err)
		return
	}

	resp := h.listInvitationsOutputToResponse(out)
	c.JSON(http.StatusOK, resp)
}

type listInvitationsParams struct {
	CollectionID *string `form:"collectionId" binding:"required_without_all=InviterID InviteeID,omitempty,uuid"`
	InviterID    *string `form:"inviterId" binding:"required_without_all=InviteeID CollectionID,omitempty,uuid"`
	InviteeID    *string `form:"inviteeId" binding:"required_without_all=InviterID CollectionID,omitempty,uuid"`
	Status       *string `form:"status" binding:"omitempty,oneof=pending accepted declined"`

	PageSize  int    `form:"pageSize" binding:"omitempty,gte=1,lte=100"`
	PageToken string `form:"pageToken"`
}

type listInvitationsResponse struct {
	Invitations   []invitationResponse `json:"invitations"`
	NextPageToken string               `json:"nextPageToken,omitempty"`
}
