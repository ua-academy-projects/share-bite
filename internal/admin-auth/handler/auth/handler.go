package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	authsvc "github.com/ua-academy-projects/share-bite/internal/admin-auth/service/auth"
)

type Handler struct {
	service authsvc.Service
}

func NewHandler(service authsvc.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	tokens, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})
}

// RecoverAccess godoc
// @Summary Recover access
// @Description Sends password reset instructions if account exists
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RecoverAccessRequest true "Recover access payload"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/recover-access [post]
func (h *Handler) RecoverAccess(c *gin.Context) {
	var req RecoverAccessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.service.RecoverAccess(c.Request.Context(), req.Email)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "If the email exists, recovery instructions have been sent"})
}

// ResetPassword godoc
// @Summary Reset password
// @Description Resets password by token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "Reset password payload"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/reset-password [post]
func (h *Handler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.service.ResetPassword(c.Request.Context(), req.Token, req.NewPassword)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "password has been reset"})
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	tokens, err := h.service.Register(
		c.Request.Context(),
		req.Email,
		req.Password,
		req.Slug,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})
}

func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest
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
