package business

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

type createLocationURI struct {
	BrandID int `uri:"id" binding:"required,min=1"`
}

type createLocationRequest struct {
	Name        string   `json:"name" binding:"required,min=2,max=255"`
	Avatar      *string  `json:"avatar"`
	Banner      *string  `json:"banner"`
	Description *string  `json:"description"`
	Latitude    *float32 `json:"latitude" binding:"omitempty,gte=-90,lte=90"`
	Longitude   *float32 `json:"longitude" binding:"omitempty,gte=-180,lte=180"`
	TagIDs      []int    `json:"tagIds"`
}

type locationResponse struct {
	Id           int      `json:"id"`
	OrgAccountId string   `json:"orgAccountId"`
	ProfileType  string   `json:"profileType"`
	Name         string   `json:"name"`
	Avatar       *string  `json:"avatar"`
	Banner       *string  `json:"banner"`
	Description  *string  `json:"description"`
	ParentId     *int     `json:"parentId,omitempty"`
	Latitude     *float32 `json:"latitude"`
	Longitude    *float32 `json:"longitude"`
	Tags         []string `json:"tags"`
}

func toLocationResponse(loc *entity.OrgUnit) locationResponse {
	return locationResponse{
		Id:           loc.Id,
		OrgAccountId: loc.OrgAccountId.String(),
		ProfileType:  loc.ProfileType,
		Name:         loc.Name,
		Avatar:       loc.Avatar,
		Banner:       loc.Banner,
		Description:  loc.Description,
		ParentId:     loc.ParentId,
		Latitude:     loc.Latitude,
		Longitude:    loc.Longitude,
		Tags:         normalizeTags(loc.Tags),
	}
}

// createLocation creates a new venue/location under a brand.
//
//	@Summary		Create location
//	@Description	Creates a venue/location under the specified brand (business owner only).
//	@Tags			locations
//	@Accept			json
//	@Produce		json
//	@Param			id				path		int						true	"Brand ID"
//	@Param			Authorization	header		string					true	"Bearer access token"
//	@Param			request			body		createLocationRequest	true	"Create location payload"
//	@Success		201				{object}	locationResponse
//	@Failure		400				{object}	errorResponse
//	@Failure		401				{object}	errorResponse
//	@Failure		403				{object}	errorResponse
//	@Failure		404				{object}	errorResponse
//	@Failure		500				{object}	errorResponse
//	@Router			/business/{id}/locations [post]
func (h *handler) createLocation(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.Error(apperror.Unauthorized("unauthorized"))
		return
	}

	var uri createLocationURI
	if err := c.ShouldBindUri(&uri); err != nil {
		logger.ErrorKV(ctx, "failed to bind create location uri", "error", err)
		c.Error(apperror.BadRequest("invalid brand id"))
		return
	}

	var req createLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorKV(ctx, "failed to bind create location request", "brandId", uri.BrandID, "error", err)
		c.Error(apperror.BadRequest("invalid request"))
		return
	}

	location, err := h.service.CreateLocation(ctx, uri.BrandID, userID, dto.CreateLocationInput{
		Name:        req.Name,
		Avatar:      req.Avatar,
		Banner:      req.Banner,
		Description: req.Description,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		TagIDs:      req.TagIDs,
	})
	if err != nil {
		sum := sha256.Sum256([]byte(userID))
		actorHash := hex.EncodeToString(sum[:])[:12]

		logger.ErrorKV(ctx, "failed to create location", "brandId", uri.BrandID, "actorHash", actorHash, "error", err)
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, toLocationResponse(location))
}

func normalizeTags(tags []string) []string {
	if tags == nil {
		return []string{}
	}
	return tags
}
