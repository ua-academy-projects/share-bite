package business

import (
	"context"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
	"github.com/ua-academy-projects/share-bite/pkg/database/pagination"
)

type handler struct {
	service businessService
}

type businessService interface {
	CreatePost(ctx context.Context, userID string, unitID int, description string, images []*multipart.FileHeader) (*entity.PostWithPhotos, error)
	CheckOwnership(ctx context.Context, userID string, unitID int) error
	UpdatePost(ctx context.Context, postID int64, userID string, content string) (*entity.PostWithPhotos, error)
	DeletePost(ctx context.Context, postID int64, userID string) error

	Get(ctx context.Context, id int) (*entity.OrgUnit, error)
	List(ctx context.Context, brandId, skip, limit int) (pagination.Result[entity.OrgUnit], error)
	GetPosts(ctx context.Context, skip, limit int) (pagination.Result[entity.PostWithPhotos], error)

	CreateLocation(ctx context.Context, brandID int, ownerUserID string, in dto.CreateLocationInput) (*entity.OrgUnit, error)
	UpdateLocation(ctx context.Context, locationID int, ownerUserID string, in dto.UpdateLocationInput) (*entity.OrgUnit, error)
	DeleteLocation(ctx context.Context, locationID int, ownerUserID string) error

	ListNearbyBoxes(ctx context.Context, offset, limit int, lat, lon float64, categoryID *int) (pagination.Result[entity.BoxWithDistance], error)
	GetVenuesByIDs(ctx context.Context, ids []int) ([]entity.OrgUnit, error)

	CreateBox(ctx context.Context, userID string, req dto.CreateBoxRequest) (*entity.Box, error)
	ReserveBox(ctx context.Context, userID string, boxID int64) (*entity.BoxReservation, error)
	Rating(ctx context.Context, id int) (float32, error)
}

func RegisterHandlers(
	r *gin.RouterGroup,
	service businessService,
	parser middleware.AccessTokenParser,
) {
	h := &handler{
		service: service,
	}

	auth := middleware.Auth(parser)

	orgUnits := r.Group("/org-units")
	{
		orgUnits.GET("/:id", h.get)
		orgUnits.GET("/:id/locations", h.list)
		orgUnits.GET("/:id/rating", h.rating)
		orgUnits.POST("/venues", h.getVenuesByIDs)
	}

	r.GET("/posts", h.GetPosts)
	r.GET("/nearby-boxes", h.ListNearbyBoxes)

	businessPosts := r.Group("/posts").
		Use(auth).
		Use(middleware.RequireRoles("business"))
	{
		businessPosts.PUT("/:id", h.UpdatePost)
		businessPosts.DELETE("/:id", h.DeletePost)
		businessPosts.POST("/:id", h.CreatePost)
	}

	businessLocations := r.Group("").
		Use(auth).
		Use(middleware.RequireRoles("business"))
	{
		businessLocations.POST("/:id/locations", h.createLocation)
		businessLocations.PATCH("/locations/:id", h.updateLocation)
		businessLocations.DELETE("/locations/:id", h.deleteLocation)
	}

	boxes := r.Group("/boxes").
		Use(auth).
		Use(middleware.RequireRoles("business"))
	{
		boxes.POST("", h.CreateBox)
	}

	reservations := r.Group("/boxes").
		Use(auth)
	{
		reservations.PATCH("/:boxID/reserve", h.reserveBox)
	}
}

type errorResponse struct {
	Error string `json:"error" example:"not found"`
}
