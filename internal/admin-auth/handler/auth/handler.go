package auth

import (
	"context"

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

type tokenResponce struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type errResponse struct {
	ErrorMessage string `json:"error"`
}

func RegisterHandlers(r *gin.RouterGroup, service authService) {
	h := &handler{service: service}

	r.POST("/login", h.login)
	r.POST("/register", h.register)
	r.POST("/refresh", h.refresh)
}
