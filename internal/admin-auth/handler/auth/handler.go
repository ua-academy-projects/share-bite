package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/handler"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/models"
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
// @Param        request  body      handler.LoginRequest    true  "Login credentials"
// @Success      200      {object}  handler.TokensResponse  "Success. Returns access and refresh tokens."
// @Failure      400      {object}  handler.ErrorResponse   "Validation error."
// @Failure      401      {object}  handler.ErrorResponse   "Invalid credentials."
// @Failure      500      {object}  handler.ErrorResponse   "Internal server error."
// @Router       /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req handler.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, handler.ErrorResponse{Error: err.Error()})
		return
	}

	tokens, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, handler.TokensResponse{
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
// @Param        request  body      handler.RecoverAccessRequest  true  "Recover access payload"
// @Success      200      {object}  handler.MessageResponse       "Success message."
// @Failure      400      {object}  handler.ErrorResponse         "Validation error."
// @Failure      500      {object}  handler.ErrorResponse         "Internal server error."
// @Router       /auth/recover-access [post]
func (h *Handler) RecoverAccess(c *gin.Context) {
	var req handler.RecoverAccessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, handler.ErrorResponse{Error: err.Error()})
		return
	}

	err := h.service.RecoverAccess(c.Request.Context(), req.Email)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, handler.MessageResponse{
		Message: "If the email exists, recovery instructions have been sent",
	})
}

// ResetPassword godoc
// @Summary      Reset password
// @Description  Resets the password using a valid token and revokes all old sessions.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      handler.ResetPasswordRequest  true  "Reset password payload"
// @Success      200      {object}  handler.MessageResponse       "Success message."
// @Failure      400      {object}  handler.ErrorResponse         "Validation error."
// @Failure      500      {object}  handler.ErrorResponse         "Internal server error."
// @Router       /auth/reset-password [post]
func (h *Handler) ResetPassword(c *gin.Context) {
	var req handler.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, handler.ErrorResponse{Error: err.Error()})
		return
	}

	err := h.service.ResetPassword(c.Request.Context(), req.Token, req.NewPassword)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, handler.MessageResponse{Message: "Password has been reset successfully."})
}

// Register godoc
// @Summary      User registration
// @Description  Creates a new user account and immediately returns a pair of tokens.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      handler.RegisterRequest  true  "Registration payload"
// @Success      201      {object}  handler.TokensResponse   "Success. Returns access and refresh tokens."
// @Failure      400      {object}  handler.ErrorResponse    "Validation error."
// @Failure      409      {object}  handler.ErrorResponse    "User already exists."
// @Failure      422      {object}  handler.ErrorResponse    "Role not found."
// @Failure      500      {object}  handler.ErrorResponse    "Internal server error."
// @Router       /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req handler.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, handler.ErrorResponse{Error: err.Error()})
		return
	}

	tokens, err := h.service.Register(c.Request.Context(), req.Email, req.Password, req.Slug)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, handler.TokensResponse{
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
// @Param        request  body      handler.RefreshRequest  true  "Refresh token payload"
// @Success      200      {object}  handler.TokensResponse  "Success. Returns new access and refresh tokens."
// @Failure      400      {object}  handler.ErrorResponse   "Validation error."
// @Failure      401      {object}  handler.ErrorResponse   "Invalid or expired refresh token."
// @Failure      500      {object}  handler.ErrorResponse   "Internal server error."
// @Router       /auth/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	var req handler.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, handler.ErrorResponse{Error: err.Error()})
		return
	}

	tokens, err := h.service.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, handler.TokensResponse{
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
// @Param        request  body      handler.OAuthCallbackRequest  true  "OAuth callback payload"
// @Success      200      {object}  handler.TokensResponse        "Success. Returns access and refresh tokens."
// @Failure      400      {object}  handler.ErrorResponse         "Unsupported provider or invalid request."
// @Failure      422      {object}  handler.ErrorResponse         "Role not found."
// @Failure      500      {object}  handler.ErrorResponse         "Internal server error."
// @Failure      502      {object}  handler.ErrorResponse         "Error exchanging code with the provider."
// @Router       /auth/oauth/{provider}/callback [post]
func (h *Handler) OAuthCallback(c *gin.Context) {
	ctx := c.Request.Context()
	providerName := c.Param("provider")

	p, err := h.providerFactory.Get(providerName)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var req handler.OAuthCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, handler.ErrorResponse{Error: "Invalid request payload."})
		return
	}

	tokens, err := h.service.OAuthLogin(ctx, p, req.Code, req.Slug)
	if err != nil {
		_ = c.Error(err)
		return
	}

	logger.InfoKV(ctx, "user oauth login success", "provider", providerName)

	c.JSON(http.StatusOK, handler.TokensResponse{
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
// @Param        request  body      handler.OAuthLinkRequest  true  "OAuth link payload"
// @Success      200      {object}  handler.MessageResponse   "Success. Returns a confirmation message."
// @Failure      400      {object}  handler.ErrorResponse     "Invalid request."
// @Failure      401      {object}  handler.ErrorResponse     "Unauthorized access."
// @Failure      409      {object}  handler.ErrorResponse     "Provider is already linked to the account."
// @Failure      500      {object}  handler.ErrorResponse     "Internal server error."
// @Failure      502      {object}  handler.ErrorResponse     "Error exchanging code with the provider."
// @Router       /user/link/{provider} [post]
func (h *Handler) OAuthLinkAccount(c *gin.Context) {
	ctx := c.Request.Context()

	userIDVal, exists := c.Get(middleware.CtxUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, handler.ErrorResponse{Error: "Unauthorized access."})
		return
	}
	userID, ok := userIDVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, handler.ErrorResponse{Error: "Internal server error."})
		return
	}

	providerName := c.Param("provider")
	p, err := h.providerFactory.Get(providerName)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var req handler.OAuthLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, handler.ErrorResponse{Error: "Invalid request payload. 'code' is required."})
		return
	}

	err = h.service.LinkProvider(ctx, userID, p, req.Code)
	if err != nil {
		_ = c.Error(err)
		return
	}

	logger.InfoKV(ctx, "social account successfully linked", "user_id", userID, "provider", providerName)

	c.JSON(http.StatusOK, handler.MessageResponse{Message: "Social account successfully linked."})
}

