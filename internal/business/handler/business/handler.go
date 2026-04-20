package business

import (
	"context"

	"strconv"

	"github.com/google/uuid"
	"github.com/ua-academy-projects/share-bite/internal/middleware"

	"github.com/gin-gonic/gin"

	"github.com/ua-academy-projects/share-bite/internal/business/entity"
)

type handler struct {
	service businessService
}

type businessService interface {
	Get(ctx context.Context, id int) (*entity.OrgUnit, error)
	List(ctx context.Context, brandId, page, limit int) ([]entity.OrgUnit, error)
	Create(ctx context.Context, in entity.OrgUnit) (int, error)
	UpdateOrg(ctx context.Context, id int, orgAccountID uuid.UUID, in entity.UpdateOrgUnitInput) (*entity.OrgUnit, error)
	DeleteOrg(ctx context.Context, id int, orgAccountID uuid.UUID) error
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
	r.POST("/", h.createOrgUnit)
	r.PUT("/:id", h.updateOrgUnit)
	r.PATCH("/:id", h.updateOrgUnit)
	r.DELETE("/:id", h.deleteOrgUnit)

}

// errorResponse is used for swagger documentation.
type errorResponse struct {
	Error string `json:"error" example:"not found"`
}

func getUserID(c *gin.Context) (int64, bool) {
	val, exists := c.Get(middleware.CtxUserID)
	if !exists {
		return 0, false
	}

	userIDStr, ok := val.(string)
	if !ok {
		return 0, false
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0, false
	}

	return userID, true
}

func checkBusinessRole(c *gin.Context) bool {
	val, exists := c.Get(middleware.CtxUserRole)
	if !exists {
		return false
	}

	role, ok := val.(string)
	if !ok {
		return false
	}

	return role == "business"
}
