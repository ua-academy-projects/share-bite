package auth

import (
	"github.com/gin-gonic/gin"
)

type handler struct {
	service authService
}

type authService interface {
}

func RegisterHandlers(
	r *gin.RouterGroup,
	service authService,
) {
	h := &handler{
		service: service,
	}
	_ = h
}
