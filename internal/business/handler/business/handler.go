package business

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
)

type handler struct {
	service businessService
}

type businessService interface {
	Get(ctx context.Context, id int) (*entity.OrgUnit, error)
	List(ctx context.Context, brandId, page, limit int) ([]entity.OrgUnit, error)
}

func RegisterHandlers(
	r *gin.RouterGroup,
	service businessService,
) {
	h := &handler{
		service: service,
	}

	r.GET("/:id", h.get)
	r.GET("/:id/locations", h.list)
}

// errorResponse is used for swagger documentation.
type errorResponse struct {
	Error string `json:"error" example:"location with id 999 was not found"`
}
