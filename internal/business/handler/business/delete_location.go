package business

import (
	"net/http"
	"crypto/sha256"
	"encoding/hex"

	"github.com/gin-gonic/gin"
	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

type deleteLocationURI struct {
	ID int `uri:"id" binding:"required,min=1"`
}

// deleteLocation deletes a venue/location.
//
//	@Summary		Delete location
//	@Description	Deletes a location under the owner's brand (business owner only).
//	@Tags			locations
//	@Produce		json
//	@Param			id				path		int		true	"Location ID"
//	@Param			Authorization	header		string	true	"Bearer access token"
//	@Success		204				{string}	string	""
//	@Failure		400				{object}	errorResponse
//	@Failure		401				{object}	errorResponse
//	@Failure		404				{object}	errorResponse
//	@Failure		403				{object}	errorResponse
//	@Failure		500				{object}	errorResponse
//	@Router			/business/locations/{id} [delete]
func (h *handler) deleteLocation(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.Error(apperror.Unauthorized("unauthorized"))
		return
	}

	var uri deleteLocationURI
	if err := c.ShouldBindUri(&uri); err != nil {
		logger.ErrorKV(ctx, "failed to bind delete location uri", "error", err)
		c.Error(apperror.BadRequest("invalid location id"))
		return
	}

	if err := h.service.DeleteLocation(ctx, uri.ID, userID); err != nil {
		sum := sha256.Sum256([]byte(userID))
		actorHash := hex.EncodeToString(sum[:])[:12]
		logger.ErrorKV(ctx, "failed to delete location", "locationId", uri.ID, "actorHash", actorHash, "error", err)
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
