package comment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/internal/guest/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/guest/util/request"
)

type updateUriRequest struct {
	PostID    int64 `uri:"id" binding:"required,numeric"`
	CommentID int64 `uri:"comment_id" binding:"required,numeric"`
}

type updateRequest struct {
	Text string `json:"text" binding:"required,max=1000"`
}

func (h *handler) update(c *gin.Context) {
	var uriReq updateUriRequest
	if err := request.BindUri(c, &uriReq); err != nil {
		c.Error(err)
		return
	}

	var req updateRequest
	if err := request.BindJSON(c, &req); err != nil {
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

	in := entity.UpdateCommentInput{
		CommentID:  uriReq.CommentID,
		CustomerID: customer.ID,
		Text:       req.Text,
	}

	comment, err := h.service.Update(ctx, uriReq.PostID, in)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"comment": commentToResponse(comment, customer)})
}
