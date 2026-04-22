package comment

import (
	"net/http"

	"github.com/ua-academy-projects/share-bite/internal/guest/dto"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
)

type listUriRequest struct {
	PostID int64 `uri:"id" binding:"required,numeric"`
}

type listQueryRequest struct {
	Take int `form:"take,default=20" binding:"gte=1,lte=100"`
	Skip int `form:"skip,default=0" binding:"gte=0"`
}

type listResponse struct {
	Total    int               `json:"total"`
	Entities []commentResponse `json:"entities"`
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
	in := dto.ListCommentsInput{
		PostID: uriReq.PostID,
		Limit:  queryReq.Take,
		Offset: queryReq.Skip,
	}

	out, err := h.service.List(ctx, in)
	if err != nil {
		c.Error(err)
		return
	}

	resp := listResponse{
		Total:    out.Total,
		Entities: make([]commentResponse, 0, len(out.Comments)),
	}

	for _, item := range out.Comments {
		resp.Entities = append(resp.Entities, commentToResponse(item.Comment, item.Customer))
	}

	c.JSON(http.StatusOK, resp)
}
