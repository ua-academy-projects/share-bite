package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

type getRequest struct {
	ID int `uri:"id" binding:"required"`
}

type brandResponse struct {
	Id     int     `json:"id" example:"1"`
	Name   string  `json:"name" example:"ShareBite HQ"`
	Avatar *string `json:"avatar" example:"https://cdn.example.com/brand-avatar.png"`
}

type getResponse struct {
	Id          int            `json:"id" example:"42"`
	Name        string         `json:"name" example:"ShareBite Downtown"`
	Avatar      *string        `json:"avatar" example:"https://cdn.example.com/avatar.png"`
	Banner      *string        `json:"banner" example:"https://cdn.example.com/banner.png"`
	Description *string        `json:"description" example:"A cozy place in the city center."`
	Latitude    *float32       `json:"latitude" example:"50.4501"`
	Longitude   *float32       `json:"longitude" example:"30.5234"`
	Brand       *brandResponse `json:"brand,omitempty"`
	Tags        []string       `json:"tags"`
}

// get returns a single location (venue) with its parent brand.
//
//	@Summary		Get location by ID
//	@Description	Returns a single venue/location that belongs to a brand, including the parent brand info.
//	@Tags			locations
//	@Produce		json
//	@Param			id	path		int	true	"Location ID"
//	@Success		200	{object}	getResponse
//	@Failure		400	{object}	errorResponse
//	@Failure		404	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/business/org-units/{id} [get]
func (h *handler) get(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	req := new(getRequest)
	if err := c.ShouldBindUri(req); err != nil || req.ID < 1 {
		c.Error(apperror.BadRequest("invalid location id"))
		return
	}

	log.Info("get location", "id", req.ID)

	location, err := h.service.Get(ctx, req.ID)
	if err != nil {
		log.Error("failed to get location", "id", req.ID, "error", err)
		c.Error(err)
		return
	}

	if location.ParentId == nil {
		log.Error("org unit is a brand, not a location", "id", req.ID)
		c.Error(apperror.OrgUnitNotFoundID(req.ID))
		return
	}

	resp := getResponse{
		Id:          location.Id,
		Name:        location.Name,
		Avatar:      location.Avatar,
		Banner:      location.Banner,
		Description: location.Description,
		Latitude:    location.Latitude,
		Longitude:   location.Longitude,
		Tags:        location.Tags,
	}

	if location.ParentId != nil {
		brand, err := h.service.Get(ctx, *location.ParentId)
		if err != nil {
			log.Error("failed to get brand for location", "brandId", *location.ParentId, "error", err)
			c.Error(err)
			return
		}
		resp.Brand = &brandResponse{
			Id:     brand.Id,
			Name:   brand.Name,
			Avatar: brand.Avatar,
		}
	}

	c.JSON(http.StatusOK, resp)
}
