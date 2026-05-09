package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/guest/error/code"
	"github.com/ua-academy-projects/share-bite/internal/guest/util/response"
	"github.com/ua-academy-projects/share-bite/pkg/validator"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()
		if err == nil {
			return
		}

		if c.Writer.Written() {
			return
		}

		status, resp := mapErrToResp(err.Err)
		c.JSON(status, resp)
	}
}

func mapErrToResp(err error) (int, response.ErrorResponse) {
	status := http.StatusInternalServerError
	resp := response.ErrorResponse{
		Message: "internal server error",
	}

	var validationErr *validator.ValidationError
	if errors.As(err, &validationErr) {
		details := make([]response.ErrorDetail, 0, len(validationErr.Errors))
		for _, e := range validationErr.Errors {
			details = append(details, response.ErrorDetail{
				Field:   e.Field,
				Message: e.Message,
			})
		}
		return http.StatusBadRequest, response.ErrorResponse{
			Message: validationErr.Error(),
			Details: details,
		}
	}

	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		return appErrStatus(appErr.Code), response.ErrorResponse{
			Message: appErr.Error(),
		}
	}

	return status, resp
}

func appErrStatus(c code.Code) int {
	switch c {
	case code.NotFound:
		return http.StatusNotFound
	case code.InvalidJSON, code.InvalidRequest, code.BadRequest, code.EmptyUpdate:
		return http.StatusBadRequest
	case code.UpstreamError:
		return http.StatusBadGateway
	case code.AlreadyExists:
		return http.StatusConflict
	case code.Forbidden:
		return http.StatusForbidden
	case code.TooManyRequests:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}
