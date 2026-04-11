package post

import (
	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"io"
	"net/http"
	"strconv"
)

const (
	maxAvatarSizeBytes = 5 * 1024 * 1024
	fileSniffSizeBytes = 512
)

type createResponse struct {
	Post postResponse `json:"post"`
}

func (h *handler) create(c *gin.Context) {
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

	venueID := c.PostForm("venue_id")
	if venueID == "" {
		c.Error(apperror.BadRequest("venue_id is required"))
		return
	}

	text := c.PostForm("text")
	if text == "" {
		c.Error(apperror.BadRequest("text is required"))
		return
	}
	if len(text) > 2000 {
		c.Error(apperror.BadRequest("text is too long"))
		return
	}

	ratingStr := c.PostForm("rating")
	ratingInt, err := strconv.Atoi(ratingStr)
	if err != nil {
		c.Error(apperror.BadRequest("invalid rating"))
		return
	}
	if ratingInt < 1 || ratingInt > 5 {
		c.Error(apperror.BadRequest("rating must be between 1 and 5"))
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.Error(apperror.BadRequest("invalid multipart form"))
		return
	}

	files := form.File["images"]
	if len(files) > 5 {
		c.Error(apperror.BadRequest("too many images, maximum is 5"))
		return
	}

	var images []entity.UploadImageInput

	for _, f := range files {
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
		VenueID:    venueID,
		Text:       text,
		Rating:     int16(ratingInt),
		Images:     images,
	}
	post, err := h.service.Create(ctx, in)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, createResponse{Post: postToResponse(post, h.storage)})
}

func isAllowedImageContentType(contentType string) bool {
	switch contentType {
	case "image/jpeg", "image/png":
		return true
	default:
		return false
	}
}
