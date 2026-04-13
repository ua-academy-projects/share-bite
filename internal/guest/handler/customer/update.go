package customer

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/ua-academy-projects/share-bite/internal/guest/util/response"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
)

// @Summary		Update customer profile
// @Description	Updates the customer profile of the authenticated user.
// @Description	All fields are optional; only provided fields will be updated.
//
// @Tags			customers
// @Accept			json
// @Produce		json
// @Security		BearerAuth
//
// @Param			request	body		updateRequest				true	"Customer update data"
//
// @Success		200		{object}	updateResponse				"Successfully updated customer profile"
// @Failure		400		{object}	response.ErrorResponse		"Invalid JSON, validation error, or empty update payload"
// @Failure		401		{object}	response.AuthErrorResponse	"Unauthorized: Missing or invalid token"
// @Failure		404		{object}	response.ErrorResponse		"Not Found: Customer profile does not exist"
// @Failure		409		{object}	response.ErrorResponse		"Conflict: Username already taken"
// @Failure		500		{object}	response.ErrorResponse		"Internal server error"
//
// @Router			/customers [patch]
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

	resp := updateResponse{Customer: h.toResponse(customer)}
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
