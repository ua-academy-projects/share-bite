package customer

import (
	"strings"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/pkg/validator"
)

func createRequestToCreateCustomer(req createRequest, userID string) (entity.CreateCustomer, error) {
	var bio *string
	if req.Bio != nil {
		trimmedBio := strings.TrimSpace(*req.Bio)
		if trimmedBio != "" {
			bio = &trimmedBio
		}
	}

	out := entity.CreateCustomer{
		UserID: userID,

		UserName:  strings.ToLower(req.UserName),
		FirstName: strings.TrimSpace(req.FirstName),
		LastName:  strings.TrimSpace(req.LastName),

		Bio: bio,
	}

	var valErrors []validator.ValidationErrorItem
	if len(out.FirstName) < 2 {
		valErrors = append(valErrors, validator.ValidationErrorItem{
			Field:   "firstName",
			Message: "This field must be at least 2 characters long",
		})
	}
	if len(out.LastName) < 2 {
		valErrors = append(valErrors, validator.ValidationErrorItem{
			Field:   "lastName",
			Message: "This field must be at least 2 characters long",
		})
	}

	if len(valErrors) > 0 {
		return entity.CreateCustomer{}, &validator.ValidationError{Errors: valErrors}
	}

	return out, nil
}

func updateRequestToUpdateCustomer(req updateRequest, userID string) (entity.UpdateCustomer, error) {
	var bio *string
	if req.Bio != nil {
		trimmedBio := strings.TrimSpace(*req.Bio)
		if trimmedBio != "" {
			bio = &trimmedBio
		}
	}

	out := entity.UpdateCustomer{
		UserID: userID,

		UserName:  lowerPtr(req.UserName),
		FirstName: trimSpacePtr(req.FirstName),
		LastName:  trimSpacePtr(req.LastName),

		Bio:             bio,
		AvatarObjectKey: req.AvatarObjectKey,

		IsFollowersPublic: req.IsFollowersPublic,
		IsFollowingPublic: req.IsFollowingPublic,
	}

	var valErrors []validator.ValidationErrorItem
	if out.FirstName != nil && len(*out.FirstName) < 2 {
		valErrors = append(valErrors, validator.ValidationErrorItem{
			Field:   "firstName",
			Message: "This field must be at least 2 characters long",
		})
	}
	if out.LastName != nil && len(*out.LastName) < 2 {
		valErrors = append(valErrors, validator.ValidationErrorItem{
			Field:   "lastName",
			Message: "This field must be at least 2 characters long",
		})
	}

	if len(valErrors) > 0 {
		return entity.UpdateCustomer{}, &validator.ValidationError{Errors: valErrors}
	}

	return out, nil
}

func trimSpacePtr(s *string) *string {
	if s == nil {
		return nil
	}
	val := strings.TrimSpace(*s)
	return &val
}

func lowerPtr(s *string) *string {
	if s == nil {
		return nil
	}
	val := strings.ToLower(*s)
	return &val
}
