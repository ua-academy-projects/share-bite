package post

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/internal/guest/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/guest/util/request"
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
	_, err := httpctx.GetUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	//customer, err := h.customerService.GetByUserID(ctx, userID)
	//if err != nil {
	//	c.Error(err)
	//	return
	//}

	in := entity.CreatePostInput{
		CustomerID: "7b2e1a58-4d9c-4f3a-8b1e-2c7a5d9f0b34",
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
