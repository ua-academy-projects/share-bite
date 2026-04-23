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
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/handler/auth"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
)

func TestHandler_Logout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Bad JSON body returns 400", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := auth.NewHandler(mockSvc, nil)

		r := gin.New()
		r.POST("/auth/logout", h.Logout)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/logout", bytes.NewBufferString("{invalid}"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "Logout", mock.Anything, mock.Anything)
	})

	t.Run("Success returns 200", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := auth.NewHandler(mockSvc, nil)

		mockSvc.On("Logout", mock.Anything, "valid-refresh-token").Return(nil)

		r := gin.New()
		r.POST("/auth/logout", h.Logout)

		body, _ := json.Marshal(auth.RefreshRequest{RefreshToken: "valid-refresh-token"})
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/logout", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message": "Successfully logged out."}`, w.Body.String())
		mockSvc.AssertExpectations(t)
	})

	t.Run("Invalid token returns 401 via error middleware", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := auth.NewHandler(mockSvc, nil)

		mockSvc.On("Logout", mock.Anything, "expired-token").
			Return(apperr.ErrInvalidToken)

		r := gin.New()
		r.Use(testErrorMiddleware())
		r.POST("/auth/logout", h.Logout)

		body, _ := json.Marshal(auth.RefreshRequest{RefreshToken: "expired-token"})
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/logout", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, apperr.ErrInvalidToken.HTTPStatus(), w.Code)
		assert.Contains(t, w.Body.String(), apperr.ErrInvalidToken.Message)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Internal service error returns 500 via error middleware", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := auth.NewHandler(mockSvc, nil)

		serviceErr := &apperr.AppError{Code: http.StatusInternalServerError, Message: "failed to logout"}
		mockSvc.On("Logout", mock.Anything, "some-token").Return(serviceErr)

		r := gin.New()
		r.Use(testErrorMiddleware())
		r.POST("/auth/logout", h.Logout)

		body, _ := json.Marshal(auth.RefreshRequest{RefreshToken: "some-token"})
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/logout", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestHandler_RevokeAllSessions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Missing UserID in context returns 401", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := auth.NewHandler(mockSvc, nil)

		r := gin.New()
		r.POST("/user/sessions/revoke-all", h.RevokeAllSessions)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/user/sessions/revoke-all", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, `{"error": "Unauthorized access."}`, w.Body.String())
		mockSvc.AssertNotCalled(t, "RevokeAllSessions", mock.Anything, mock.Anything)
	})

	t.Run("UserID wrong type in context returns 500", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := auth.NewHandler(mockSvc, nil)

		r := gin.New()
		r.Use(func(c *gin.Context) {
			// навмисно кладемо не string
			c.Set(middleware.CtxUserID, 12345)
			c.Next()
		})
		r.POST("/user/sessions/revoke-all", h.RevokeAllSessions)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/user/sessions/revoke-all", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error": "Internal server error."}`, w.Body.String())
		mockSvc.AssertNotCalled(t, "RevokeAllSessions", mock.Anything, mock.Anything)
	})

	t.Run("Success returns 200", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := auth.NewHandler(mockSvc, nil)

		mockSvc.On("RevokeAllSessions", mock.Anything, "user-42").Return(nil)

		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set(middleware.CtxUserID, "user-42")
			c.Next()
		})
		r.POST("/user/sessions/revoke-all", h.RevokeAllSessions)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/user/sessions/revoke-all", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message": "All sessions have been successfully revoked."}`, w.Body.String())
		mockSvc.AssertExpectations(t)
	})

	t.Run("Service error is forwarded to gin context", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := auth.NewHandler(mockSvc, nil)

		serviceErr := &apperr.AppError{Code: http.StatusInternalServerError, Message: "failed to revoke all sessions"}
		mockSvc.On("RevokeAllSessions", mock.Anything, "user-42").Return(serviceErr)

		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set(middleware.CtxUserID, "user-42")
			c.Next()
		})
		r.Use(testErrorMiddleware())
		r.POST("/user/sessions/revoke-all", h.RevokeAllSessions)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/user/sessions/revoke-all", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}
