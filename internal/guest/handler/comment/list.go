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
	Limit  int `form:"limit,default=20" binding:"gte=1,lte=100"`
	Offset int `form:"offset,default=0" binding:"gte=0,lte=1000"`
}

type listResponse struct {
	Comments []commentResponse `json:"comments"`
	Total    int               `json:"total"`
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
		PostID: uriReq.PostID,
		Limit:  queryReq.Limit,
		Offset: queryReq.Offset,
	}

	out, err := h.service.List(ctx, in)
	if err != nil {
		c.Error(err)
		return
	}

	resp := listResponse{
		Comments: make([]commentResponse, 0, len(out.Comments)),
		Total:    out.Total,
	}

	for _, comment := range out.Comments {
		resp.Comments = append(resp.Comments, commentToResponse(comment))
	}

	c.JSON(http.StatusOK, resp)
}
