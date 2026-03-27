package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) get(c *gin.Context) {
	req := new(getRequest)
	if err := c.ShouldBindUri(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid business id"})
		return
	}

	ctx := c.Request.Context()
	orgUnit, err := h.service.Get(ctx, req.ID)
	if err != nil {
		c.Error(err)
		return
	}

	parentUnit, err := h.service.Get(ctx, orgUnit.ParentId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, getResponse{
		Id:          orgUnit.Id,
		Name:        orgUnit.Name,
		Avatar:      orgUnit.Avatar,
		Banner:      orgUnit.Banner,
		Description: orgUnit.Description,
		Latitude:    orgUnit.Latitude,
		Longitude:   orgUnit.Longitude,
		Parent: struct {
			Id     int    `json:"id"`
			Name   string `json:"name"`
			Avatar string `json:"avatar"`
		}{
			Id:     parentUnit.Id,
			Name:   parentUnit.Name,
			Avatar: parentUnit.Avatar,
		},
	})
}

type getRequest struct {
	ID int `uri:"id" binding:"required"`
}

type getResponse struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Avatar      string  `json:"avatar"`
	Banner      string  `json:"banner"`
	Description string  `json:"description"`
	Latitude    float32 `json:"latitude"`
	Longitude   float32 `json:"longitude"`
	Parent      struct {
		Id     int    `json:"id"`
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
	} `json:"parent"`
}
