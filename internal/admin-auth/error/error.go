package apperror

import (
	"errors"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/admin-auth/error/code"
)

var (
	ErrInvalidResetToken = newError(code.InvalidRequest, "invalid or expired reset token")
)

type Error struct {
	Code code.Code
	Err  error
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func (e *Error) Unwrap() error {
	return e.Err
}

func newError(code code.Code, err string) *Error {
	return &Error{
		Code: code,
		Err:  errors.New(err),
	}
}

func UserNotFoundEmail(email string) *Error {
	msg := fmt.Sprintf("user with email %q was not found", email)
	return newError(code.NotFound, msg)
}
