package customer

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
)

func (h *handler) update(c *gin.Context) {
	var req updateRequest
	if err := request.BindJSON(c, &req); err != nil {
		c.Error(err)
		return
	}

	userID, err := httpctx.GetUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	in, err := updateRequestToUpdateCustomer(req, userID)
	if err != nil {
		c.Error(err)
		return
	}

	customer, err := h.service.Update(ctx, in)
	if err != nil {
		c.Error(err)
		return
	}

	resp := updateResponse{Customer: customerToResponse(customer)}
	c.JSON(http.StatusOK, resp)
}

type updateRequest struct {
	UserName  *string `json:"userName" binding:"omitempty,alphanum,min=3,max=30"`
	FirstName *string `json:"firstName" binding:"omitempty,min=2,max=50"`
	LastName  *string `json:"lastName" binding:"omitempty,min=2,max=50"`

	Bio             *string `json:"bio" binding:"omitempty,max=500"`
	AvatarObjectKey *string `json:"avatarObjectKey" binding:"omitempty,max=1024"`
}

type updateResponse struct {
	Customer customerResponse `json:"customer"`
}
