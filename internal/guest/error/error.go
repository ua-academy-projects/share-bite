package apperror

import (
	"errors"
	"fmt"
	"time"

	"github.com/ua-academy-projects/share-bite/internal/guest/error/code"
)

var (
	ErrInvalidJSON       = newError(code.InvalidJSON, "invalid request body format")
	ErrInvalidParam      = newError(code.InvalidRequest, "invalid path parameter")
	ErrInvalidQueryParam = newError(code.InvalidRequest, "invalid query parameter")

	ErrForbidden = newError(code.Forbidden, "forbidden: you are not the owner of this resource")

	ErrInvalidReorderParams = newError(code.InvalidRequest, "invalid reorder parameters")
	ErrInvalidPageToken     = newError(code.InvalidRequest, "invalid page token")

	ErrUpstreamError   = newError(code.UpstreamError, "upstream service error")
	ErrInvalidPostData = newError(code.InvalidRequest, "invalid post data")

	ErrCustomerAlreadyExists    = newError(code.AlreadyExists, "customer profile already exists")
	ErrVenueAlreadyInCollection = newError(code.AlreadyExists, "this venue is already in the collection")

	ErrEmptyUpdate = newError(code.EmptyUpdate, "nothing to update")

	ErrCollectionAccessDenied = newError(code.Forbidden, "you are not allowed to do that action in the collection")

	ErrImageRequired        = newError(code.BadRequest, "image is required")
	ErrStorageNotConfigured = newError(code.Internal, "storage is not configured")
	ErrUnsupportedImageType = newError(code.BadRequest, "unsupported image type. only JPEG and PNG are supported")
	ErrMultipartFormData    = newError(code.BadRequest, "content type must be multipart/form-data")

	ErrCannotListOthersOutboundInvites = newError(code.Forbidden, "you can only view your own outbound invitations")
	ErrCannotListOthersInboundInvites  = newError(code.Forbidden, "you can only view your own inbound invitations")

	ErrInvitationExpired          = newError(code.BadRequest, "this invitation has expired, please ask for a new one")
	ErrInvitationAlreadyProcessed = newError(code.AlreadyExists, "this invitation has already been processed")
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

func VenueNotFoundID(venueID int64) *Error {
	msg := fmt.Sprintf("venue with id %d was not found", venueID)
	return newError(code.NotFound, msg)
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

func InvalidPostStatusTransition(from, to string) *Error {
	msg := fmt.Sprintf("invalid post status transition: %s -> %s", from, to)
	return newError(code.InvalidRequest, msg)
}

func CollectionVenuesLimitReached(limit int) *Error {
	msg := fmt.Sprintf("collection has reached the limit of %d venues", limit)
	return newError(code.InvalidRequest, msg)
}

func CollaboratorNotFound(customerID string) *Error {
	msg := fmt.Sprintf("customer with id %q is not a collaborator in this collection", customerID)
	return newError(code.NotFound, msg)
}

func CustomerAlreadyCollaborator(customerID string) *Error {
	msg := fmt.Sprintf("customer with id %q is already a collaborator in this collection", customerID)
	return newError(code.AlreadyExists, msg)
}

func CollectionCollaboratorsLimitReached(limit int) *Error {
	msg := fmt.Sprintf("collection has reached the limit of %d collaborators", limit)
	return newError(code.InvalidRequest, msg)
}

func InviteeCustomerNotFoundID(inviteeID string) *Error {
	msg := fmt.Sprintf("invitee customer with id %q does not exist", inviteeID)
	return newError(code.NotFound, msg)
}

func InvitationAlreadySent(collectionID string, inviteeID string) *Error {
	msg := fmt.Sprintf("invitation for collection %q has already been sent to customer %q", collectionID, inviteeID)
	return newError(code.AlreadyExists, msg)
}

func InvitationCooldown(cooldown time.Duration) *Error {
	hours := int(cooldown.Hours())
	minutes := int(cooldown.Minutes()) % 60

	var msg string
	if hours > 0 {
		msg = fmt.Sprintf("please wait %d hour(s) and %d minute(s) before resending this invitation", hours, minutes)
	} else {
		if minutes < 1 {
			minutes = 1
		}
		msg = fmt.Sprintf("please wait %d minute(s) before resending this invitation", minutes)
	}

	return newError(code.TooManyRequests, msg)
}

func InvitationNotFoundID(invitationID string) *Error {
	msg := fmt.Sprintf("invitation with id %q does not exist", invitationID)
	return newError(code.NotFound, msg)
}

func InvitationNotFoundForInvitee(collectionID string, inviteeID string) *Error {
	msg := fmt.Sprintf("invitation for collection %q and invitee %q was not found", collectionID, inviteeID)
	return newError(code.NotFound, msg)
}
