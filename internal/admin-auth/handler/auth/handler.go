package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	apperr "github.com/ua-academy-projects/share-bite/internal/admin-auth/error"
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
// @Summary      Авторизація користувача
// @Description  Перевіряє email та пароль.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      LoginRequest  true  "Дані для входу"
// @Success      200      {object}  object  "Успіх. Повертає JSON: {'access_token': '...', 'refresh_token': '...'}"
// @Failure      400      {object}  object  "Помилка валідації: {'message': '...'}"
// @Failure      401      {object}  object  "Невірні облікові дані: {'error': '...'}"
// @Failure      500      {object}  object  "Внутрішня помилка сервера: {'error': '...'}"
// @Router       /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
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

// Register godoc
// @Summary      Реєстрація користувача
// @Description  Створює нового користувача та одразу повертає пару токенів.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      RegisterRequest  true  "Дані для реєстрації"
// @Success      201      {object}  object  "Успіх. Повертає JSON: {'access_token': '...', 'refresh_token': '...'}"
// @Failure      400      {object}  object  "Помилка валідації: {'message': '...'}"
// @Failure      409      {object}  object  "Користувач вже існує: {'error': '...'}"
// @Failure      422      {object}  object  "Роль не знайдена: {'error': '...'}"
// @Failure      500      {object}  object  "Внутрішня помилка сервера: {'error': '...'}"
// @Router       /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
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
// @Summary      Оновлення токенів
// @Description  Генерує нову пару токенів на основі валідного refresh_token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      RefreshRequest  true  "Refresh токен"
// @Success      200      {object}  object  "Успіх. Повертає JSON: {'access_token': '...', 'refresh_token': '...'}"
// @Failure      400      {object}  object  "Помилка валідації: {'message': '...'}"
// @Failure      401      {object}  object  "Невалідний або прострочений токен: {'error': '...'}"
// @Failure      500      {object}  object  "Внутрішня помилка сервера: {'error': '...'}"
// @Router       /auth/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
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
// @Summary      OAuth Авторизація / Реєстрація
// @Description  Обмінює код від провайдера на токени доступу.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        provider path      string                true  "Назва провайдера (google, github)"
// @Param        request  body      OAuthCallbackRequest  true  "Код від провайдера та роль"
// @Success      200      {object}  object  "Успіх. Повертає JSON: {'access_token': '...', 'refresh_token': '...'}"
// @Failure      400      {object}  object  "Непідтримуваний провайдер: {'error': '...'}"
// @Failure      502      {object}  object  "Помилка обміну коду з провайдером: {'error': '...'}"
// @Failure      500      {object}  object  "Внутрішня помилка сервера: {'error': '...'}"
// @Router       /auth/oauth/{provider}/callback [post]
func (h *Handler) OAuthCallback(c *gin.Context) {
	ctx := c.Request.Context()
	providerName := c.Param("provider")

	p, err := h.providerFactory.Get(providerName)
	if err != nil {
		logger.WarnKV(ctx, "unsupported provider requested", "provider", providerName, "error", err.Error())
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported or invalid authentication provider."})
		return
	}

	var req OAuthCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WarnKV(ctx, "invalid json payload", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload."})
		return
	}

	tokens, err := h.service.OAuthLogin(ctx, p, req.Code, req.Slug)
	if err != nil {
		_ = c.Error(err)
		var appErr *apperr.AppError
		if errors.As(err, &appErr) {
			logger.WarnKV(ctx, "oauth login failed (client/domain error)", "provider", providerName, "slug", req.Slug, "error", appErr.Error())
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
			return
		}
		logger.ErrorKV(ctx, "oauth login failed (internal error)", "provider", providerName, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error. Please try again later.",
		})
		return
	}

	logger.InfoKV(ctx, "user oauth login success", "provider", providerName)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})
}

// OAuthLinkAccount godoc
// @Summary      Прив'язка соцмережі до акаунту
// @Description  Прив'язує Google або GitHub до вже авторизованого користувача.
// @Tags         User
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        provider path      string            true  "Назва провайдера (google, github)"
// @Param        request  body      OAuthLinkRequest  true  "Код від провайдера"
// @Success      200      {object}  object  "Успіх. Повертає: {'message': 'Social account successfully linked.'}"
// @Failure      400      {object}  object  "Невалідний запит: {'error': '...'}"
// @Failure      401      {object}  object  "Неавторизований доступ: {'error': '...'}"
// @Failure      409      {object}  object  "Провайдер вже прив'язаний: {'error': '...'}"
// @Failure      500      {object}  object  "Внутрішня помилка сервера: {'error': '...'}"
// @Router       /user/link/{provider} [post]
func (h *Handler) OAuthLinkAccount(c *gin.Context) {
	ctx := c.Request.Context()

	userIDVal, exists := c.Get(middleware.CtxUserID)
	if !exists {
		logger.WarnKV(ctx, "missing user_id in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access."})
		return
	}
	userID, ok := userIDVal.(string)
	if !ok {
		logger.ErrorKV(ctx, "invalid user id type in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
		return
	}

	providerName := c.Param("provider")
	p, err := h.providerFactory.Get(providerName)
	if err != nil {
		logger.WarnKV(ctx, "unsupported provider for linking", "provider", providerName, "error", err.Error())
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported authentication provider."})
		return
	}

	var req OAuthLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WarnKV(ctx, "invalid json payload during account link", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload. 'code' is required."})
		return
	}

	err = h.service.LinkProvider(ctx, userID, p, req.Code)
	if err != nil {
		_ = c.Error(err)

		var appErr *apperr.AppError
		if errors.As(err, &appErr) {
			logger.WarnKV(ctx, "failed to link account (domain error)", "user_id", userID, "provider", providerName, "error", appErr.Error())
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
			return
		}
		logger.ErrorKV(ctx, "failed to link account (internal error)", "user_id", userID, "provider", providerName, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error.",
		})
		return
	}

	logger.InfoKV(ctx, "social account successfully linked", "user_id", userID, "provider", providerName)

	c.JSON(http.StatusOK, gin.H{
		"message": "Social account successfully linked.",
	})
}
