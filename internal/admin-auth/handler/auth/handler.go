package auth

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	authsvc "github.com/ua-academy-projects/share-bite/internal/admin-auth/service/auth"
)

type authService interface {
	Login(ctx context.Context, email, password string) (*authsvc.Tokens, error)
	Register(ctx context.Context, email, password, slug string) (*authsvc.Tokens, error)
	Refresh(ctx context.Context, refreshToken string) (*authsvc.Tokens, error)
}
type handler struct {
	service authService
}

func RegisterHandlers(r *gin.RouterGroup, service authService) {
	h := &handler{service: service}

	r.POST("/login", h.login)
	r.POST("/register", h.register)
	r.POST("/refresh", h.refresh)
}

type loginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type registerRequest struct {
	Email    string `json:"email" binding:"required,email,max=254" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=8,max=72" example:"strong-password-123"`
	Slug     string `json:"slug" binding:"required,oneof=user business" enums:"user,business" example:"user"`
}

type authTokensResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
}

type errorResponse struct {
	Message string `json:"message" example:"invalid request payload"`
}

func (h *handler) login(c *gin.Context) {
	var req loginRequest
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

// register godoc
//
//	@Summary		Register admin user
//	@Description	Create a new admin-auth account and return access and refresh tokens
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		registerRequest		true	"Registration payload"
//	@Success		201		{object}	authTokensResponse
//	@Failure		400		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/auth/register [post]
func (h *handler) register(c *gin.Context) {
	var req registerRequest
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
