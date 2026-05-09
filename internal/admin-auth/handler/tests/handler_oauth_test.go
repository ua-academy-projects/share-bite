package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	apperr "github.com/ua-academy-projects/share-bite/internal/admin-auth/error"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/handler/auth"         // пакет з твоїми хендлерами та Request-структурами
	authsvc "github.com/ua-academy-projects/share-bite/internal/admin-auth/service/auth" // пакет з інтерфейсами провайдера та сервісу
	"github.com/ua-academy-projects/share-bite/internal/middleware"
)

func TestHandler_OAuthCallback(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Unsupported Provider", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		mockFactory := new(MockProviderFactory)

		h := auth.NewHandler(mockSvc, mockFactory)

		expectedErr := &apperr.AppError{
			Code:    http.StatusBadRequest,
			Message: "Unsupported or invalid authentication provider",
		}
		mockFactory.On("Get", "unknown").Return(nil, expectedErr)

		r := gin.New()
		r.Use(testErrorMiddleware())
		r.POST("/oauth/:provider/callback", h.OAuthCallback)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/oauth/unknown/callback", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Unsupported or invalid authentication provider")
	})

	t.Run("Successful Login", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		mockFactory := new(MockProviderFactory)
		mockProvider := new(MockOAuthProvider)
		h := auth.NewHandler(mockSvc, mockFactory)

		mockFactory.On("Get", "google").Return(mockProvider, nil)

		expectedTokens := &authsvc.Tokens{
			AccessToken:  "mock-access-token",
			RefreshToken: "mock-refresh-token",
		}
		mockSvc.On("OAuthLogin", mock.Anything, mockProvider, "valid-code", "user").Return(expectedTokens, nil)

		r := gin.New()
		r.POST("/oauth/:provider/callback", h.OAuthCallback)

		reqBody := auth.OAuthCallbackRequest{Code: "valid-code", Slug: "user"}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/oauth/google/callback", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"access_token": "mock-access-token", "refresh_token": "mock-refresh-token"}`, w.Body.String())
	})

	t.Run("Service Domain Error", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		mockFactory := new(MockProviderFactory)
		mockProvider := new(MockOAuthProvider)
		h := auth.NewHandler(mockSvc, mockFactory)

		mockFactory.On("Get", "google").Return(mockProvider, nil)

		appError := &apperr.AppError{Code: http.StatusForbidden, Message: "Role not allowed"}
		mockSvc.On("OAuthLogin", mock.Anything, mockProvider, "valid-code", "user").Return(nil, appError)

		r := gin.New()
		r.Use(testErrorMiddleware())
		r.POST("/oauth/:provider/callback", h.OAuthCallback)

		reqBody := auth.OAuthCallbackRequest{Code: "valid-code", Slug: "user"}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/oauth/google/callback", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.JSONEq(t, `{"error": "Role not allowed"}`, w.Body.String())
	})
}

func TestHandler_OAuthLinkAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Missing UserID in Context", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		mockFactory := new(MockProviderFactory)
		h := auth.NewHandler(mockSvc, mockFactory)

		r := gin.New()
		r.POST("/user/link/:provider", h.OAuthLinkAccount)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/user/link/google", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Successful Account Link", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		mockFactory := new(MockProviderFactory)
		mockProvider := new(MockOAuthProvider)
		h := auth.NewHandler(mockSvc, mockFactory)

		mockFactory.On("Get", "google").Return(mockProvider, nil)
		mockSvc.On("LinkProvider", mock.Anything, "user-123", mockProvider, "valid-code").Return(nil)

		r := gin.New()

		r.Use(func(c *gin.Context) {
			c.Set(middleware.CtxUserID, "user-123")
			c.Next()
		})
		r.POST("/user/link/:provider", h.OAuthLinkAccount)

		reqBody := auth.OAuthLinkRequest{Code: "valid-code"}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/user/link/google", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message": "Social account successfully linked."}`, w.Body.String())
	})
}
