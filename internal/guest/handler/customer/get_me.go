package customer

import (
	"net/http"

	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"

	"github.com/gin-gonic/gin"
	_ "github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

// @Summary		Get current customer profile
// @Description	Retrieves the customer profile associated with the currently authenticated user.
//
// @Tags			customers
// @Produce		json
// @Security		BearerAuth
//
// @Success		200	{object}	getMeResponse				"Successfully retrieved customer profile"
// @Failure		401	{object}	response.AuthErrorResponse	"Unauthorized: Missing or invalid token"
// @Failure		404	{object}	response.ErrorResponse		"Not Found: Customer profile does not exist for this user"
// @Failure		500	{object}	response.ErrorResponse		"Internal server error"
//
// @Router			/customers [get]
func (h *handler) getMe(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := httpctx.GetUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	customer, err := h.service.GetByUserID(ctx, userID)
	if err != nil {
		c.Error(err)
		return
	}

	resp := getMeResponse{Customer: h.toResponse(customer)}
	c.JSON(http.StatusOK, resp)
}

type getMeResponse struct {
	Customer customerResponse `json:"customer"`
}
