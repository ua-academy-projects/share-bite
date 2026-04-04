package request

import (
	"errors"

	"github.com/gin-gonic/gin"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/pkg/validator"
)

func BindJSON(c *gin.Context, req any) error {
	if err := c.ShouldBindJSON(req); err != nil {
		var valErr *validator.ValidationError
		if errors.As(err, &valErr) {
			return valErr
		}

		return apperror.ErrInvalidJSON
	}

	return nil
}

func BindUri(c *gin.Context, req any) error {
	if err := c.ShouldBindUri(req); err != nil {
		var valErr *validator.ValidationError
		if errors.As(err, &valErr) {
			return valErr
		}

		return apperror.ErrInvalidParam
	}

	return nil
}

func BindQuery(c *gin.Context, req any) error {
	if err := c.ShouldBindQuery(req); err != nil {
		var valErr *validator.ValidationError
		if errors.As(err, &valErr) {
			return valErr
		}

		return apperror.ErrInvalidQueryParam
	}

	return nil
}
