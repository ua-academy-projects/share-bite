package customer

import (
	"net/http"

	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"

	"github.com/gin-gonic/gin"
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

	resp := getMeResponse{Customer: h.toResponse(customer)}
	c.JSON(http.StatusOK, resp)
}

type getMeResponse struct {
	Customer customerResponse `json:"customer"`
}
