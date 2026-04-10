package post

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
)

type updateRequest struct {
	VenueID *int64  `json:"venue_id" binding:"omitempty"`
	Text    *string `json:"text" binding:"omitempty,max=2000"`
	Rating  *int16  `json:"rating" binding:"omitempty,min=1,max=5"`
}

type updateURIRequest struct {
	ID string `uri:"id" binding:"required,numeric"`
}

type updateResponse struct {
	Post postResponse `json:"post"`
}

func (h *handler) update(c *gin.Context) {
	var uriReq updateURIRequest
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

	in := entity.UpdatePostInput{
		ID:      uriReq.ID,
		VenueID: req.VenueID,
		Text:    req.Text,
		Rating:  req.Rating,
	}

	post, err := h.service.Update(ctx, in)
	if err != nil {
		c.Error(err)
		return
	}

	resp := updateResponse{Post: postToResponse(post)}
	c.JSON(http.StatusOK, resp)
}
