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
	ErrUpstreamError     = newError(code.UpstreamError, "upstream service error")
	ErrInvalidPostData   = newError(code.InvalidRequest, "invalid post data")

	ErrCustomerAlreadyExists = newError(code.AlreadyExists, "customer profile already exists")

	ErrEmptyUpdate = newError(code.EmptyUpdate, "nothing to update")

	ErrImageRequired        = newError(code.BadRequest, "image is required")
	ErrStorageNotConfigured = newError(code.Internal, "storage is not configured")
	ErrUnsupportedImageType = newError(code.BadRequest, "unsupported image type. only JPEG and PNG are supported")
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

func VenueNotFoundID(venueID string) *Error {
	msg := fmt.Sprintf("venue with id %q was not found", venueID)
	return newError(code.NotFound, msg)
}

func PostNotFoundID(postID string) *Error {
	msg := fmt.Sprintf("post with id %q was not found", postID)
	return newError(code.NotFound, msg)
}

func CustomerNotFoundUserID(userID string) *Error {
	msg := fmt.Sprintf("customer with user_id %q was not found", userID)
	return newError(code.NotFound, msg)
}

func CustomerNotFoundID(customerID string) *Error {
	msg := fmt.Sprintf("customer with id %q was not found", customerID)
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

func BadRequest(msg string) *Error {
	return newError(code.BadRequest, msg)
}

func Internal(msg string) *Error {
	return newError(code.Internal, msg)
}
