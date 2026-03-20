package apperror

import (
	"errors"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/guest/error/code"
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
