package business

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/internal/storage/mediatype"
)

func (h *handler) uploadBanner(c *gin.Context) {
	reqURI := new(getRequest)
	if err := c.ShouldBindUri(reqURI); err != nil {
		c.Error(apperror.BadRequest("invalid id"))
		return
	}

	orgAccountID, err := h.extractUserUUID(c)
	if err != nil {
		c.Error(err)
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, mediatype.DefaultMaxImageSizeBytes+multipartOverheadBytes)
	file, err := c.FormFile("image")
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			c.Error(apperror.BadRequest("image file is too large"))
			return
		}
		c.Error(apperror.BadRequest("image is required"))
		return
	}

	updated, err := h.service.UploadBanner(c.Request.Context(), reqURI.ID, orgAccountID, file)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.UpdateOrgResponse{
		Id:          updated.Id,
		Name:        updated.Name,
		Avatar:      h.presign(c.Request.Context(), updated.Avatar),
		Banner:      h.presign(c.Request.Context(), updated.Banner),
		Description: updated.Description,
	})
}
