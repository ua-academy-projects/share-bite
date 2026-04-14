package customer

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/storage"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	"io"
	"net/http"
	"time"
)

const (
	maxAvatarSizeBytes = 5 * 1024 * 1024
	fileSniffSizeBytes = 512
)

func (h *handler) create(c *gin.Context) {
	var req createRequest
	if err := request.BindJSON(c, &req); err != nil {
		c.Error(err)
		return
	}

	userID, err := httpctx.GetUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	in, err := createRequestToCreateCustomer(req, userID)
	if err != nil {
		c.Error(err)
		return
	}

	customerID, err := h.service.Create(ctx, in)
	if err != nil {
		c.Error(err)
		return
	}

	resp := createResponse{CustomerID: customerID}
	c.JSON(http.StatusCreated, resp)
}

type createRequest struct {
	UserName  string `json:"userName" binding:"required,alphanum,min=3,max=30"`
	FirstName string `json:"firstName" binding:"required,min=2,max=50"`
	LastName  string `json:"lastName" binding:"required,min=2,max=50"`

	Bio *string `json:"bio" binding:"omitempty,max=500"`
}

type createResponse struct {
	CustomerID string `json:"customerId"`
}

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

	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.Error(apperror.ErrImageRequired)
		return
	}

	if fileHeader.Size <= 0 {
		c.Error(apperror.BadRequest("image file is empty"))
		return
	}
	if fileHeader.Size > maxAvatarSizeBytes {
		c.Error(apperror.BadRequest("image file is too large"))
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
	if !isAllowedAvatarContentType(contentType) {
		c.Error(apperror.ErrUnsupportedImageType)
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

	ext := extensionFromContentType(contentType)
	if ext == "" {
		c.Error(apperror.ErrUnsupportedImageType)
		return
	}
	objectKey := generateAvatarKey(userID, ext)

	uploadedKey, err := h.storage.Upload(
		c.Request.Context(),
		objectKey,
		contentType,
		file,
	)
	if err != nil {
		c.Error(err)
		return
	}

	customer, err := h.service.Update(c.Request.Context(), entity.UpdateCustomer{
		UserID:          userID,
		AvatarObjectKey: &uploadedKey,
	})
	if err != nil {
		cleanupDelete(h.storage, uploadedKey)
		c.Error(err)
		return
	}

	if currentCustomer.AvatarObjectKey != nil && *currentCustomer.AvatarObjectKey != uploadedKey {
		go cleanupDelete(h.storage, *currentCustomer.AvatarObjectKey)
	}

	c.JSON(http.StatusOK, h.toResponse(customer))
}

func generateAvatarKey(userID string, ext string) string {
	return fmt.Sprintf("avatars/%s/%s.%s", userID, uuid.New().String(), ext)
}

func extensionFromContentType(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return "jpg"
	case "image/png":
		return "png"
	default:
		return ""
	}
}

func isAllowedAvatarContentType(contentType string) bool {
	switch contentType {
	case "image/jpeg", "image/png":
		return true
	default:
		return false
	}
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
