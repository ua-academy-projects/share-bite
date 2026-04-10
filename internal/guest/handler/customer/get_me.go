package customer

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
)

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

	resp := getMeResponse{Customer: CustomerToResponse(customer)}
	c.JSON(http.StatusOK, resp)
}

type getMeResponse struct {
	Customer CustomerResponse `json:"customer"`
}
