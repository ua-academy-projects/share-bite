package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) list(c *gin.Context) {
	req := new(listRequest)
	if err := c.ShouldBindQuery(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	orgUnits, err := h.service.List(ctx, req.Page, req.Limit)
	if err != nil {
		_ = c.Error(err)
		return
	}

	items := make([]listItem, 0, len(orgUnits))
	for _, u := range orgUnits {
		items = append(items, listItem{
			Id:          u.Id,
			Name:        u.Name,
			Avatar:      u.Avatar,
			Description: u.Description,
			Latitude:    u.Latitude,
			Longitude:   u.Longitude,
		})
	}

	c.JSON(http.StatusOK, listResponse{
		Items: items,
	})
}

type listRequest struct {
	Page  int `form:"page" binding:"required,gte=1"`
	Limit int `form:"limit" binding:"required,gte=1,lte=100"`
}

type listItem struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Avatar      string  `json:"avatar"`
	Description string  `json:"description"`
	Latitude    float32 `json:"latitude"`
	Longitude   float32 `json:"longitude"`
}

type listResponse struct {
	Items []listItem `json:"items"`
}
