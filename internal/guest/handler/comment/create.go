package comment

import (
	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
)

type uriRequest struct {
	ID int64 `uri:"id" binding:"required,numeric"`
}

type createRequest struct {
	Text string `json:"text" binding:"required,max=1000"`
}

type createResponse struct {
	Comment commentResponse `json:"comment"`
}

func (h *handler) create(c *gin.Context) {
	var uriReq uriRequest
	if err := request.BindUri(c, &uriReq); err != nil {
		c.Error(err)
		return
	}

	var req createRequest
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

	in := dto.CreateCommentInput{
		PostID:     uriReq.ID,
		CustomerID: customer.ID,
		Text:       req.Text,
	}

	comment, err := h.service.Create(ctx, in)
	if err != nil {
		c.Error(err)
		return
	}

	resp := createResponse{Comment: commentToResponse(comment, customer)}
	c.JSON(http.StatusCreated, resp)
}
