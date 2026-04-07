package apperror

import (
	"errors"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/guest/error/code"
)

var (
	ErrImageRequired        = newError(code.BadRequest, "image is required")
	ErrStorageNotConfigured = newError(code.Internal, "storage is not configured")
	ErrUnsupportedImageType = newError(code.BadRequest, "unsupported image type")
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

func PostNotFoundID(postID string) *Error {
	msg := fmt.Sprintf("post with id %q was not found", postID)
	return newError(code.NotFound, msg)
}

func BadRequest(msg string) *Error {
	return newError(code.BadRequest, msg)
}

func Internal(msg string) *Error {
	return newError(code.Internal, msg)
}