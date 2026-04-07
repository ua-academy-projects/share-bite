package error

import (
	"errors"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/business/error/code"
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

func OrgUnitNotFoundID(id int) *Error {
	msg := fmt.Sprintf("org unit with id %d was not found", id)
	return newError(code.NotFound, msg)
}

func LocationNotFoundID(id int) *Error {
	msg := fmt.Sprintf("location with id %d was not found", id)
	return newError(code.NotFound, msg)
}

func BadRequest(msg string) *Error {
	return newError(code.BadRequest, msg)
}
