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

		respCode := http.StatusInternalServerError
		resp := response.ErrorResponse{
			Message: "internal server error",
		}

		var validationErr *validator.ValidationError
		if errors.As(err, &validationErr) {
			respCode = http.StatusBadRequest

			details := make([]response.ErrorDetail, 0, len(validationErr.Errors))
			for _, e := range validationErr.Errors {
				details = append(details, response.ErrorDetail{
					Field:   e.Field,
					Message: e.Message,
				})
			}
			resp = response.ErrorResponse{
				Message: validationErr.Error(),
				Details: details,
			}

			c.JSON(respCode, resp)
			return
		}

		var appErr *apperror.Error
		if errors.As(err, &appErr) {
			switch appErr.Code {
			case code.NotFound:
				respCode = http.StatusNotFound

			case code.InvalidJSON,
				code.InvalidRequest,
				code.EmptyUpdate:
				respCode = http.StatusBadRequest

			case code.UpstreamError:
				respCode = http.StatusBadGateway

			case code.AlreadyExists:
				respCode = http.StatusConflict

			case code.Forbidden:
				respCode = http.StatusForbidden

			default:
				respCode = http.StatusInternalServerError
			}

			resp.Message = appErr.Error()
		}

		c.JSON(respCode, resp)
	}
}
