package business

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type handler struct {
	service businessService
}

type businessService interface {
}

func RegisterHandlers(
	r *gin.RouterGroup,
	service businessService,
) {
	h := &handler{
		service: service,
	}

	r.POST("/:id", h.getById)
	r.POST("/list", h.List)
}

func (h *handler) refresh(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	tokens, err := h.service.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})
}
