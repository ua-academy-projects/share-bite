package customer

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/ua-academy-projects/share-bite/internal/guest/util/response"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
)

// @Summary		Create a customer profile
// @Description	Creates a new customer profile for the authenticated user.
// @Description	The username must be unique across the system.
//
// @Tags			customers
// @Accept			json
// @Produce		json
// @Security		BearerAuth
//
// @Param			request	body		createRequest				true	"Customer creation data"
//
// @Success		201		{object}	createResponse				"Customer profile successfully created"
// @Failure		400		{object}	response.ErrorResponse		"Invalid JSON payload or validation error"
// @Failure		401		{object}	response.AuthErrorResponse	"Unauthorized: Missing or invalid token"
// @Failure		409		{object}	response.ErrorResponse		"Conflict: Username already taken"
// @Failure		500		{object}	response.ErrorResponse		"Internal server error"
//
// @Router			/customers [post]
func (h *handler) create(c *gin.Context) {
	var req createRequest
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
	in, err := createRequestToCreateCustomer(req, userID)
	if err != nil {
		c.Error(err)
		return
	}

	customerID, err := h.service.Create(ctx, in)
	if err != nil {
		c.Error(err)
		return
	}

	resp := createResponse{CustomerID: customerID}
	c.JSON(http.StatusCreated, resp)
}

type createRequest struct {
	UserName  string `json:"userName" binding:"required,alphanum,min=3,max=30"`
	FirstName string `json:"firstName" binding:"required,min=2,max=50"`
	LastName  string `json:"lastName" binding:"required,min=2,max=50"`

	Bio *string `json:"bio" binding:"omitempty,max=500"`
}

type createResponse struct {
	CustomerID string `json:"customerId"`
}
