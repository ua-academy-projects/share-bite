package auth

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func (s *service) Login(ctx context.Context, email, password string) (*Tokens, error) {
	u, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	at, rt, err := s.tokenManager.GeneratePair(u.ID)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	return &Tokens{AccessToken: at, RefreshToken: rt}, nil
}
