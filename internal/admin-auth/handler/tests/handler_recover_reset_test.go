package tests

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/dto"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/models"
	"github.com/ua-academy-projects/share-bite/internal/admin-auth/pkg"
	"golang.org/x/crypto/bcrypt"
)

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
	if !strings.Contains(c.Errors.Last().Error(), "failed to fetch user") {
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
	userID := "user-1"
	expectedTokenHash := pkg.HashToken(rawToken)

	repo := &mockUserRepository{}
	emailSender := &mockEmailSender{}

	repo.On("ResetPassword", mock.Anything, expectedTokenHash, mock.MatchedBy(func(passwordHash string) bool {
		return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(newPassword)) == nil
	})).Return("user-1", true, nil).Once()

	repo.On("RevokeAllUserTokens", mock.Anything, userID).
		Return(nil).
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
	if resp["message"] != "Password has been reset successfully." {
		t.Fatalf("unexpected message: %q", resp["message"])
	}
}

func TestResetPasswordRepositoryError(t *testing.T) {
	repo := &mockUserRepository{}
	emailSender := &mockEmailSender{}

	repo.On("ResetPassword", mock.Anything, mock.Anything, mock.Anything).
		Return("", false, errors.New("db error")).Once()

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
