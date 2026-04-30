package ghAuth

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

func (s *JWTSessionStore) Create(ctx context.Context, userID string) (string, error) {
	accessToken, _, err := s.tokenManager.GenerateToken(userID, "admin")
	if err != nil {
		return "", fmt.Errorf("generate session token: %w", err)
	}
	return accessToken, nil
}