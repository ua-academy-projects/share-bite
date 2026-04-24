package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	authsvc "github.com/ua-academy-projects/share-bite/internal/admin-auth/service/auth"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

type ProviderFactory interface {
	Get(name string) (authsvc.OAuthProvider, error)
}
type Handler struct {
	service         authsvc.Service
	providerFactory ProviderFactory
}

func NewHandler(service authsvc.Service, providerFactory ProviderFactory) *Handler {
	return &Handler{service: service, providerFactory: providerFactory}
}

// Login godoc
// @Summary      User login
// @Description  Authenticates a user using email and password.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      LoginRequest  true  "Login credentials"
// @Success      200      {object}  TokensResponse  "Success. Returns access and refresh tokens."
// @Failure      400      {object}  ErrorResponse   "Validation error."
// @Failure      401      {object}  ErrorResponse   "Invalid credentials."
// @Failure      500      {object}  ErrorResponse   "Internal server error."
// @Router       /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	tokens, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, TokensResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

// RecoverAccess godoc
// @Summary      Recover access
// @Description  Sends password reset instructions if the account exists.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      RecoverAccessRequest  true  "Recover access payload"
// @Success      200      {object}  MessageResponse       "Success message."
// @Failure      400      {object}  ErrorResponse         "Validation error."
// @Failure      500      {object}  ErrorResponse         "Internal server error."
// @Router       /auth/recover-access [post]
func (h *Handler) RecoverAccess(c *gin.Context) {
	var req RecoverAccessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.service.RecoverAccess(c.Request.Context(), req.Email)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, MessageResponse{
		Message: "If the email exists, recovery instructions have been sent",
	})
}

// ResetPassword godoc
// @Summary      Reset password
// @Description  Resets the password using a valid token and revokes all old sessions.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      ResetPasswordRequest  true  "Reset password payload"
// @Success      200      {object}  MessageResponse       "Success message."
// @Failure      400      {object}  ErrorResponse         "Validation error."
// @Failure      500      {object}  ErrorResponse         "Internal server error."
// @Router       /auth/reset-password [post]
func (h *Handler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.service.ResetPassword(c.Request.Context(), req.Token, req.NewPassword)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "Password has been reset successfully."})
}

// Register godoc
// @Summary      User registration
// @Description  Creates a new user account and immediately returns a pair of tokens.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      RegisterRequest  true  "Registration payload"
// @Success      201      {object}  TokensResponse   "Success. Returns access and refresh tokens."
// @Failure      400      {object}  ErrorResponse    "Validation error."
// @Failure      409      {object}  ErrorResponse    "User already exists."
// @Failure      422      {object}  ErrorResponse    "Role not found."
// @Failure      500      {object}  ErrorResponse    "Internal server error."
// @Router       /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	tokens, err := h.service.Register(c.Request.Context(), req.Email, req.Password, req.Slug)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, TokensResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

// Refresh godoc
// @Summary      Refresh tokens
// @Description  Generates a new pair of tokens based on a valid refresh_token (Token Rotation).
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      RefreshRequest  true  "Refresh token payload"
// @Success      200      {object}  TokensResponse  "Success. Returns new access and refresh tokens."
// @Failure      400      {object}  ErrorResponse   "Validation error."
// @Failure      401      {object}  ErrorResponse   "Invalid or expired refresh token."
// @Failure      500      {object}  ErrorResponse   "Internal server error."
// @Router       /auth/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	tokens, err := h.service.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, TokensResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

