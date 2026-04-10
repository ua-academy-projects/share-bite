package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ua-academy-projects/share-bite/internal/admin-auth/dto"
	apperr "github.com/ua-academy-projects/share-bite/internal/admin-auth/error"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/handler/auth"         // пакет з твоїми хендлерами та Request-структурами
	authsvc "github.com/ua-academy-projects/share-bite/internal/admin-auth/service/auth" // пакет з інтерфейсами провайдера та сервісу
	"github.com/ua-academy-projects/share-bite/internal/middleware"
)

type MockProviderFactory struct {
	mock.Mock
}

func (m *MockProviderFactory) Get(name string) (authsvc.OAuthProvider, error) {
	args := m.Called(name)
	if args.Get(0) != nil {
		return args.Get(0).(authsvc.OAuthProvider), args.Error(1)
	}
	return nil, args.Error(1)
}

type MockOAuthProvider struct {
	mock.Mock
}

func (m *MockOAuthProvider) ExchangeCode(ctx context.Context, code string) (*dto.OAuthUserInfo, error) {
	args := m.Called(ctx, code)
	if args.Get(0) != nil {
		return args.Get(0).(*dto.OAuthUserInfo), args.Error(1)
	}
	return nil, args.Error(1)
}

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(ctx context.Context, email string, password string) (*authsvc.Tokens, error) {
	args := m.Called(ctx, email, password)
	if args.Get(0) != nil {
		return args.Get(0).(*authsvc.Tokens), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthService) Register(ctx context.Context, email string, password string, slug string) (*authsvc.Tokens, error) {
	args := m.Called(ctx, email, password, slug)
	if args.Get(0) != nil {
		return args.Get(0).(*authsvc.Tokens), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthService) Refresh(ctx context.Context, refreshToken string) (*authsvc.Tokens, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) != nil {
		return args.Get(0).(*authsvc.Tokens), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthService) OAuthLogin(ctx context.Context, provider authsvc.OAuthProvider, code string, slug string) (*authsvc.Tokens, error) {
	args := m.Called(ctx, provider, code, slug)
	if args.Get(0) != nil {
		return args.Get(0).(*authsvc.Tokens), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthService) LinkProvider(ctx context.Context, userID string, provider authsvc.OAuthProvider, code string) error {
	args := m.Called(ctx, userID, provider, code)
	return args.Error(0)
}

func (m *MockAuthService) RecoverAccess(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockAuthService) ResetPassword(ctx context.Context, token string, newPassword string) error {
	args := m.Called(ctx, token, newPassword)
	return args.Error(0)
}

func testErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.Last()
		if err == nil {
			return
		}

		var appErr *apperr.AppError
		if errors.As(err.Err, &appErr) {
			c.JSON(appErr.HTTPStatus(), gin.H{"error": appErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

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
