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

	"github.com/ua-academy-projects/share-bite/internal/admin-auth/dto"
	apperr "github.com/ua-academy-projects/share-bite/internal/admin-auth/error"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/handler/admin"
)

func TestAdminHandler_GetUserDetails(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success returns 200 and user details including customer profile", func(t *testing.T) {
		mockSvc := new(MockAdminService)
		h := admin.NewHandler(mockSvc)

		expectedUser := &dto.FullUserDetails{
			ID:       "68d15417-f340-445f-a3a1-4d934fea5bbc",
			Email:    "test1234@gmail.com",
			RoleSlug: "user",
			CustomerProfile: &dto.CustomerProfileData{
				FirstName: "petro",
				LastName:  "petrov",
				Bio:       "ejvnjinbv",
			},
		}

		mockSvc.On("GetUserDetails", mock.Anything, "68d15417-f340-445f-a3a1-4d934fea5bbc").Return(expectedUser, nil)

		r := gin.New()
		r.GET("/admin/users/:id", h.GetUserDetails)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/admin/users/68d15417-f340-445f-a3a1-4d934fea5bbc", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "petro")
		mockSvc.AssertExpectations(t)
	})

	t.Run("User not found returns 404 via middleware", func(t *testing.T) {
		mockSvc := new(MockAdminService)
		h := admin.NewHandler(mockSvc)
		fakeUUID := "00000000-0000-0000-0000-000000000000"

		mockSvc.On("GetUserDetails", mock.Anything, fakeUUID).Return(nil, apperr.ErrUserNotFound)

		r := gin.New()
		r.Use(testErrorMiddleware())
		r.GET("/admin/users/:id", h.GetUserDetails)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/admin/users/"+fakeUUID, nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, apperr.ErrUserNotFound.HTTPStatus(), w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestAdminHandler_ChangeUserRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	type ChangeRoleRequest struct {
		RoleSlug string `json:"role_slug"`
	}

	targetUUID := "11111111-1111-1111-1111-111111111111"

	t.Run("Success returns 200", func(t *testing.T) {
		mockSvc := new(MockAdminService)
		h := admin.NewHandler(mockSvc)

		mockSvc.On("ChangeUserRole", mock.Anything, targetUUID, "moderator").Return(nil)

		r := gin.New()
		r.PATCH("/admin/users/:id/role", h.ChangeUserRole)

		reqBody := ChangeRoleRequest{RoleSlug: "moderator"}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/admin/users/"+targetUUID+"/role", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Invalid transition returns 400 via middleware", func(t *testing.T) {
		mockSvc := new(MockAdminService)
		h := admin.NewHandler(mockSvc)

		appError := &apperr.AppError{Code: http.StatusBadRequest, Message: "invalid role transition from 'user' to 'admin'"}
		mockSvc.On("ChangeUserRole", mock.Anything, targetUUID, "admin").Return(appError)

		r := gin.New()
		r.Use(testErrorMiddleware())
		r.PATCH("/admin/users/:id/role", h.ChangeUserRole)

		reqBody := ChangeRoleRequest{RoleSlug: "admin"}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/admin/users/"+targetUUID+"/role", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid role transition")
		mockSvc.AssertExpectations(t)
	})
}
