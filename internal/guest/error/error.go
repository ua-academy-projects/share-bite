package apperror

import (
	"errors"

	"github.com/ua-academy-projects/share-bite/internal/guest/error/code"
)

var (
	ErrInvalidJSON       = newError(code.InvalidJSON, "invalid request body format")
	ErrInvalidParam      = newError(code.InvalidRequest, "invalid path parameter")
	ErrInvalidQueryParam = newError(code.InvalidRequest, "invalid query parameter")
	ErrUpstreamError     = newError(code.UpstreamError, "upstream service error")
	ErrInvalidPostData   = newError(code.InvalidRequest, "invalid post data")

	ErrCustomerAlreadyExists = newError(code.AlreadyExists, "customer profile already exists")
	ErrCustomerUserNameTaken = newError(code.AlreadyExists, "customer username already taken")
	ErrCustomerNotFound      = newError(code.NotFound, "customer not found")
	ErrVenueNotFound         = newError(code.NotFound, "venue not found")
	ErrPostNotFound          = newError(code.NotFound, "post not found")

	ErrEmptyUpdate = newError(code.EmptyUpdate, "nothing to update")
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
