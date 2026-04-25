package customer

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/ua-academy-projects/share-bite/internal/guest/util/response"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
)

// @Summary		Get customer by username
// @Description	Retrieves a public customer profile by their unique username.
//
// @Tags			customers
// @Produce		json
//
// @Param			username	path		string					true	"Customer username"
//
// @Success		200			{object}	getByUserNameResponse	"Successfully retrieved customer profile"
// @Failure		400			{object}	response.ErrorResponse	"Bad Request: Invalid username format"
// @Failure		404			{object}	response.ErrorResponse	"Not Found: Customer with this username does not exist"
// @Failure		500			{object}	response.ErrorResponse	"Internal server error"
//
// @Router			/customers/{username} [get]
func (h *handler) getByUserName(c *gin.Context) {
	var req getByUserNameRequest
	if err := request.BindUri(c, &req); err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	customer, err := h.service.GetByUserName(ctx, req.UserName)
	if err != nil {
		c.Error(err)
		return
	}

	resp := getByUserNameResponse{Customer: h.toResponse(customer)}
	c.JSON(http.StatusOK, resp)
}

type getByUserNameRequest struct {
	UserName string `uri:"username" binding:"required,alphanum,min=3,max=30"`
}

type getByUserNameResponse struct {
	Customer customerResponse `json:"customer"`
}
