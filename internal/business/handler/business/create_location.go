package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
)

type createLocationURI struct {
	BrandID int `uri:"id" binding:"required,min=1"`
}

type createLocationRequest struct {
	Name        string   `json:"name" binding:"required,min=2,max=255"`
	Avatar      *string  `json:"avatar"`
	Banner      *string  `json:"banner"`
	Description *string  `json:"description"`
	Latitude    *float32 `json:"latitude"`
	Longitude   *float32 `json:"longitude"`
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
	})
	if err != nil {
		logger.ErrorKV(ctx, "failed to create location", "brandId", uri.BrandID, "ownerUserID", userID, "error", err)
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, toLocationResponse(location))
}
