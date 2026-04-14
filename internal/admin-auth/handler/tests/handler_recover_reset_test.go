package tests

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/dto"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/handler/auth"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/models"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/repository/user"
	authsvc "github.com/ua-academy-projects/share-bite/internal/admin-auth/service/auth"
	"github.com/ua-academy-projects/share-bite/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepository struct {
	mock.Mock
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

func (m *mockUserRepository) CreateUser(ctx context.Context, params user.CreateUser) (*user.CreatedUser, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*user.CreatedUser), args.Error(1)
}

func (m *mockUserRepository) FindBySocialProvider(ctx context.Context, provider string, providerID string) (*dto.UserWithRole, error) {
	args := m.Called(ctx, provider, providerID)
	userWithRole := args.Get(0)
	if userWithRole == nil {
		return nil, args.Error(1)
	}
	return userWithRole.(*dto.UserWithRole), args.Error(1)
}

func (m *mockUserRepository) CreateWithSocial(ctx context.Context, params dto.CreateUserWithSocialParams) (*dto.CreatedUser, error) {
	args := m.Called(ctx, params)
	createdUser := args.Get(0)
	if createdUser == nil {
		return nil, args.Error(1)
	}
	return createdUser.(*dto.CreatedUser), args.Error(1)
}

func (m *mockUserRepository) LinkSocialAccount(ctx context.Context, params dto.CreateSocialAccountParams) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}

func (m *mockUserRepository) CreatePasswordResetToken(ctx context.Context, params dto.CreatePasswordResetTokenParams) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}

func (m *mockUserRepository) ResetPassword(ctx context.Context, tokenHash, passwordHash string) (bool, error) {
	args := m.Called(ctx, tokenHash, passwordHash)
	return args.Bool(0), args.Error(1)
}

func (m *mockUserRepository) AssignRole(ctx context.Context, userID string, roleID int) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}
type stubTokenProvider struct{}

func (s stubTokenProvider) GenerateToken(_ string, _ string) (string, string, error) {
	return "", "", errors.New("unexpected GenerateToken call")
}

func (s stubTokenProvider) ParseRefreshToken(_ string) (string, string, error) {
	return "", "", errors.New("unexpected ParseRefreshToken call")
}

type mockEmailSender struct {
	mock.Mock
}

func (s *mockEmailSender) SendPasswordResetToken(ctx context.Context, toEmail, token string) error {
	args := s.Called(ctx, toEmail, token)
	return args.Error(0)
}

type noopTxManager struct{}

func (noopTxManager) ReadCommitted(ctx context.Context, fn database.Handler) error {
	return fn(ctx)
}

func buildRecoverResetHandler(repo *mockUserRepository, emailSender *mockEmailSender) *auth.Handler {
	service := authsvc.New(repo, stubTokenProvider{}, emailSender, noopTxManager{})
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

func TestRecoverAccessSuccess(t *testing.T) {
	userEmail := "admin@example.com"
	userID := "user-1"

	repo := &mockUserRepository{}
	emailSender := &mockEmailSender{}

	repo.On("FindByEmail", mock.Anything, userEmail).
		Return(&dto.UserWithRole{User: models.User{ID: userID, Email: userEmail}}, nil).
		Once()
	repo.On("CreatePasswordResetToken", mock.Anything, mock.MatchedBy(func(params dto.CreatePasswordResetTokenParams) bool {
		return params.UserID == userID && params.TokenHash != "" && !params.ExpiresAt.IsZero()
	})).
		Return(nil).
		Once()
	emailSender.On("SendPasswordResetToken", mock.Anything, userEmail, mock.MatchedBy(func(token string) bool {
		return token != ""
	})).
		Return(nil).
		Once()

	h := buildRecoverResetHandler(repo, emailSender)
	c, w := buildGinContext(t, http.MethodPost, "/auth/recover-access", map[string]string{"email": userEmail})

	h.RecoverAccess(c)

	if len(c.Errors) != 0 {
		t.Fatalf("expected no gin errors, got: %v", c.Errors.String())
	}
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status code: %d", w.Code)
	}
	repo.AssertExpectations(t)
	emailSender.AssertExpectations(t)

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp["message"] != "If the email exists, recovery instructions have been sent" {
		t.Fatalf("unexpected message: %q", resp["message"])
	}
}

func TestRecoverAccessFindByEmailFails(t *testing.T) {
	repo := &mockUserRepository{}
	emailSender := &mockEmailSender{}

	repo.On("FindByEmail", mock.Anything, "admin@example.com").
		Return((*dto.UserWithRole)(nil), errors.New("db error")).
		Once()

	h := buildRecoverResetHandler(repo, emailSender)
	c, w := buildGinContext(t, http.MethodPost, "/auth/recover-access", map[string]string{"email": "admin@example.com"})

	h.RecoverAccess(c)

	if len(c.Errors) == 0 {
		t.Fatal("expected gin error to be recorded")
	}
	if !strings.Contains(c.Errors.Last().Error(), "find user by email") {
		t.Fatalf("unexpected error: %s", c.Errors.Last().Error())
	}
	if w.Body.Len() != 0 {
		t.Fatalf("expected empty body without error middleware, got: %s", w.Body.String())
	}
	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "CreatePasswordResetToken", mock.Anything, mock.Anything)
	emailSender.AssertNotCalled(t, "SendPasswordResetToken", mock.Anything, mock.Anything, mock.Anything)
}

func TestResetPasswordSuccess(t *testing.T) {
	rawToken := "token-abc"
	newPassword := "new-password-123"
	expectedTokenHash := hashToken(rawToken)

	repo := &mockUserRepository{}
	emailSender := &mockEmailSender{}

	repo.On("ResetPassword", mock.Anything, expectedTokenHash, mock.MatchedBy(func(passwordHash string) bool {
		return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(newPassword)) == nil
	})).
		Return(true, nil).
		Once()

	h := buildRecoverResetHandler(repo, emailSender)
	c, w := buildGinContext(t, http.MethodPost, "/auth/reset-password", map[string]string{
		"token":        rawToken,
		"new_password": newPassword,
	})

	h.ResetPassword(c)

	if len(c.Errors) != 0 {
		t.Fatalf("expected no gin errors, got: %v", c.Errors.String())
	}
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status code: %d", w.Code)
	}
	repo.AssertExpectations(t)
	emailSender.AssertNotCalled(t, "SendPasswordResetToken", mock.Anything, mock.Anything, mock.Anything)

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp["message"] != "password has been reset" {
		t.Fatalf("unexpected message: %q", resp["message"])
	}
}

func TestResetPasswordRepositoryError(t *testing.T) {
	repo := &mockUserRepository{}
	emailSender := &mockEmailSender{}

	repo.On("ResetPassword", mock.Anything, mock.Anything, mock.Anything).
		Return(false, errors.New("db error")).
		Once()

	h := buildRecoverResetHandler(repo, emailSender)
	c, w := buildGinContext(t, http.MethodPost, "/auth/reset-password", map[string]string{
		"token":        "token-abc",
		"new_password": "new-password-123",
	})

	h.ResetPassword(c)

	if len(c.Errors) == 0 {
		t.Fatal("expected gin error to be recorded")
	}
	if !strings.Contains(c.Errors.Last().Error(), "reset password") {
		t.Fatalf("unexpected error: %s", c.Errors.Last().Error())
	}
	if w.Body.Len() != 0 {
		t.Fatalf("expected empty body without error middleware, got: %s", w.Body.String())
	}
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
