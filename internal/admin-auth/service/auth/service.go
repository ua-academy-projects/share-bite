package auth

import (
	"context"

	"github.com/ua-academy-projects/share-bite/internal/admin-auth/repository/user"
	"github.com/ua-academy-projects/share-bite/pkg/jwt"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Service interface {
	Login(ctx context.Context, email, password string) (*Tokens, error)
	Register(ctx context.Context, email, password, slug string) (*Tokens, error)
	Refresh(ctx context.Context, refreshToken string) (*Tokens, error)
}

type service struct {
	userRepo     user.Repository
	tokenManager *jwt.TokenManager
}

func New(userRepo user.Repository, tokenManager *jwt.TokenManager) Service {
	return &service{
		userRepo:     userRepo,
		tokenManager: tokenManager,
	}
}
