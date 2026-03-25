package auth

import (
	"context"
	"fmt"
)

func (s *service) Refresh(ctx context.Context, refreshToken string) (*Tokens, error) {
	userID, err := s.tokenManager.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	at, rt, err := s.tokenManager.GeneratePair(userID)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	return &Tokens{AccessToken: at, RefreshToken: rt}, nil
}
