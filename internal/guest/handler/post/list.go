package post

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

func (h *handler) list(c *gin.Context) {
	req := new(listRequest)
	if err := c.ShouldBindQuery(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	in := entity.ListPostsInput{
		Limit:  req.Limit,
		Offset: req.Offset,
	}
	out, err := h.service.List(ctx, in)
	if err != nil {
		c.Error(err)
		return
	}

	resp := listPostsOutToResponse(out)
	c.JSON(http.StatusOK, resp)
}

type listRequest struct {
	Limit  int `form:"limit" binding:"required,gte=1,lte=100"`
	Offset int `form:"offset" binding:"gte=0,lte=1000"`
}

type listResponse struct {
	Posts []postResponse `json:"posts"`
	Total int            `json:"total"`
}

func listPostsOutToResponse(out entity.ListPostsOutput) listResponse {
	list := make([]postResponse, 0, len(out.Posts))
	for _, p := range out.Posts {
		list = append(list, postToResponse(p))
	}

	return listResponse{
		Posts: list,
		Total: out.Total,
	}
}
