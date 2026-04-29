package post

import (
	"mime/multipart"
	"net/http"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	_ "github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

type createRequest struct {
	VenueID  int64                   `form:"venue_id" binding:"required"`
	Text     string                  `form:"text" binding:"required,max=2000"`
	Rating   int16                   `form:"rating" binding:"required,min=1,max=5"`
	Images   []*multipart.FileHeader `form:"images" binding:"omitempty"`
	Mentions []string                `form:"mentions" binding:"omitempty"`
}

type createResponse struct {
	Post postResponse `json:"post"`
}

type errorResponse struct {
	Error string `json:"error" example:"not found"`
}

// create creates a guest post with optional images.
//
//	@Summary		Create post
//	@Description	Creates a post for the authenticated customer.
//	@Tags			guest-posts
//	@Accept			mpfd
//	@Produce		json
//	@Security		BearerAuth
//	@Param			venue_id	formData	int		true	"Venue ID"
//	@Param			text		formData	string	true	"Post text"
//	@Param			rating		formData	int		true	"Rating (1..5)"
//	@Param			images		formData	file	false	"Post images (jpeg/png, up to 5)"
//	@Success		201			{object}	createResponse
//	@Failure		400			{object}	response.ErrorResponse
//	@Failure		401			{object}	response.ErrorResponse
//	@Failure		403			{object}	response.ErrorResponse
//	@Failure		404			{object}	response.ErrorResponse
//	@Failure		502			{object}	response.ErrorResponse
//	@Failure		500			{object}	response.ErrorResponse
//	@Router			/posts/ [post]
func (h *handler) create(c *gin.Context) {
	if c.ContentType() != gin.MIMEMultipartPOSTForm {
		c.Error(apperror.ErrMultipartFormData)
		return
	}

	var req createRequest
	if err := c.ShouldBind(&req); err != nil {
		c.Error(apperror.BadRequest(err.Error()))
		return
	}

	ctx := c.Request.Context()

	customer, err := h.getAuthenticatedCustomer(c)
	if err != nil {
		c.Error(err)
		return
	}

	req.Text = strings.TrimSpace(req.Text)
	if req.Text == "" {
		c.Error(apperror.BadRequest("text is required"))
		return
	}

	images, err := buildUploadImages(req.Images)
	if err != nil {
		c.Error(err)
		return
	}

	dtoImages := make([]dto.UploadImageInput, 0, len(images))
	for _, image := range images {
		dtoImages = append(dtoImages, dto.UploadImageInput{
			File:        image.File,
			ContentType: image.ContentType,
			FileSize:    image.FileSize,
		})
	}

	in := dto.CreatePostInput{
		CustomerID: customer.ID,
		UserID:     customer.UserID,
		VenueID:    req.VenueID,
		Text:       req.Text,
		Rating:     req.Rating,
		Images:     dtoImages,
		Mentions:   req.Mentions,
	}

	post, err := h.service.Create(ctx, in)
	if err != nil {
		c.Error(err)
		return
	}

	resp := createResponse{Post: postToResponse(post, h.storage, customer)}
	c.JSON(http.StatusCreated, resp)
}
