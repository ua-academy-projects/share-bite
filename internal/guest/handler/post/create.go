package post

import (
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
)

type createRequest struct {
	VenueID int64                   `form:"venue_id" binding:"required"`
	Text    string                  `form:"text" binding:"required,max=2000"`
	Rating  int16                   `form:"rating" binding:"required,min=1,max=5"`
	Images  []*multipart.FileHeader `form:"images" binding:"omitempty"`
}

type createResponse struct {
	Post postResponse `json:"post"`
}

func (h *handler) create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBind(&req); err != nil {
		c.Error(apperror.BadRequest(err.Error()))
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

	images, err := buildUploadImages(req.Images)
	if err != nil {
		c.Error(err)
		return
	}

	in := entity.CreatePostInput{
		CustomerID: customer.ID,
		VenueID:    req.VenueID,
		Text:       req.Text,
		Rating:     req.Rating,
		Images:     images,
	}

	post, err := h.service.Create(ctx, in)
	if err != nil {
		c.Error(err)
		return
	}

	resp := createResponse{Post: postToResponse(post, h.storage)}
	c.JSON(http.StatusCreated, resp)
}
