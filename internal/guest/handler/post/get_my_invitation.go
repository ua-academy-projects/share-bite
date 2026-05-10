package post

import (
	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"net/http"
	"time"
)

func (h *handler) getMyInvitations(c *gin.Context) {
	customer, err := h.getAuthenticatedCustomer(c)
	if err != nil {
		c.Error(err)
		return
	}

	invitations, err := h.service.GetMyPostInvitations(c.Request.Context(), customer.ID)
	if err != nil {
		c.Error(err)
		return
	}

	resp := getInvitationsResponse{
		Invitations: invitationsToResponse(invitations),
		Count:       len(invitations),
	}

	c.JSON(http.StatusOK, resp)
}

type getInvitationsResponse struct {
	Invitations []postInvitationResponse `json:"invitations"`
	Count       int                      `json:"count"`
}

type postInvitationResponse struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	InvitedBy string    `json:"invited_by"`
	Status    string    `json:"status"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

func invitationsToResponse(invitations []entity.PostCollaborator) []postInvitationResponse {
	resp := make([]postInvitationResponse, 0, len(invitations))
	for _, invitation := range invitations {
		resp = append(resp, postInvitationResponse{
			ID:        invitation.ID,
			PostID:    invitation.PostID,
			InvitedBy: invitation.InvitedBy,
			Status:    string(invitation.Status),
			ExpiresAt: invitation.ExpiresAt,
			CreatedAt: invitation.CreatedAt,
		})
	}
	return resp
}
