package github

import (
	"context"
	"fmt"

	"github.com/ua-academy-projects/share-bite/pkg/jwt"
)

type JWTSessionStore struct {
	tokenManager *jwt.TokenManager
}

func NewJWTSessionStore(tokenManager *jwt.TokenManager) *JWTSessionStore {
	return &JWTSessionStore{
		tokenManager: tokenManager,
	}
}

func (s *JWTSessionStore) Create(ctx context.Context, userID, role string) (string, error) {
	if role == "" {
		role = "user"
	}
	accessToken, _, err := s.tokenManager.GenerateToken(userID, role, jwt.UserStatusActive)
	if err != nil {
		return "", fmt.Errorf("generate session token: %w", err)
	}
	return accessToken, nil
}
