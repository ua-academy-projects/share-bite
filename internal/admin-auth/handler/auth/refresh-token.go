package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
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
