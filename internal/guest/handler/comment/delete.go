package comment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/guest/util/request"
)

type deleteUriRequest struct {
	PostID    int64 `uri:"id" binding:"required,numeric"`
	CommentID int64 `uri:"comment_id" binding:"required,numeric"`
}

func (h *handler) delete(c *gin.Context) {
	var uriReq deleteUriRequest
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

	err = h.service.Delete(ctx, uriReq.CommentID, customer.ID)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
