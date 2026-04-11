package post

import (
	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
)

const (
	maxAvatarSizeBytes = 5 * 1024 * 1024
	fileSniffSizeBytes = 512
)

type createRequest struct {
	VenueID int64                   `form:"venue_id" binding:"required"`
	Text    string                  `form:"text" binding:"required,max=2000"`
	Rating  int16                   `form:"rating" binding:"required,min=1,max=5"`
	Images  []*multipart.FileHeader `form:"images" binding:"omitempty,max=5"`
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

	var images []entity.UploadImageInput

	for _, f := range req.Images {
		if f.Size > maxAvatarSizeBytes {
			c.Error(apperror.BadRequest("image too large"))
			return
		}
		file, err := f.Open()
		if err != nil {
			c.Error(err)
			return
		}

		buffer := make([]byte, fileSniffSizeBytes)
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			file.Close()
			c.Error(err)
			return
		}

		contentType := http.DetectContentType(buffer[:n])
		if !isAllowedImageContentType(contentType) {
			file.Close()
			c.Error(apperror.ErrUnsupportedImageType)
			return
		}
		seeker, ok := file.(io.Seeker)
		if !ok {
			file.Close()
			c.Error(apperror.Internal("uploaded file is not seekable"))
			return
		}

		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			file.Close()
			c.Error(err)
			return
		}

		images = append(images, entity.UploadImageInput{
			File:        file,
			ContentType: contentType,
			FileSize:    f.Size,
		})
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

func isAllowedImageContentType(contentType string) bool {
	switch contentType {
	case "image/jpeg", "image/png":
		return true
	default:
		return false
	}
}
