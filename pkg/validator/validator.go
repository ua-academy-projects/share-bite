package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func (e *ValidationError) Error() string {
	return "request validation failed"
}

type CustomValidator struct {
	v *validator.Validate
}

func New(tag string) *CustomValidator {
	v := validator.New()

	if tag != "" {
		v.SetTagName(tag)
	}

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		tags := []string{"json", "uri", "form"}
		for _, tag := range tags {
			name := strings.SplitN(fld.Tag.Get(tag), ",", 2)[0]
			if name == "-" {
				return ""
			}
			if name != "" {
				return name
			}
		}

		return fld.Name
	})

	return &CustomValidator{
		v: v,
	}
}

func (cv *CustomValidator) ValidateStruct(i any) error {
	if err := cv.v.Struct(i); err != nil {
		var validationErrors []ValidationErrorItem

		var valErrs validator.ValidationErrors
		if errors.As(err, &valErrs) {
			for _, e := range valErrs {
				validationErrors = append(validationErrors, ValidationErrorItem{
					Field:   e.Field(),
					Message: msgForTag(e),
				})
			}
		}

		return &ValidationError{Errors: validationErrors}
	}

	return nil
}

func (cv *CustomValidator) Engine() any {
	return cv.v
}

func msgForTag(err validator.FieldError) string {
	param := err.Param()

	switch err.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return fmt.Sprintf("This field must be at least %s characters long", param)
	case "max":
		return fmt.Sprintf("This field must be at most %s characters long", param)
	case "email":
		return "This field must be a valid email address"
	case "url":
		return "This field must be a valid URL"
	case "uuid":
		return "This field must be a valid UUID"
	case "gte":
		return fmt.Sprintf("This field must be greater than or equal to %s", param)
	case "gt":
		return fmt.Sprintf("This field must be greater than %s", param)
	case "lte":
		return fmt.Sprintf("This field must be less than or equal to %s", param)
	case "e164":
		return "This field must be a valid phone number"
	case "required_without":
		return fmt.Sprintf("This field is required if %s is missing", param)
	case "alphanum":
		return "This field can only contain letters and numbers"
	default:
		return "This field is invalid"
	}
}
