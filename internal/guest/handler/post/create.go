package post

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

type createRequest struct {
	CustomerID string `json:"customer_id" binding:"required,uuid"`
	VenueID    string `json:"venue_id" binding:"required,uuid"`
	Text       string `json:"text" binding:"required,max=2000"`
	Rating     int16  `json:"rating" binding:"required,min=1,max=5"`
}

type createResponse struct {
	Post postResponse `json:"post"`
}

func (h *handler) create(c *gin.Context) {
	req := new(createRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	in := entity.CreatePostInput{
		CustomerID: req.CustomerID,
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