// Logout godoc
// @Summary      User logout
// @Description  Deletes a specific session (refresh token) from the database.
// @Tags         Auth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      handler.RefreshRequest   true  "Refresh token to delete"
// @Success      200      {object}  handler.MessageResponse  "Success message."
// @Failure      400      {object}  handler.ErrorResponse    "Validation error."
// @Failure      401      {object}  handler.ErrorResponse    "Unauthorized access."
// @Failure      403      {object}  handler.ErrorResponse    "Forbidden. Token belongs to another user."
// @Failure      500      {object}  handler.ErrorResponse    "Internal server error."
// @Router       /user/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	userIDVal, exists := c.Get(middleware.CtxUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, handler.ErrorResponse{Error: "Unauthorized access."})
		return
	}

	userIDStr, ok := userIDVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, handler.ErrorResponse{Error: "Internal server error."})
		return
	}

	var req handler.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, handler.ErrorResponse{Error: err.Error()})
		return
	}

	err := h.service.Logout(c.Request.Context(), userIDStr, req.RefreshToken)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, handler.MessageResponse{Message: "Successfully logged out."})
}

// RevokeAllSessions godoc
// @Summary      Revoke all sessions
// @Description  Invalidates all active sessions for the user (deletes all refresh tokens).
// @Tags         User
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200      {object}  handler.MessageResponse  "Success message."
// @Failure      401      {object}  handler.ErrorResponse    "Unauthorized access."
// @Failure      500      {object}  handler.ErrorResponse    "Internal server error."
// @Router       /user/sessions/revoke-all [post]
func (h *Handler) RevokeAllSessions(c *gin.Context) {
	ctx := c.Request.Context()
	userIDVal, exists := c.Get(middleware.CtxUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, handler.ErrorResponse{Error: "Unauthorized access."})
		return
	}
	userID, ok := userIDVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, handler.ErrorResponse{Error: "Internal server error."})
		return
	}
	err := h.service.RevokeAllSessions(ctx, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	logger.InfoKV(ctx, "all user sessions revoked", "user_id", userID)
	c.JSON(http.StatusOK, handler.MessageResponse{Message: "All sessions have been successfully revoked."})
}

func (h *Handler) GetUserStatus(c *gin.Context) {
	requesterUserID, requesterRole, ok := getRequester(c)
	if !ok {
		return
	}

	targetUserID := c.Param("userId")
	status, err := h.service.GetUserStatus(c.Request.Context(), requesterUserID, requesterRole, targetUserID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, handler.UserStatusResponse{Status: string(status)})
}

func (h *Handler) UpdateUserStatus(c *gin.Context) {
	requesterUserID, requesterRole, ok := getRequester(c)
	if !ok {
		return
	}

	targetUserID := c.Param("userId")

	var req handler.UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, handler.ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.service.UpdateUserStatus(c.Request.Context(), requesterUserID, requesterRole, targetUserID, models.UserStatus(req.Status)); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, handler.MessageResponse{Message: "user status has been updated"})
}

func getRequester(c *gin.Context) (string, string, bool) {
	requesterUserIDVal, exists := c.Get(middleware.CtxUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, handler.ErrorResponse{Error: "unauthorized"})
		return "", "", false
	}

	requesterRoleVal, exists := c.Get(middleware.CtxUserRole)
	if !exists {
		c.JSON(http.StatusUnauthorized, handler.ErrorResponse{Error: "unauthorized"})
		return "", "", false
	}

	requesterUserID, userIDOk := requesterUserIDVal.(string)
	requesterRole, roleOk := requesterRoleVal.(string)
	if !userIDOk || !roleOk {
		c.JSON(http.StatusInternalServerError, handler.ErrorResponse{Error: "internal server error"})
		return "", "", false
	}

	return requesterUserID, requesterRole, true
}
