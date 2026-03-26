package apperror

import (
	"errors"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/guest/error/code"
)

var (
	ErrInvalidJSON       = newError(code.InvalidJSON, "invalid request body format")
	ErrInvalidParam      = newError(code.InvalidRequest, "invalid path parameter")
	ErrInvalidQueryParam = newError(code.InvalidRequest, "invalid query parameter")

	CustomerAlreadyExists = newError(code.AlreadyExists, "customer profile already exists")

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

func PostNotFoundID(postID string) *Error {
	msg := fmt.Sprintf("post with id %q was not found", postID)
	return newError(code.NotFound, msg)
}

func CustomerNotFoundUserID(userID string) *Error {
	msg := fmt.Sprintf("customer with user_id %q was not found", userID)
	return newError(code.NotFound, msg)
}

func CustomerNotFoundUserName(userName string) *Error {
	msg := fmt.Sprintf("customer with username %q was not found", userName)
	return newError(code.NotFound, msg)
}

func CustomerUserNameTaken(userName string) *Error {
	msg := fmt.Sprintf("customer with username %q already exists", userName)
	return newError(code.AlreadyExists, msg)
}
