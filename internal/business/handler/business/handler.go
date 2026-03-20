package business

import "github.com/gin-gonic/gin"

type handler struct {
	service businessService
}

type businessService interface {
}

func RegisterHandlers(
	r *gin.RouterGroup,
	service businessService,
) {
	h := &handler{
		service: service,
	}

	_ = h
}
