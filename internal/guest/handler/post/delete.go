package post

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
)

type deleteURIRequest struct {
	ID string `uri:"id" binding:"required,numeric"`
}

func (h *handler) delete(c *gin.Context) {
	var uriReq deleteURIRequest
	if err := request.BindUri(c, &uriReq); err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	userID, err := httpctx.GetUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	customer, err := h.customerService.GetByUserID(ctx, userID)
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.service.Delete(ctx, uriReq.ID, customer.ID); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
