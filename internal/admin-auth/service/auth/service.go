package auth

import (
	"context"
	"errors"
	"net/http"

	"time"

	"github.com/ua-academy-projects/share-bite/internal/admin-auth/dto"
	apperr "github.com/ua-academy-projects/share-bite/internal/admin-auth/error"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/pkg"
	emailsvc "github.com/ua-academy-projects/share-bite/internal/admin-auth/provider/email"
	"github.com/ua-academy-projects/share-bite/pkg/logger"

	"github.com/ua-academy-projects/share-bite/internal/admin-auth/repository/user"
	"github.com/ua-academy-projects/share-bite/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type TokenProvider interface {
	GenerateToken(userID string, role string) (string, string, error)
	ParseRefreshToken(token string) (string, string, error)
	GetRefreshTTL() time.Duration
}

type OAuthProvider interface {
	ExchangeCode(ctx context.Context, code string) (*dto.OAuthUserInfo, error)
}

type Service interface {
	Login(ctx context.Context, email, password string) (*Tokens, error)
	Register(ctx context.Context, email, password, slug string) (*Tokens, error)
	Refresh(ctx context.Context, refreshToken string) (*Tokens, error)
	Logout(ctx context.Context, userID, refreshToken string) error
	RevokeAllSessions(ctx context.Context, userID string) error
	OAuthLogin(ctx context.Context, provider OAuthProvider, code string, slug string) (*Tokens, error)
	LinkProvider(ctx context.Context, userID string, provider OAuthProvider, code string) error
	RecoverAccess(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
}

type service struct {
	userRepo         user.Repository
	tokenProvider    TokenProvider
	emailSender      emailsvc.Sender
	txManager        database.TxManager
	passwordResetTTL time.Duration
	maxSessions      int
}

func New(userRepo user.Repository, tokenProvider TokenProvider, emailSender emailsvc.Sender, txManager database.TxManager, resetTTL time.Duration, maxSessions int) Service {
	return &service{
		userRepo:         userRepo,
		tokenProvider:    tokenProvider,
		emailSender:      emailSender,
		txManager:        txManager,
		passwordResetTTL: resetTTL,
		maxSessions:      maxSessions,
	}
}

func (s *service) Login(ctx context.Context, email, password string) (*Tokens, error) {
	u, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, apperr.Wrap(http.StatusInternalServerError, "failed to query user", err)
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

	return s.issueTokens(ctx, u.ID, u.RoleSlug)
}

func (s *service) Register(ctx context.Context, email, password, slug string) (*Tokens, error) {
	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, apperr.Wrap(http.StatusInternalServerError, "failed to query user by email", err)
	}
	if existingUser != nil {
		return nil, apperr.ErrUserAlreadyExists
	}

	role, err := s.userRepo.FindRoleBySlug(ctx, slug)
	if err != nil {
		return nil, apperr.Wrap(http.StatusInternalServerError, "failed to query role", err)
	}
	if role == nil {
		return nil, apperr.ErrRoleNotFound
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperr.Wrap(http.StatusInternalServerError, "failed to hash password", err)
	}
	var createdUserId string
	if txErr := s.txManager.ReadCommitted(ctx, func(txCtx context.Context) error {
		createdUser, err := s.userRepo.CreateUser(txCtx, user.CreateUser{
			Email:        email,
			PasswordHash: string(passwordHash),
		})
		if err != nil {
			return apperr.Wrap(http.StatusInternalServerError, "failed to create user", err)
		}

		if err := s.userRepo.AssignRole(txCtx, createdUser.ID, role.ID); err != nil {
			return apperr.Wrap(http.StatusInternalServerError, "failed to assign role", err)
		}
		createdUserId = createdUser.ID
		return nil
	}); txErr != nil {
		return nil, txErr
	}

	if createdUserId == "" {
		return nil, apperr.Wrap(http.StatusInternalServerError, "created user not found after transaction", nil)
	}

	return s.issueTokens(ctx, createdUserId, role.Slug)
}

func (s *service) Refresh(ctx context.Context, refreshToken string) (*Tokens, error) {
	_, role, err := s.tokenProvider.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, apperr.ErrInvalidToken
	}

	hashedToken := pkg.HashToken(refreshToken)
	userID, err := s.userRepo.GetUserIDByRefreshToken(ctx, hashedToken)
	if err != nil {
		if errors.Is(err, apperr.ErrInvalidToken) {
			return nil, err
		}
		return nil, apperr.Wrap(http.StatusInternalServerError, "failed to verify refresh token in db", err)
	}
	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, apperr.Wrap(http.StatusInternalServerError, "failed to fetch user", err)
	}
	if u == nil {
		return nil, apperr.ErrUserNotFound
	}

	if err := s.userRepo.RevokeRefreshToken(ctx, hashedToken); err != nil {
		return nil, apperr.Wrap(http.StatusInternalServerError, "failed to revoke old token", err)
	}

	return s.issueTokens(ctx, userID, role)
}

func (s *service) Logout(ctx context.Context, userID string, refreshToken string) error {
	tokenHash := pkg.HashToken(refreshToken)

	ownerID, err := s.userRepo.GetUserIDByRefreshToken(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, apperr.ErrInvalidToken) {
			return nil
		}
		return err
	}
	if ownerID != userID {
		return apperr.ErrForbidden
	}

	err = s.userRepo.RevokeRefreshToken(ctx, tokenHash)
	if err != nil {
		return apperr.Wrap(http.StatusInternalServerError, "failed to logout", err)
	}
	return nil
}

