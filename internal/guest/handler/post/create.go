package post

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
)

type createRequest struct {
	VenueID string `json:"venue_id" binding:"required,uuid"`
	Text    string `json:"text" binding:"required,max=2000"`
	Rating  int16  `json:"rating" binding:"required,min=1,max=5"`
}

type createResponse struct {
	Post postResponse `json:"post"`
}

func (h *handler) create(c *gin.Context) {
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

	in := entity.CreatePostInput{
		CustomerID: customer.ID,
		VenueID:    req.VenueID,
		Text:       req.Text,
		Rating:     req.Rating,
	}

	post, err := h.service.Create(ctx, in)
	if err != nil {
		c.Error(err)
		return
	}

	resp := createResponse{Post: postToResponse(post)}
	c.JSON(http.StatusCreated, resp)
}
