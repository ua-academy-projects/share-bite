package collection

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
)

const (
	defaultInvitationsLimit = 20
	maxInvitationsLimit     = 100
)

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

	PageSize  int    `form:"pageSize" binding:"omitempty,min=1,max=100"`
	PageToken string `form:"pageToken"`
}

type listInvitationsResponse struct {
	Invitations   []invitationResponse `json:"invitations"`
	NextPageToken string               `json:"nextPageToken,omitempty"`
}
