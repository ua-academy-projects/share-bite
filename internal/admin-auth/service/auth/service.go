package auth

import (
	"context"
	"net/http"

	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ua-academy-projects/share-bite/internal/admin-auth/dto"
	apperr "github.com/ua-academy-projects/share-bite/internal/admin-auth/error"
	emailsvc "github.com/ua-academy-projects/share-bite/internal/admin-auth/provider/email"

	"github.com/ua-academy-projects/share-bite/internal/admin-auth/repository/user"
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

type OAuthProvider interface {
	ExchangeCode(ctx context.Context, code string) (*dto.OAuthUserInfo, error)
}

type Service interface {
	Login(ctx context.Context, email, password string) (*Tokens, error)
	Register(ctx context.Context, email, password, slug string) (*Tokens, error)
	Refresh(ctx context.Context, refreshToken string) (*Tokens, error)
	OAuthLogin(ctx context.Context, provider OAuthProvider, code string, slug string) (*Tokens, error)
	LinkProvider(ctx context.Context, userID string, provider OAuthProvider, code string) error
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
		return nil, apperr.ErrInvalidCredentials
	}
	if u.PasswordHash == nil {
		return nil, apperr.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*u.PasswordHash), []byte(password)); err != nil {
		return nil, apperr.ErrInvalidCredentials
	}

	return s.issueTokens(u.ID, u.RoleSlug)
}

func (s *service) Register(ctx context.Context, email, password, slug string) (*Tokens, error) {
	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	if existingUser != nil {
		return nil, apperr.ErrUserAlreadyExists
	}

	role, err := s.userRepo.FindRoleBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("find role by slug: %w", err)
	}
	if role == nil {
		return nil, apperr.ErrRoleNotFound
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
	return s.issueTokens(createdUser.ID, role.Slug)
}

func (s *service) Refresh(_ context.Context, refreshToken string) (*Tokens, error) {
	userID, role, err := s.tokenProvider.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, apperr.ErrInvalidToken
	}

	return s.issueTokens(userID, role)
}

func (s *service) OAuthLogin(ctx context.Context, provider OAuthProvider, code string, slug string) (*Tokens, error) {
	info, err := provider.ExchangeCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchange code: %w", err)
	}

	existing, err := s.userRepo.FindBySocialProvider(ctx, info.Provider, info.ProviderID)
	if err != nil {
		return nil, fmt.Errorf("find by social provider: %w", err)
	}
	if existing != nil {
		return s.issueTokens(existing.ID, existing.RoleSlug)
	}

	byEmail, err := s.userRepo.FindByEmail(ctx, info.Email)
	if err != nil {
		return nil, fmt.Errorf("find by email: %w", err)
	}
	if byEmail != nil {
		if !info.EmailVerified {
			return nil, apperr.ErrEmailNotVerified
		}
		err := s.userRepo.LinkSocialAccount(ctx, dto.CreateSocialAccountParams{
			UserID:     byEmail.ID,
			Provider:   info.Provider,
			ProviderID: info.ProviderID,
			Email:      info.Email,
		})
		if err != nil {
			return nil, err
		}
		return s.issueTokens(byEmail.ID, byEmail.RoleSlug)
	}
	if !info.EmailVerified {
		return nil, apperr.ErrEmailNotVerified
	}

	role, err := s.userRepo.FindRoleBySlug(ctx, slug)
	if err != nil {
		return nil, apperr.Wrap(http.StatusInternalServerError, "failed to find role by slug", err)
	}
	if role == nil {
		return nil, apperr.ErrRoleNotFound
	}
	createUser, err := s.userRepo.CreateWithSocial(ctx, dto.CreateUserWithSocialParams{
		Email:      info.Email,
		Provider:   info.Provider,
		ProviderID: info.ProviderID,
		RoleID:     role.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return s.issueTokens(createUser.ID, role.Slug)
}

func (s *service) LinkProvider(ctx context.Context, userID string, provider OAuthProvider, code string) error {
	info, err := provider.ExchangeCode(ctx, code)
	if err != nil {
		return fmt.Errorf("exchange code: %w", err)
	}

	existing, err := s.userRepo.FindBySocialProvider(ctx, info.Provider, info.ProviderID)
	if err != nil {
		return fmt.Errorf("find by social provider: %w", err)
	}
	if existing != nil {
		return apperr.ErrProviderAlreadyLinked
	}

	return s.userRepo.LinkSocialAccount(ctx, dto.CreateSocialAccountParams{
		UserID:     userID,
		Provider:   info.Provider,
		ProviderID: info.ProviderID,
		Email:      info.Email,
	})
}

func (s *service) issueTokens(userID, role string) (*Tokens, error) {
	access, refresh, err := s.tokenProvider.GenerateToken(userID, role)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}
	return &Tokens{AccessToken: access, RefreshToken: refresh}, nil
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

	err = s.txManager.ReadCommited(ctx, func(txCtx context.Context) error {
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

	err = s.txManager.ReadCommited(ctx, func(txCtx context.Context) error {
		updated, err := s.userRepo.ResetPassword(txCtx, hashToken(token), string(passwordHash))
		if err != nil {
			return fmt.Errorf("reset password: %w", err)
		}

		if !updated {
			return apperr.ErrInvalidResetToken
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