// OAuthCallback godoc
// @Summary      OAuth Login / Registration
// @Description  Exchanges an authorization code from a provider for access tokens.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        provider path      string                true  "Provider name (e.g., google)"
// @Param        request  body      OAuthCallbackRequest  true  "OAuth callback payload"
// @Success      200      {object}  TokensResponse        "Success. Returns access and refresh tokens."
// @Failure      400      {object}  ErrorResponse         "Unsupported provider or invalid request."
// @Failure      422      {object}  ErrorResponse         "Role not found."
// @Failure      500      {object}  ErrorResponse         "Internal server error."
// @Failure      502      {object}  ErrorResponse         "Error exchanging code with the provider."
// @Router       /auth/oauth/{provider}/callback [post]
func (h *Handler) OAuthCallback(c *gin.Context) {
	ctx := c.Request.Context()
	providerName := c.Param("provider")

	p, err := h.providerFactory.Get(providerName)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var req OAuthCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload."})
		return
	}

	tokens, err := h.service.OAuthLogin(ctx, p, req.Code, req.Slug)
	if err != nil {
		_ = c.Error(err)
		return
	}

	logger.InfoKV(ctx, "user oauth login success", "provider", providerName)

	c.JSON(http.StatusOK, TokensResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

// OAuthLinkAccount godoc
// @Summary      Link social account
// @Description  Links a Google or GitHub account to an already authenticated user.
// @Tags         User
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        provider path      string            true  "Provider name (e.g., google)"
// @Param        request  body      OAuthLinkRequest  true  "OAuth link payload"
// @Success      200      {object}  MessageResponse   "Success. Returns a confirmation message."
// @Failure      400      {object}  ErrorResponse     "Invalid request."
// @Failure      401      {object}  ErrorResponse     "Unauthorized access."
// @Failure      409      {object}  ErrorResponse     "Provider is already linked to the account."
// @Failure      500      {object}  ErrorResponse     "Internal server error."
// @Failure      502      {object}  ErrorResponse     "Error exchanging code with the provider."
// @Router       /user/link/{provider} [post]
func (h *Handler) OAuthLinkAccount(c *gin.Context) {
	ctx := c.Request.Context()

	userIDVal, exists := c.Get(middleware.CtxUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized access."})
		return
	}
	userID, ok := userIDVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error."})
		return
	}

	providerName := c.Param("provider")
	p, err := h.providerFactory.Get(providerName)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var req OAuthLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload. 'code' is required."})
		return
	}

	err = h.service.LinkProvider(ctx, userID, p, req.Code)
	if err != nil {
		_ = c.Error(err)
		return
	}

	logger.InfoKV(ctx, "social account successfully linked", "user_id", userID, "provider", providerName)

	c.JSON(http.StatusOK, MessageResponse{Message: "Social account successfully linked."})
}

// Logout godoc
// @Summary      User logout
// @Description  Deletes a specific session (refresh token) from the database.
// @Tags         Auth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      RefreshRequest   true  "Refresh token to delete"
// @Success      200      {object}  MessageResponse  "Success message."
// @Failure      400      {object}  ErrorResponse    "Validation error."
// @Failure      401      {object}  ErrorResponse    "Unauthorized access."
// @Failure      500      {object}  ErrorResponse    "Internal server error."
// @Router       /auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	userIDVal, exists := c.Get(middleware.CtxUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized access."})
		return
	}

	userIDStr, ok := userIDVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error."})
		return
	}

	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.service.Logout(c.Request.Context(), userIDStr, req.RefreshToken)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "Successfully logged out."})
}

// RevokeAllSessions godoc
// @Summary      Revoke all sessions
// @Description  Invalidates all active sessions for the user (deletes all refresh tokens).
// @Tags         User
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200      {object}  MessageResponse  "Success message."
// @Failure      401      {object}  ErrorResponse    "Unauthorized access."
// @Failure      500      {object}  ErrorResponse    "Internal server error."
// @Router       /user/sessions/revoke-all [post]
func (h *Handler) RevokeAllSessions(c *gin.Context) {
	ctx := c.Request.Context()
	userIDVal, exists := c.Get(middleware.CtxUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized access."})
		return
	}
	userID, ok := userIDVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error."})
		return
	}
	err := h.service.RevokeAllSessions(ctx, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	logger.InfoKV(ctx, "all user sessions revoked", "user_id", userID)
	c.JSON(http.StatusOK, MessageResponse{Message: "All sessions have been successfully revoked."})
}
