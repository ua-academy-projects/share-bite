package customer

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	_ "github.com/ua-academy-projects/share-bite/internal/guest/util/response"
	"github.com/ua-academy-projects/share-bite/internal/storage"
	"github.com/ua-academy-projects/share-bite/internal/storage/key"
	"github.com/ua-academy-projects/share-bite/internal/storage/mediatype"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

const (
	fileSniffSizeBytes     = 512
	multipartOverheadBytes = 1 << 20 // 1MB overhead for multipart boundaries and headers
)

// @Summary		Upload customer avatar
// @Description	Uploads an image to be used as the customer's avatar.
// @Description	Supported formats: JPEG, PNG. Max size: 5MB.
//
// @Tags			customers
// @Accept			multipart/form-data
// @Produce		json
// @Security		BearerAuth
//
// @Param			image	formData	file						true	"Avatar image file (JPEG or PNG, max 5MB)"
//
// @Success		200		{object}	customerResponse			"Successfully uploaded avatar and updated profile"
// @Failure		400		{object}	response.ErrorResponse		"Bad Request: Missing image, file too large, or unsupported format"
// @Failure		401		{object}	response.AuthErrorResponse	"Unauthorized: Missing or invalid token"
// @Failure		404		{object}	response.ErrorResponse		"Not Found: Customer profile does not exist"
// @Failure		500		{object}	response.ErrorResponse		"Internal server error or storage failure"
//
// @Router			/customers/avatar [post]
func (h *handler) uploadAvatar(c *gin.Context) {
	if h.storage == nil {
		c.Error(apperror.Internal("storage is not configured"))
		return
	}

	userID, err := httpctx.GetUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	currentCustomer, err := h.service.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, mediatype.DefaultMaxImageSizeBytes+multipartOverheadBytes)
	fileHeader, err := c.FormFile("image")
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			c.Error(apperror.BadRequest("image file is too large"))
			return
		}
		c.Error(apperror.ErrImageRequired)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.Error(err)
		return
	}
	defer file.Close()

	buffer := make([]byte, fileSniffSizeBytes)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		c.Error(err)
		return
	}

	contentType := http.DetectContentType(buffer[:n])
	if err := mediatype.DefaultImageValidator.Validate(contentType, fileHeader.Size); err != nil {
		if errors.Is(err, mediatype.ErrUnsupportedType) {
			c.Error(apperror.ErrUnsupportedImageType)
			return
		}

		if errors.Is(err, mediatype.ErrFileTooLarge) {
			c.Error(apperror.BadRequest(err.Error()))
			return
		}

		c.Error(err)
		return
	}

	seeker, ok := file.(io.Seeker)
	if !ok {
		c.Error(apperror.Internal("uploaded file is not seekable"))
		return
	}
	if _, err := seeker.Seek(0, io.SeekStart); err != nil {
		c.Error(err)
		return
	}

	ext, ok := mediatype.ExtFromContentType(contentType)
	if !ok {
		c.Error(apperror.ErrUnsupportedImageType)
		return
	}
	objectKey := key.CustomerAvatarKey(currentCustomer.ID, uuid.NewString(), ext)

	if err := h.storage.Upload(
		c.Request.Context(),
		objectKey,
		contentType,
		file,
	); err != nil {
		c.Error(err)
		return
	}

	customer, err := h.service.Update(c.Request.Context(), entity.UpdateCustomer{
		UserID:          userID,
		AvatarObjectKey: &objectKey,
	})
	if err != nil {
		cleanupDelete(h.storage, objectKey)
		c.Error(err)
		return
	}

	if currentCustomer.AvatarObjectKey != nil && *currentCustomer.AvatarObjectKey != objectKey {
		go cleanupDelete(h.storage, *currentCustomer.AvatarObjectKey)
	}

	c.JSON(http.StatusOK, h.toResponse(c.Request.Context(), customer))
}

func cleanupDelete(storage storage.ObjectStorage, key string) {
	if storage == nil || key == "" {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := storage.Delete(ctx, key); err != nil {
		logger.ErrorKV(ctx, "failed to cleanup avatar object",
			"key", key,
			"error", err,
		)
	}
}
