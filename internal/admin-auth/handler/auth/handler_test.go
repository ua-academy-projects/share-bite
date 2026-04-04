package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	authsvc "github.com/ua-academy-projects/share-bite/internal/admin-auth/service/auth"
)

type stubAuthService struct {
	registerFunc func(ctx context.Context, email, password, slug string) (*authsvc.Tokens, error)
}

func (s stubAuthService) Login(context.Context, string, string) (*authsvc.Tokens, error) {
	return nil, nil
}

func (s stubAuthService) Register(ctx context.Context, email, password, slug string) (*authsvc.Tokens, error) {
	return s.registerFunc(ctx, email, password, slug)
}

func (s stubAuthService) Refresh(context.Context, string) (*authsvc.Tokens, error) {
	return nil, nil
}

func TestRegisterReturnsTokens(t *testing.T) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	var gotEmail string
	var gotPassword string
	var gotSlug string

	router := gin.New()
	RegisterHandlers(router.Group("/auth"), stubAuthService{
		registerFunc: func(ctx context.Context, email, password, slug string) (*authsvc.Tokens, error) {
			gotEmail = email
			gotPassword = password
			gotSlug = slug

			return &authsvc.Tokens{
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
			}, nil
		},
	})

	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/register",
		strings.NewReader(`{"email":"user@example.com","password":"strongpass123","slug":"user"}`),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	if gotEmail != "user@example.com" || gotPassword != "strongpass123" || gotSlug != "user" {
		t.Fatalf("unexpected register args: email=%q password=%q slug=%q", gotEmail, gotPassword, gotSlug)
	}

	var resp map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if resp["access_token"] != "access-token" {
		t.Fatalf("expected access token %q, got %q", "access-token", resp["access_token"])
	}

	if resp["refresh_token"] != "refresh-token" {
		t.Fatalf("expected refresh token %q, got %q", "refresh-token", resp["refresh_token"])
	}
}

func TestRegisterReturnsBadRequestForInvalidPayload(t *testing.T) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	called := false

	router := gin.New()
	RegisterHandlers(router.Group("/auth"), stubAuthService{
		registerFunc: func(ctx context.Context, email, password, slug string) (*authsvc.Tokens, error) {
			called = true
			return nil, nil
		},
	})

	req := httptest.NewRequest(
		http.MethodPost,
		"/auth/register",
		strings.NewReader(`{"email":"not-an-email","password":"short","slug":"admin"}`),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	if called {
		t.Fatal("register service should not be called for invalid payload")
	}

	var resp map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if resp["message"] == "" {
		t.Fatal("expected validation error message")
	}
}
