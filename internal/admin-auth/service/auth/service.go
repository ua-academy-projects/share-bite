package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/ua-academy-projects/share-bite/internal/admin-auth/dto"
	apperror "github.com/ua-academy-projects/share-bite/internal/admin-auth/error"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/repository/user"
	emailsvc "github.com/ua-academy-projects/share-bite/internal/admin-auth/service/email"
	"github.com/ua-academy-projects/share-bite/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

const passwordResetTokenTTL = 30 * time.Minute

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
	RecoverAccess(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
}

type service struct {
	userRepo      user.Repository
	tokenProvider TokenProvider
	emailSender   emailsvc.Sender
	txManager     database.TxManager
}

func New(userRepo user.Repository, tokenProvider TokenProvider, emailSender emailsvc.Sender, txManager database.TxManager) Service {
	return &service{
		userRepo:      userRepo,
		tokenProvider: tokenProvider,
		emailSender:   emailSender,
		txManager:     txManager,
	}
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

func (s *service) RecoverAccess(ctx context.Context, email string) error {
	u, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("find user by email: %w", err)
	}

	if u == nil {
		return fmt.Errorf("find user by email: empty user result")
	}

	rawToken, tokenHash, err := generatePasswordResetToken()
	if err != nil {
		return fmt.Errorf("generate password reset token: %w", err)
	}

	err = s.txManager.ReadCommitted(ctx, func(txCtx context.Context) error {
		if err := s.userRepo.CreatePasswordResetToken(txCtx, dto.CreatePasswordResetTokenParams{
			UserID:    u.ID,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(passwordResetTokenTTL),
		}); err != nil {
			return fmt.Errorf("create password reset token: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	if err := s.emailSender.SendPasswordResetToken(ctx, u.Email, rawToken); err != nil {
		return fmt.Errorf("send password reset email: %w", err)
	}

	return nil
}

func (s *service) ResetPassword(ctx context.Context, token, newPassword string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	err = s.txManager.ReadCommitted(ctx, func(txCtx context.Context) error {
		updated, err := s.userRepo.ResetPassword(txCtx, hashToken(token), string(passwordHash))
		if err != nil {
			return fmt.Errorf("reset password: %w", err)
		}

		if !updated {
			return apperror.ErrInvalidResetToken
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func generatePasswordResetToken() (string, string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}

	raw := base64.RawURLEncoding.EncodeToString(b)
	return raw, hashToken(raw), nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
