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

	ErrForbidden = newError(code.Forbidden, "forbidden: you are not the owner of this resource")

	ErrCollectionFull       = newError(code.InvalidRequest, "collection has reached the maximum limit of 100 venues")
	ErrInvalidReorderParams = newError(code.InvalidRequest, "invalid reorder parameters")
	ErrInvalidPageToken     = newError(code.InvalidRequest, "invalid page token")

	ErrUpstreamError   = newError(code.UpstreamError, "upstream service error")
	ErrInvalidPostData = newError(code.InvalidRequest, "invalid post data")

	ErrCustomerAlreadyExists    = newError(code.AlreadyExists, "customer profile already exists")
	ErrVenueAlreadyInCollection = newError(code.AlreadyExists, "this venue is already in the collection")

	ErrEmptyUpdate = newError(code.EmptyUpdate, "nothing to update")

	ErrCollectionAccessDenied = newError(code.Forbidden, "you are not allowed to manage the collection")

	ErrImageRequired        = newError(code.BadRequest, "image is required")
	ErrStorageNotConfigured = newError(code.Internal, "storage is not configured")
	ErrUnsupportedImageType = newError(code.BadRequest, "unsupported image type. only JPEG and PNG are supported")

	ErrCannotFollowYourself = newError(
		code.InvalidRequest,
		"cannot follow yourself",
	)
	ErrCannotUnfollowYourself = newError(
		code.InvalidRequest,
		"cannot unfollow yourself",
	)
	ErrFollowersListPrivate = newError(
		code.Forbidden,
		"followers list is private",
	)
	ErrFollowingListPrivate = newError(
		code.Forbidden,
		"following list is private",
	)

	ErrFollowNotFound = newError(
		code.NotFound,
		"follow relationship not found",
	)
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

func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}

	return e.Code == t.Code && e.Err.Error() == t.Err.Error()
}

func newError(code code.Code, err string) *Error {
	return &Error{
		Code: code,
		Err:  errors.New(err),
	}
}

func BadRequest(msg string) *Error {
	return newError(code.BadRequest, msg)
}

func Internal(msg string) *Error {
	return newError(code.Internal, msg)
}

func VenueNotFoundID(venueID int64) *Error {
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

func CollectionNotFoundID(collectionID string) *Error {
	msg := fmt.Sprintf("collection with id %q was not found", collectionID)
	return newError(code.NotFound, msg)
}

func VenueNotFoundInCollection(venueID int64) *Error {
	msg := fmt.Sprintf("venue with id %d was not found in the collection", venueID)
	return newError(code.NotFound, msg)
}

func CommentNotFoundID(commentID int64) *Error {
	msg := fmt.Sprintf("comment with id %d was not found", commentID)
	return newError(code.NotFound, msg)
}

func AlreadyFollowing(current, target string) *Error {
	return newError(
		code.AlreadyExists,
		fmt.Sprintf("current %q is already following target %q", current, target),
	)
}