func (s *service) RevokeAllSessions(ctx context.Context, userID string) error {
	err := s.userRepo.RevokeAllUserTokens(ctx, userID)
	if err != nil {
		return apperr.Wrap(http.StatusInternalServerError, "failed to revoke all sessions", err)
	}
	return nil
}

func (s *service) OAuthLogin(ctx context.Context, provider OAuthProvider, code string, slug string) (*Tokens, error) {
	info, err := provider.ExchangeCode(ctx, code)
	if err != nil {
		return nil, apperr.Wrap(http.StatusBadGateway, "failed to exchange code with provider", err)
	}

	existing, err := s.userRepo.FindBySocialProvider(ctx, info.Provider, info.ProviderID)
	if err != nil {
		return nil, apperr.Wrap(http.StatusInternalServerError, "failed to query social provider", err)
	}
	if existing != nil {
		return s.issueTokens(ctx, existing.ID, existing.RoleSlug)
	}

	byEmail, err := s.userRepo.FindByEmail(ctx, info.Email)
	if err != nil {
		return nil, apperr.Wrap(http.StatusInternalServerError, "failed to query user by email", err)
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
			return nil, apperr.Wrap(http.StatusInternalServerError, "failed to link social account", err)
		}
		return s.issueTokens(ctx, byEmail.ID, byEmail.RoleSlug)
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
		return nil, apperr.Wrap(http.StatusInternalServerError, "failed to create user via social", err)
	}
	return s.issueTokens(ctx, createUser.ID, role.Slug)
}

func (s *service) LinkProvider(ctx context.Context, userID string, provider OAuthProvider, code string) error {
	info, err := provider.ExchangeCode(ctx, code)
	if err != nil {
		return apperr.Wrap(http.StatusBadGateway, "failed to exchange code with provider", err)
	}

	existing, err := s.userRepo.FindBySocialProvider(ctx, info.Provider, info.ProviderID)
	if err != nil {
		return apperr.Wrap(http.StatusInternalServerError, "failed to check social provider", err)
	}
	if existing != nil {
		return apperr.ErrProviderAlreadyLinked
	}

	err = s.userRepo.LinkSocialAccount(ctx, dto.CreateSocialAccountParams{
		UserID:     userID,
		Provider:   info.Provider,
		ProviderID: info.ProviderID,
		Email:      info.Email,
	})
	if err != nil {
		return apperr.Wrap(http.StatusInternalServerError, "failed to link provider", err)
	}
	return nil
}

func (s *service) RecoverAccess(ctx context.Context, email string) error {
	u, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return apperr.Wrap(http.StatusInternalServerError, "failed to fetch user", err)
	}

	if u == nil {
		return apperr.UserNotFoundEmail(email)
	}

	rawToken, tokenHash, err := pkg.GeneratePasswordResetToken()
	if err != nil {
		return apperr.Wrap(http.StatusInternalServerError, "failed to generate reset token", err)
	}

	err = s.txManager.ReadCommitted(ctx, func(txCtx context.Context) error {
		if err := s.userRepo.CreatePasswordResetToken(txCtx, dto.CreatePasswordResetTokenParams{
			UserID:    u.ID,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().Add(s.passwordResetTTL),
		}); err != nil {
			return apperr.Wrap(http.StatusInternalServerError, "failed to store reset token", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	if err := s.emailSender.SendPasswordResetToken(ctx, u.Email, rawToken); err != nil {
		return apperr.Wrap(http.StatusInternalServerError, "failed to send password reset email", err)
	}

	return nil
}

func (s *service) ResetPassword(ctx context.Context, token, newPassword string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return apperr.Wrap(http.StatusInternalServerError, "failed to hash new password", err)
	}

	err = s.txManager.ReadCommitted(ctx, func(txCtx context.Context) error {
		userID, updated, err := s.userRepo.ResetPassword(txCtx, pkg.HashToken(token), string(passwordHash))
		if err != nil {
			return apperr.Wrap(http.StatusInternalServerError, "failed to reset password in db", err)
		}

		if !updated {
			return apperr.ErrInvalidResetToken
		}

		if err := s.userRepo.RevokeAllUserTokens(txCtx, userID); err != nil {
			return apperr.Wrap(http.StatusInternalServerError, "failed to revoke existing sessions during password reset", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *service) issueTokens(ctx context.Context, userID, role string) (*Tokens, error) {
	if err := s.userRepo.EnforceMaxSessions(ctx, userID, s.maxSessions); err != nil {
		logger.ErrorKV(ctx, "failed to enforce max sessions limit", err.Error())
	}
	access, refresh, err := s.tokenProvider.GenerateToken(userID, role)
	if err != nil {
		return nil, apperr.Wrap(http.StatusInternalServerError, "failed to generate jwt tokens", err)
	}
	ttl := s.tokenProvider.GetRefreshTTL()
	err = s.userRepo.StoreRefreshToken(ctx, dto.StoreRefreshTokenParams{
		TokenHash: pkg.HashToken(refresh),
		UserID:    userID,
		ExpiresAt: time.Now().Add(ttl),
	})
	if err != nil {
		return nil, apperr.Wrap(http.StatusInternalServerError, "failed to save refresh token to db", err)
	}

	return &Tokens{AccessToken: access, RefreshToken: refresh}, nil
}
