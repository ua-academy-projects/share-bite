package tests

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/dto"
	apperr "github.com/ua-academy-projects/share-bite/internal/admin-auth/error"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/handler/auth"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/models"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/repository/user"
	authsvc "github.com/ua-academy-projects/share-bite/internal/admin-auth/service/auth"
	"github.com/ua-academy-projects/share-bite/pkg/database"
)

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) ResetPassword(ctx context.Context, tokenHash, passwordHash string) (string, bool, error) {
	args := m.Called(ctx, tokenHash, passwordHash)
	return args.String(0), args.Bool(1), args.Error(2)
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*dto.UserWithRole, error) {
	args := m.Called(ctx, email)
	userWithRole := args.Get(0)
	if userWithRole == nil {
		return nil, args.Error(1)
	}
	return userWithRole.(*dto.UserWithRole), args.Error(1)
}

func (m *mockUserRepository) FindRoleBySlug(ctx context.Context, slug string) (*models.Role, error) {
	args := m.Called(ctx, slug)
	role := args.Get(0)
	if role == nil {
		return nil, args.Error(1)
	}
	return role.(*models.Role), args.Error(1)
}

func (m *mockUserRepository) FindByID(ctx context.Context, id string) (*dto.UserWithRole, error) {
	args := m.Called(ctx, id)
	userWithRole := args.Get(0)
	if userWithRole == nil {
		return nil, args.Error(1)
	}
	return userWithRole.(*dto.UserWithRole), args.Error(1)
}

func (m *mockUserRepository) CreateUser(ctx context.Context, params user.CreateUser) (*user.CreatedUser, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.CreatedUser), args.Error(1)
}

func (m *mockUserRepository) AssignRole(ctx context.Context, userID string, roleID int) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

func (m *mockUserRepository) FindBySocialProvider(ctx context.Context, provider, providerID string) (*dto.UserWithRole, error) {
	args := m.Called(ctx, provider, providerID)
	res := args.Get(0)
	if res == nil {
		return nil, args.Error(1)
	}
	return res.(*dto.UserWithRole), args.Error(1)
}

func (m *mockUserRepository) CreateWithSocial(ctx context.Context, params dto.CreateUserWithSocialParams) (*dto.CreatedUser, error) {
	args := m.Called(ctx, params)
	res := args.Get(0)
	if res == nil {
		return nil, args.Error(1)
	}
	return res.(*dto.CreatedUser), args.Error(1)
}

func (m *mockUserRepository) LinkSocialAccount(ctx context.Context, params dto.CreateSocialAccountParams) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}

func (m *mockUserRepository) CreatePasswordResetToken(ctx context.Context, params dto.CreatePasswordResetTokenParams) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}

func (m *mockUserRepository) StoreRefreshToken(ctx context.Context, params dto.StoreRefreshTokenParams) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}

func (m *mockUserRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	args := m.Called(ctx, tokenHash)
	return args.Error(0)
}

func (m *mockUserRepository) GetUserIDByRefreshToken(ctx context.Context, tokenHash string) (string, error) {
	args := m.Called(ctx, tokenHash)
	return args.String(0), args.Error(1)
}

func (m *mockUserRepository) RevokeAllUserTokens(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *mockUserRepository) EnforceMaxSessions(ctx context.Context, userID string, maxSessions int) error {
	args := m.Called(ctx, userID, maxSessions)
	return args.Error(0)
}

func (m *mockUserRepository) DeleteExpiredTokens(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type MockAuthService struct{ mock.Mock }

func (m *MockAuthService) Login(ctx context.Context, email, password string) (*authsvc.Tokens, error) {
	args := m.Called(ctx, email, password)
	if args.Get(0) != nil {
		return args.Get(0).(*authsvc.Tokens), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthService) Register(ctx context.Context, email, password, slug string) (*authsvc.Tokens, error) {
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

func (m *MockAuthService) Logout(ctx context.Context, userID string, refreshToken string) error {
	args := m.Called(ctx, userID, refreshToken)
	return args.Error(0)
}

func (m *MockAuthService) RevokeAllSessions(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockAuthService) OAuthLogin(ctx context.Context, provider authsvc.OAuthProvider, code, slug string) (*authsvc.Tokens, error) {
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

func (m *MockAuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
	args := m.Called(ctx, token, newPassword)
	return args.Error(0)
}

type MockProviderFactory struct{ mock.Mock }

func (m *MockProviderFactory) Get(name string) (authsvc.OAuthProvider, error) {
	args := m.Called(name)
	if args.Get(0) != nil {
		return args.Get(0).(authsvc.OAuthProvider), args.Error(1)
	}
	return nil, args.Error(1)
}

type MockOAuthProvider struct{ mock.Mock }

func (m *MockOAuthProvider) ExchangeCode(ctx context.Context, code string) (*dto.OAuthUserInfo, error) {
	args := m.Called(ctx, code)
	if args.Get(0) != nil {
		return args.Get(0).(*dto.OAuthUserInfo), args.Error(1)
	}
	return nil, args.Error(1)
}

type mockEmailSender struct {
	mock.Mock
}

func (s *mockEmailSender) SendPasswordResetToken(ctx context.Context, toEmail, token string) error {
	args := s.Called(ctx, toEmail, token)
	return args.Error(0)
}

type stubTokenProvider struct{}

func (s stubTokenProvider) GenerateToken(_ string, _ string) (string, string, error) {
	return "", "", errors.New("unexpected GenerateToken call")
}

func (s stubTokenProvider) ParseRefreshToken(_ string) (string, string, error) {
	return "", "", errors.New("unexpected ParseRefreshToken call")
}

func (s stubTokenProvider) GetRefreshTTL() time.Duration {
	return time.Hour * 24
}

type noopTxManager struct{}

func (noopTxManager) ReadCommitted(ctx context.Context, fn database.Handler) error {
	return fn(ctx)
}

func buildRecoverResetHandler(repo *mockUserRepository, emailSender *mockEmailSender) *auth.Handler {
	service := authsvc.New(repo, stubTokenProvider{}, emailSender, noopTxManager{}, time.Hour, 5)
	return auth.NewHandler(service, nil)
}

func buildGinContext(t *testing.T, method, target string, body any) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal body: %v", err)
	}

	c.Request = httptest.NewRequest(method, target, strings.NewReader(string(payload)))
	c.Request.Header.Set("Content-Type", "application/json")

	return c, w
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
