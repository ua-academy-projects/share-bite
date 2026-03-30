package customer

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
)

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

	resp := getByUserNameResponse{Customer: customerToResponse(customer, h.storage)}
	c.JSON(http.StatusOK, resp)
}

type getByUserNameRequest struct {
	UserName string `uri:"username" binding:"required,alphanum,min=3,max=30"`
}

type getByUserNameResponse struct {
	Customer customerResponse `json:"customer"`
}
