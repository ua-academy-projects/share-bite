package post

import (
	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
	"github.com/ua-academy-projects/share-bite/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
)

func (h *handler) list(c *gin.Context) {
	var req listRequest
	if err := request.BindQuery(c, &req); err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	customerID := getOptionalCustomerID(c, h.customerService)
	in := dto.ListPostsInput{
		Limit:      req.Limit,
		Offset:     req.Offset,
		CustomerID: customerID,
	}
	out, err := h.service.List(ctx, in)
	if err != nil {
		c.Error(err)
		return
	}

	resp := listPostsOutToResponse(out, h.storage)
	c.JSON(http.StatusOK, resp)
}

type listRequest struct {
	Limit  int `form:"limit,default=20" binding:"gte=1,lte=100"`
	Offset int `form:"offset,default=0" binding:"gte=0,lte=1000"`
}

type listResponse struct {
	Posts []postResponse `json:"posts"`
	Total int            `json:"total"`
}

func listPostsOutToResponse(out dto.ListPostsOutput, storage storage.ObjectStorage) listResponse {
	list := make([]postResponse, 0, len(out.Posts))
	for _, p := range out.Posts {
		list = append(list, postToResponse(p, storage))
	}

	return listResponse{
		Posts: list,
		Total: out.Total,
	}
}
