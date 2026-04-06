package business

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database/pagination"
)

type handler struct {
	service businessService
}

type businessService interface {
	Get(ctx context.Context, id int) (*entity.OrgUnit, error)
	List(ctx context.Context, brandId, skip, limit int) (pagination.Result[entity.OrgUnit], error)
	GetVenuesByIDs(ctx context.Context, ids []int) ([]entity.OrgUnit, error)
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
	r.POST("/venues", h.getVenuesByIDs)
}

// errorResponse is used for swagger documentation.
type errorResponse struct {
	Error string `json:"error" example:"not found"`
}
