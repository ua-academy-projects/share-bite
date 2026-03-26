package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/admin-auth/dto"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/repository/user"
	"golang.org/x/crypto/bcrypt"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type TokenProvider interface {
	GenerateToken(userID string, role string) (string, string, error)
	ParseRefreshToken(token string) (string, string, error)
}

type Service interface {
	Login(ctx context.Context, email, password string) (*Tokens, error)
	Register(ctx context.Context, email, password, slug string) (*Tokens, error)
	Refresh(ctx context.Context, refreshToken string) (*Tokens, error)
}

type service struct {
	userRepo      user.Repository
	tokenProvider TokenProvider
}

func New(userRepo user.Repository, tokenProvider TokenProvider) Service {
	return &service{userRepo: userRepo, tokenProvider: tokenProvider}
}

func (s *service) Login(ctx context.Context, email, password string) (*Tokens, error) {
	u, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	if u == nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	accessToken, refreshToken, err := s.tokenProvider.GenerateToken(u.ID, u.RoleSlug)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	return &Tokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

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

	createdUser, err := s.userRepo.CreateWithRole(ctx, dto.CreateWithRoleParams{
		Email:        email,
		PasswordHash: string(passwordHash),
		RoleID:       role.ID,
	})

	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	accessToken, refreshToken, err := s.tokenProvider.GenerateToken(createdUser.ID, role.Slug)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	return &Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func (s *service) Refresh(_ context.Context, refreshToken string) (*Tokens, error) {
	userID, role, err := s.tokenProvider.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	newAccessToken, newRefreshToken, err := s.tokenProvider.GenerateToken(userID, role)
	if err != nil {
		return nil, fmt.Errorf("generate new tokens: %w", err)
	}

	return &Tokens{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
