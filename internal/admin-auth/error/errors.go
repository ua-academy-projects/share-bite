package error

import (
	"errors"
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Err }

func (e *AppError) HTTPStatus() int { return e.Code }

func (e *AppError) Is(target error) bool {
	var t *AppError
	if errors.As(target, &t) {
		return e.Message == t.Message
	}
	return false
}

func New(code int, msg string) *AppError {
	return &AppError{Code: code, Message: msg}
}

func Wrap(code int, msg string, cause error) *AppError {
	return &AppError{Code: code, Message: msg, Err: cause}
}

var (
	ErrInvalidCredentials = New(http.StatusUnauthorized, "invalid credentials")
	ErrInvalidToken       = New(http.StatusUnauthorized, "invalid or expired token")
	ErrUserNotFound       = New(http.StatusNotFound, "user not found")
	ErrUserAlreadyExists  = New(http.StatusConflict, "user with this email already exists")
	ErrRoleNotFound       = New(http.StatusUnprocessableEntity, "role not found")
	ErrForbidden          = New(http.StatusForbidden, "user doesn't have permission to access this resource")

	ErrProviderExchangeFail  = New(http.StatusBadGateway, "failed to exchange code with provider")
	ErrProviderUserInfoFail  = New(http.StatusBadGateway, "failed to fetch user info from provider")
	ErrProviderAlreadyLinked = New(http.StatusConflict, "social provider already linked to account")
	ErrUnsupportedProvider   = New(http.StatusBadRequest, "unsupported social provider")
	ErrEmailNotVerified      = New(http.StatusForbidden, "email not verified by social provider")
	ErrInvalidResetToken     = New(http.StatusBadRequest, "invalid or expired reset token")
)

func UserNotFoundEmail(email string) *AppError {
	msg := fmt.Sprintf("user with email %q was not found", email)
	return New(http.StatusNotFound, msg)
}
