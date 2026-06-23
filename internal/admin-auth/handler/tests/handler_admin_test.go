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
	"github.com/stretchr/testify/require"

	"github.com/ua-academy-projects/share-bite/internal/admin-auth/dto"
	apperr "github.com/ua-academy-projects/share-bite/internal/admin-auth/error"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/handler/admin"
)

func TestAdminHandler_GetUserDetails(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success returns 200 and user details including customer profile", func(t *testing.T) {
		mockSvc := new(MockAdminService)
		h := admin.NewHandler(mockSvc, nil)

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
		h := admin.NewHandler(mockSvc, nil)
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
		h := admin.NewHandler(mockSvc, nil)

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
		h := admin.NewHandler(mockSvc, nil)

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

func TestAdminHandler_GetPendingBusinesses(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success returns 200 with custom pagination", func(t *testing.T) {
		mockSvc := new(MockAdminService)
		h := admin.NewHandler(mockSvc, nil)

		expectedResp := &dto.PaginatedPendingBusinessesResponse{}
		mockSvc.On("GetPendingBusinessesList", mock.Anything, 15, 5).Return(expectedResp, nil)

		r := gin.New()
		r.GET("/admin/businesses/pending", h.GetPendingBusinesses)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/admin/businesses/pending?limit=15&offset=5", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Success sanitizes out of bounds parameters", func(t *testing.T) {
		mockSvc := new(MockAdminService)
		h := admin.NewHandler(mockSvc, nil)

		expectedResp := &dto.PaginatedPendingBusinessesResponse{}
		mockSvc.On("GetPendingBusinessesList", mock.Anything, 50, 0).Return(expectedResp, nil)

		r := gin.New()
		r.GET("/admin/businesses/pending", h.GetPendingBusinesses)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/admin/businesses/pending?limit=150&offset=0", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestAdminHandler_GetPlatformStatistics(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success returns 200 and statistics payload", func(t *testing.T) {
		mockSvc := new(MockAdminService)
		h := admin.NewHandler(mockSvc, nil)

		expected := &dto.PlatformStatisticsResponse{
			TotalUsers:                   100,
			TotalAdminUsers:              2,
			TotalModeratorUsers:          3,
			TotalRegularUsers:            85,
			TotalBusinessRoleUsers:       10,
			TotalActiveUsers:             90,
			TotalMutedUsers:              5,
			TotalSuspendedUsers:          5,
			TotalCustomers:               80,
			TotalGuestPosts:              200,
			TotalGuestComments:           840,
			TotalGuestPostLikes:          1500,
			TotalCollections:             45,
			AvgPostsPerCustomer:          2.5,
			AvgCommentsPerCustomer:       4.2,
			AvgCommentsPerPost:           1.3,
			CollectionsWithCollaborators: 12,
			PostsWithCollaborators:       7,
			TotalBusinessOrgUnits:        25,
			TotalBusinessPosts:           15,
			TotalBusinessComments:        75,
			TotalBusinessLikes:           120,
			TotalBusinessBoxes:           8,
			TotalBusinessBoxItems:        32,
			AvgPostsPerBusiness:          3.1,
			AvgCommentsPerBusiness:       5.0,
			AvgBusinessCommentsPerPost:   0.8,
		}

		mockSvc.On("GetPlatformStatistics", mock.Anything).Return(expected, nil)

		r := gin.New()
		r.GET("/admin/statistics", h.GetPlatformStatistics)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/admin/statistics", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		expectedBody, err := json.Marshal(expected)
		require.NoError(t, err)
		require.JSONEq(t, string(expectedBody), w.Body.String())

		mockSvc.AssertExpectations(t)
	})

	t.Run("Service error returns 500 via middleware", func(t *testing.T) {
		mockSvc := new(MockAdminService)
		h := admin.NewHandler(mockSvc, nil)

		appError := &apperr.AppError{Code: http.StatusInternalServerError, Message: "database connection failed"}
		mockSvc.On("GetPendingBusinessesList", mock.Anything, 10, 0).Return(nil, appError)

		r := gin.New()
		r.Use(testErrorMiddleware())
		r.GET("/admin/businesses/pending", h.GetPendingBusinesses)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/admin/businesses/pending", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "database connection failed")
		mockSvc.AssertExpectations(t)
	})
}

func TestAdminHandler_ReviewBusiness(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type ReviewRequest struct {
		Status  string `json:"status"`
		Comment string `json:"comment"`
	}

	mockAdminUUID := "admin-uuid-1111-2222"
	mockAuthMw := func(c *gin.Context) {
		c.Set("userId", mockAdminUUID)
		c.Next()
	}

	t.Run("Success verification returns 200", func(t *testing.T) {
		mockSvc := new(MockAdminService)
		h := admin.NewHandler(mockSvc, nil)

		expectedParams := dto.ReviewBusinessParams{
			OrgUnitID: 42,
			NewStatus: "verified",
			AdminID:   mockAdminUUID,
			Comment:   new("All documents are fine"),
		}

		mockSvc.On("ReviewBusinessStatus", mock.Anything, expectedParams).Return(nil)

		r := gin.New()
		r.Use(mockAuthMw)
		r.PATCH("/admin/businesses/:id/review", h.ReviewBusiness)

		body := ReviewRequest{Status: "verified", Comment: "All documents are fine"}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/admin/businesses/42/review", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Business verification status updated successfully.")
		mockSvc.AssertExpectations(t)
	})

	t.Run("Invalid integer ID format returns 400", func(t *testing.T) {
		mockSvc := new(MockAdminService)
		h := admin.NewHandler(mockSvc, nil)

		r := gin.New()
		r.PATCH("/admin/businesses/:id/review", h.ReviewBusiness)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/admin/businesses/abc/review", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "must be a positive integer")
	})

	t.Run("Invalid body payload returns 400", func(t *testing.T) {
		mockSvc := new(MockAdminService)
		h := admin.NewHandler(mockSvc, nil)

		r := gin.New()
		r.PATCH("/admin/businesses/:id/review", h.ReviewBusiness)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/admin/businesses/10/review", bytes.NewBufferString("{invalid-json"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request payload")
	})

	t.Run("Unauthorized when admin context is unresolved", func(t *testing.T) {
		mockSvc := new(MockAdminService)
		h := admin.NewHandler(mockSvc, nil)

		r := gin.New()
		r.PATCH("/admin/businesses/:id/review", h.ReviewBusiness)

		body := ReviewRequest{Status: "rejected", Comment: "Bad license"}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/admin/businesses/10/review", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "admin identity could not be resolved")
	})
}
