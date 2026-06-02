package post

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/storage/mediatype"
)

const (
	maxPostImages      = 5
	fileSniffSizeBytes = 512
)

func buildUploadImages(files []*multipart.FileHeader) ([]entity.UploadImageInput, error) {
	if len(files) > maxPostImages {
		return nil, apperror.BadRequest(fmt.Sprintf("too many images: max %d", maxPostImages))
	}

	images := make([]entity.UploadImageInput, 0, len(files))

	for _, f := range files {
		file, err := f.Open()
		if err != nil {
			return nil, err
		}

		buffer := make([]byte, fileSniffSizeBytes)
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			file.Close()
			return nil, err
		}

		contentType := http.DetectContentType(buffer[:n])
		if err := mediatype.DefaultImageValidator.Validate(contentType, f.Size); err != nil {
			file.Close()

			if errors.Is(err, mediatype.ErrUnsupportedType) {
				return nil, apperror.ErrUnsupportedImageType
			}
			if errors.Is(err, mediatype.ErrFileTooLarge) {
				return nil, apperror.BadRequest(err.Error())
			}
			return nil, err
		}

		seeker, ok := file.(io.Seeker)
		if !ok {
			file.Close()
			return nil, apperror.Internal("uploaded file is not seekable")
		}

		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			file.Close()
			return nil, err
		}

		images = append(images, entity.UploadImageInput{
			File:        file,
			ContentType: contentType,
			FileSize:    f.Size,
		})
	}

	return images, nil
}

func multipartFieldProvided(c *gin.Context, field string) bool {
	if c.Request.MultipartForm == nil {
		return false
	}

	if _, ok := c.Request.MultipartForm.File[field]; ok {
		return true
	}

	if _, ok := c.Request.MultipartForm.Value[field]; ok {
		return true
	}

	return false
}
