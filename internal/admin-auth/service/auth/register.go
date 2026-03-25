package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/admin-auth/repository/user"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) Register(ctx context.Context, email, password, slug string) (*Tokens, error) {
	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	role, err := s.userRepo.FindRoleBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("find role by slug: %w", err)
	}
	if role == nil {
		return nil, errors.New("role not found")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	createdUser, err := s.userRepo.Create(ctx, user.CreateParams{
		Email:        email,
		PasswordHash: string(passwordHash),
	})
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	if err := s.userRepo.AssignRole(ctx, createdUser.ID, role.ID); err != nil {
		return nil, fmt.Errorf("assign role: %w", err)
	}

	at, rt, err := s.tokenManager.GeneratePair(createdUser.ID)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	return &Tokens{AccessToken: at, RefreshToken: rt}, nil
}
