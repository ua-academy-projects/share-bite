package business

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

type updateLocationURI struct {
	ID int `uri:"id" binding:"required,min=1"`
}

type updateLocationRequest struct {
	Name        *string   `json:"name"`
	Avatar      *string   `json:"avatar"`
	Banner      *string   `json:"banner"`
	Description *string   `json:"description"`
	Latitude    *float32  `json:"latitude" binding:"omitempty,gte=-90,lte=90"`
	Longitude   *float32  `json:"longitude" binding:"omitempty,gte=-180,lte=180"`
	TagSlugs    *[]string `json:"tagSlugs"`
}

// updateLocation updates an existing venue/location.
//
//	@Summary		Update location
//	@Description	Updates location fields (business owner only).
//	@Tags			locations
//	@Accept			json
//	@Produce		json
//	@Param			id				path		int						true	"Location ID"
//	@Param			Authorization	header		string					true	"Bearer access token"
//	@Param			request			body		updateLocationRequest	true	"Update location payload"
//	@Success		200				{object}	locationResponse
//	@Failure		400				{object}	errorResponse
//	@Failure		401				{object}	errorResponse
//	@Failure		403				{object}	errorResponse
//	@Failure		404				{object}	errorResponse
//	@Failure		500				{object}	errorResponse
//	@Router			/business/locations/{id} [patch]
func (h *handler) updateLocation(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.Error(apperror.Unauthorized("unauthorized"))
		return
	}

	var uri updateLocationURI
	if err := c.ShouldBindUri(&uri); err != nil {
		logger.ErrorKV(ctx, "failed to bind update location uri", "error", err)
		c.Error(apperror.BadRequest("invalid location id"))
		return
	}

	var req updateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorKV(ctx, "failed to bind update location request", "locationId", uri.ID, "error", err)
		c.Error(apperror.BadRequest("invalid request"))
		return
	}

	if req.Name == nil && req.Avatar == nil && req.Banner == nil &&
		req.Description == nil && req.Latitude == nil && req.Longitude == nil &&
		req.TagSlugs == nil {
		c.Error(apperror.BadRequest("empty update"))
		return
	}

	location, err := h.service.UpdateLocation(ctx, uri.ID, userID, dto.UpdateLocationInput{
		Name:        req.Name,
		Avatar:      req.Avatar,
		Banner:      req.Banner,
		Description: req.Description,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		TagSlugs:    req.TagSlugs,
	})
	if err != nil {
		sum := sha256.Sum256([]byte(userID))
		actorHash := hex.EncodeToString(sum[:])[:12]

		logger.ErrorKV(ctx, "failed to update location", "locationId", uri.ID, "actorHash", actorHash, "error", err)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, toLocationResponse(location))
}
