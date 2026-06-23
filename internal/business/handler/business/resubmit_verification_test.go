package business

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/pkg/jwt"

	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/internal/business/error/code"
)

const validTestUUID = "123e4567-e89b-12d3-a456-426614174000"
const validHackerUUID = "987e4567-e89b-12d3-a456-426614174000"

type testAccessTokenParser struct {
	mockUserID string
	mockRole   string
	mockStatus jwt.UserStatus
}

func (p testAccessTokenParser) ParseAccessToken(string) (jwt.AccessTokenPayload, error) {
	return jwt.AccessTokenPayload{
		UserID: p.mockUserID,
		Role:   p.mockRole,
		Status: p.mockStatus,
	}, nil
}

type resubmitVerificationServiceMock struct {
	businessService
	called    bool
	gotOrgID  int
	gotUserID string
	err       error
}

func (m *resubmitVerificationServiceMock) ResubmitVerification(_ context.Context, orgUnitID int, userID string) error {
	m.called = true
	m.gotOrgID = orgUnitID
	m.gotUserID = userID
	return m.err
}

func setupResubmitRouter(s businessService, mockUserID, mockRole string) *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(testBusinessErrorMiddleware())

	parser := testAccessTokenParser{
		mockUserID: mockUserID,
		mockRole:   mockRole,
		mockStatus: jwt.UserStatusActive,
	}
	RegisterHandlers(r.Group("/business"), s, parser, nil, nil)

	return r
}
func TestResubmitVerification_Success(t *testing.T) {
	mock := &resubmitVerificationServiceMock{}

	router := setupResubmitRouter(mock, validTestUUID, "business")

	req, _ := http.NewRequest(http.MethodPost, "/business/42/resubmit", nil)
	req.Header.Set("Authorization", "Bearer dynamic-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}

	if !mock.called {
		t.Fatalf("service ResubmitVerification was not called")
	}
	if mock.gotOrgID != 42 {
		t.Fatalf("expected orgUnitID=42, got %d", mock.gotOrgID)
	}
	if mock.gotUserID != validTestUUID {
		t.Fatalf("expected userID=%s, got %s", validTestUUID, mock.gotUserID)
	}
}

func TestResubmitVerification_InvalidURI_ReturnsBadRequest(t *testing.T) {
	mock := &resubmitVerificationServiceMock{}
	router := setupResubmitRouter(mock, validTestUUID, "business")

	req, _ := http.NewRequest(http.MethodPost, "/business/abc/resubmit", nil)
	req.Header.Set("Authorization", "Bearer dynamic-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body: %s", w.Code, w.Body.String())
	}
	if mock.called {
		t.Fatalf("service should not be called when URI is invalid")
	}
}

func TestResubmitVerification_ServiceReturnsForbidden(t *testing.T) {
	mockError := &apperror.Error{
		Code: code.Forbidden,
		Err:  errors.New("action forbidden by business rules"),
	}

	mock := &resubmitVerificationServiceMock{
		err: mockError,
	}

	router := setupResubmitRouter(mock, validHackerUUID, "business")

	req, _ := http.NewRequest(http.MethodPost, "/business/10/resubmit", nil)
	req.Header.Set("Authorization", "Bearer dynamic-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d, body: %s", w.Code, w.Body.String())
	}

	if !mock.called {
		t.Fatalf("service should have been called")
	}
}

func TestResubmitVerification_MissingUserUUID(t *testing.T) {
	mock := &resubmitVerificationServiceMock{}
	router := setupResubmitRouter(mock, "", "business")

	req, _ := http.NewRequest(http.MethodPost, "/business/10/resubmit", nil)
	req.Header.Set("Authorization", "Bearer dynamic-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		t.Fatalf("expected error status, got 200 OK")
	}

	if mock.called {
		t.Fatalf("service should not be called when user UUID extraction fails")
	}
}
