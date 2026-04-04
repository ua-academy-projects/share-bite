package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/repository/user"
	"github.com/ua-academy-projects/share-bite/internal/config"
	"github.com/ua-academy-projects/share-bite/pkg/database"
	"golang.org/x/crypto/bcrypt"
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
	userRepo  user.Repository
	txmanager database.TxManager
}

func New(userRepo user.Repository, txmanager database.TxManager) Service {
	return &service{userRepo: userRepo, txmanager: txmanager}
}

func (s *service) Login(ctx context.Context, email, password string) (*Tokens, error) {
	u, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return s.generateTokens(u.ID)
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
	var createdUserId string
	if txErr := s.txmanager.ReadCommited(ctx, func(txCtx context.Context) error {
		createdUser, err := s.userRepo.CreateUser(txCtx, user.CreateUser{
			Email:        email,
			PasswordHash: string(passwordHash),
		})
		if err != nil {
			return fmt.Errorf("create user: %w", err)
		}

		if err := s.userRepo.AssignRole(txCtx, createdUser.ID, role.ID); err != nil {
			return fmt.Errorf("assign role to user: %w", err)
		}
		createdUserId = createdUser.ID
		return nil
	}); txErr != nil {
		return nil, txErr
	}

	if createdUserId == "" {
		return nil, errors.New("created user not found")
	}

	return s.generateTokens(createdUserId)
}

func (s *service) Refresh(ctx context.Context, refreshToken string) (*Tokens, error) {
	cfg := config.Config().JwtToken

	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(cfg.RefreshTokenSecretKey()), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return nil, errors.New("invalid token subject")
	}

	return s.generateTokens(userID)
}

func (s *service) generateTokens(userID string) (*Tokens, error) {
	cfg := config.Config().JwtToken

	accessToken, err := s.generateToken(userID, cfg.AccessTokenSecretKey(), cfg.AccessTokenTTL())
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := s.generateToken(userID, cfg.RefreshTokenSecretKey(), cfg.RefreshTokenTTL())
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	return &Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) generateToken(userID, secretKey string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(ttl).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secretKey))
}
