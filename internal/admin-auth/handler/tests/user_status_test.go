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
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/handler"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/handler/auth"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/models"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
)

func TestHandler_GetUserStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("self access returns 200", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := auth.NewHandler(mockSvc, nil, nil)

		mockSvc.On("GetUserStatus", mock.Anything, "user-1", "user", "user-1").
			Return(models.UserStatusActive, nil)

		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set(middleware.CtxUserID, "user-1")
			c.Set(middleware.CtxUserRole, "user")
			c.Next()
		})
		r.GET("/users/:userId/status", h.GetUserStatus)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/users/user-1/status", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"status":"active"}`, w.Body.String())
		mockSvc.AssertExpectations(t)
	})

	t.Run("service forbidden error returns 403", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := auth.NewHandler(mockSvc, nil, nil)

		mockSvc.On("GetUserStatus", mock.Anything, "user-2", "user", "user-1").
			Return(models.UserStatus(""), apperr.ErrForbiddenStatusRead)

		r := gin.New()
		r.Use(testErrorMiddleware())
		r.Use(func(c *gin.Context) {
			c.Set(middleware.CtxUserID, "user-2")
			c.Set(middleware.CtxUserRole, "user")
			c.Next()
		})
		r.GET("/users/:userId/status", h.GetUserStatus)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/users/user-1/status", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestHandler_UpdateUserStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("invalid body returns 400", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := auth.NewHandler(mockSvc, nil, nil)

		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set(middleware.CtxUserID, "admin-1")
			c.Set(middleware.CtxUserRole, "admin")
			c.Next()
		})
		r.PUT("/users/:userId/status", h.UpdateUserStatus)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/users/user-1/status", bytes.NewBufferString(`{"status":"wrong"}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertNotCalled(t, "UpdateUserStatus", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("success returns 200", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := auth.NewHandler(mockSvc, nil, nil)

		mockSvc.On("UpdateUserStatus", mock.Anything, "admin-1", "admin", "user-1", models.UserStatusMuted).
			Return(nil)

		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set(middleware.CtxUserID, "admin-1")
			c.Set(middleware.CtxUserRole, "admin")
			c.Next()
		})
		r.PUT("/users/:userId/status", h.UpdateUserStatus)

		body, err := json.Marshal(handler.UpdateUserStatusRequest{Status: "muted"})
		assert.NoError(t, err)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/users/user-1/status", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"user status has been updated"}`, w.Body.String())
		mockSvc.AssertExpectations(t)
	})

	t.Run("service not found error returns 404", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h := auth.NewHandler(mockSvc, nil, nil)

		mockSvc.On("UpdateUserStatus", mock.Anything, "admin-1", "admin", "user-404", models.UserStatusSuspended).
			Return(apperr.ErrUserNotFound)

		r := gin.New()
		r.Use(testErrorMiddleware())
		r.Use(func(c *gin.Context) {
			c.Set(middleware.CtxUserID, "admin-1")
			c.Set(middleware.CtxUserRole, "admin")
			c.Next()
		})
		r.PUT("/users/:userId/status", h.UpdateUserStatus)

		body, err := json.Marshal(handler.UpdateUserStatusRequest{Status: "suspended"})
		assert.NoError(t, err)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/users/user-404/status", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockSvc.AssertExpectations(t)
	})
}
