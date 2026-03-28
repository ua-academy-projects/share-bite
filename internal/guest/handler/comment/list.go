package comment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/internal/guest/util/request"
)

type listUriRequest struct {
	PostID int64 `uri:"id" binding:"required,numeric"`
}

type listQueryRequest struct {
	PageSize  int    `form:"page_size,default=20" binding:"gte=1,lte=100"`
	PageToken string `form:"page_token"`
}

type listResponse struct {
	Comments      []commentResponse `json:"comments"`
	NextPageToken string            `json:"next_page_token,omitempty"`
}

func (h *handler) list(c *gin.Context) {
	var uriReq listUriRequest
	if err := request.BindUri(c, &uriReq); err != nil {
		c.Error(err)
		return
	}

	var queryReq listQueryRequest
	if err := request.BindQuery(c, &queryReq); err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	in := entity.ListCommentsInput{
		PostID:    uriReq.PostID,
		PageSize:  queryReq.PageSize,
		PageToken: queryReq.PageToken,
	}

	out, err := h.service.List(ctx, in)
	if err != nil {
		c.Error(err)
		return
	}

	resp := listResponse{
		Comments:      make([]commentResponse, 0, len(out.Comments)),
		NextPageToken: out.NextPageToken,
	}

	for _, item := range out.Comments {
		resp.Comments = append(resp.Comments, commentToResponse(item.Comment, item.Customer))
	}

	c.JSON(http.StatusOK, resp)
}
