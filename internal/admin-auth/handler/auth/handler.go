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
//
//	@Summary		Авторизація користувача
//	@Description	Перевіряє email та пароль.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		LoginRequest	true	"Дані для входу"
//	@Success		200		{object}	object			"Успіх. Повертає JSON: {'access_token': '...', 'refresh_token': '...'}"
//	@Failure		400		{object}	object			"Помилка валідації: {'error': '...'}"
//	@Failure		401		{object}	object			"Невірні облікові дані: {'error': '...'}"
//	@Failure		500		{object}	object			"Внутрішня помилка сервера: {'error': '...'}"
//	@Router			/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})
}

// RecoverAccess godoc
//
//	@Summary		Recover access
//	@Description	Sends password reset instructions if account exists
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		RecoverAccessRequest	true	"Recover access payload"
//	@Success		200		{object}	MessageResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/auth/recover-access [post]
func (h *Handler) RecoverAccess(c *gin.Context) {
	var req RecoverAccessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.service.RecoverAccess(c.Request.Context(), req.Email)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "If the email exists, recovery instructions have been sent"})
}

// ResetPassword godoc
//
//	@Summary		Reset password
//	@Description	Resets password by token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		ResetPasswordRequest	true	"Reset password payload"
//	@Success		200		{object}	MessageResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/auth/reset-password [post]
func (h *Handler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	err := h.service.ResetPassword(c.Request.Context(), req.Token, req.NewPassword)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "password has been reset"})
}

// Register godoc
//
//	@Summary		Реєстрація користувача
//	@Description	Створює нового користувача та одразу повертає пару токенів.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		RegisterRequest	true	"Дані для реєстрації"
//	@Success		201		{object}	object			"Успіх. Повертає JSON: {'access_token': '...', 'refresh_token': '...'}"
//	@Failure		400		{object}	object			"Помилка валідації: {'error': '...'}"
//	@Failure		409		{object}	object			"Користувач вже існує: {'error': '...'}"
//	@Failure		422		{object}	object			"Роль не знайдена: {'error': '...'}"
//	@Failure		500		{object}	object			"Внутрішня помилка сервера: {'error': '...'}"
//	@Router			/auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.service.Register(
		c.Request.Context(),
		req.Email,
		req.Password,
		req.Slug,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})
}

// Refresh godoc
//
//	@Summary		Оновлення токенів
//	@Description	Генерує нову пару токенів на основі валідного refresh_token.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		RefreshRequest	true	"Refresh токен"
//	@Success		200		{object}	object			"Успіх. Повертає JSON: {'access_token': '...', 'refresh_token': '...'}"
//	@Failure		400		{object}	object			"Помилка валідації: {'error': '...'}"
//	@Failure		401		{object}	object			"Невалідний або прострочений токен: {'error': '...'}"
//	@Failure		500		{object}	object			"Внутрішня помилка сервера: {'error': '...'}"
//	@Router			/auth/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.service.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})
}

// OAuthCallback godoc
//
//	@Summary		OAuth Авторизація / Реєстрація
//	@Description	Обмінює код від провайдера на токени доступу.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			provider	path		string					true	"Назва провайдера (google)"
//	@Param			request		body		OAuthCallbackRequest	true	"Код від провайдера та роль"
//	@Success		200			{object}	object					"Успіх. Повертає JSON: {'access_token': '...', 'refresh_token': '...'}"
//	@Failure		400			{object}	object					"Непідтримуваний провайдер: {'error': '...'}"
//	@Failure		422			{object}	object					"Роль не знайдена: {'error': '...'}"
//	@Failure		502			{object}	object					"Помилка обміну коду з провайдером: {'error': '...'}"
//	@Failure		500			{object}	object					"Внутрішня помилка сервера: {'error': '...'}"
//	@Router			/auth/oauth/{provider}/callback [post]
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload."})
		return
	}

	tokens, err := h.service.OAuthLogin(ctx, p, req.Code, req.Slug)
	if err != nil {
		_ = c.Error(err)
		return
	}

	logger.InfoKV(ctx, "user oauth login success", "provider", providerName)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})
}

// OAuthLinkAccount godoc
//
//	@Summary		Прив'язка соцмережі до акаунту
//	@Description	Прив'язує Google або GitHub до вже авторизованого користувача.
//	@Tags			User
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			provider	path		string				true	"Назва провайдера (google)"
//	@Param			request		body		OAuthLinkRequest	true	"Код від провайдера"
//	@Success		200			{object}	object				"Успіх. Повертає: {'message': 'Social account successfully linked.'}"
//	@Failure		400			{object}	object				"Невалідний запит: {'error': '...'}"
//	@Failure		401			{object}	object				"Неавторизований доступ: {'error': '...'}"
//	@Failure		409			{object}	object				"Провайдер вже прив'язаний: {'error': '...'}"
//	@Failure		500			{object}	object				"Внутрішня помилка сервера: {'error': '...'}"
//	@Failure		502			{object}	object				"Помилка обміну коду з провайдером: {'error': '...'}"
//	@Router			/user/link/{provider} [post]
func (h *Handler) OAuthLinkAccount(c *gin.Context) {
	ctx := c.Request.Context()

	userIDVal, exists := c.Get(middleware.CtxUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access."})
		return
	}
	userID, ok := userIDVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload. 'code' is required."})
		return
	}

	err = h.service.LinkProvider(ctx, userID, p, req.Code)
	if err != nil {
		_ = c.Error(err)
		return
	}

	logger.InfoKV(ctx, "social account successfully linked", "user_id", userID, "provider", providerName)

	c.JSON(http.StatusOK, gin.H{
		"message": "Social account successfully linked.",
	})
}
