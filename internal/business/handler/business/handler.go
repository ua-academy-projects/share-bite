package business

import (
	"context"
	"errors"
	"mime/multipart"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
	"github.com/ua-academy-projects/share-bite/internal/storage"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/pkg/database/pagination"
	common_middleware "github.com/ua-academy-projects/share-bite/pkg/middleware"
)

const (
	multipartOverheadBytes = 1 << 20 // 1MB overhead
)

type handler struct {
	service businessService
	storage storage.ObjectStorage
}

func (h *handler) extractUserUUID(c *gin.Context) (uuid.UUID, error) {
	userUUID, err := httpctx.GetUserUUID(c)
	if err != nil {
		if errors.Is(err, httpctx.ErrMissingContext) {
			return uuid.Nil, apperror.Unauthorized("unauthorized")
		}
		return uuid.Nil, apperror.Unauthorized("invalid user identity")
	}
	return userUUID, nil
}

type businessService interface {
	CreatePost(ctx context.Context, userID string, unitID int, description string, images []*multipart.FileHeader) (*entity.PostWithPhotos, error)
	CheckOwnership(ctx context.Context, userID string, unitID int) error
	UpdatePost(ctx context.Context, postID int64, userID string, content string) (*entity.PostWithPhotos, error)
	DeletePost(ctx context.Context, postID int64, userID string) error
	Get(ctx context.Context, id int) (*entity.OrgUnit, error)
	ToggleLike(ctx context.Context, postID int64, customerID string) (bool, error)
	GetLikes(ctx context.Context, postID int64, limit, offset int) ([]entity.LikeWithAuthor, error)
	CreateComment(ctx context.Context, postID int64, authorID, content string) (*entity.Comment, error)
	UpdateComment(ctx context.Context, postID, commentID int64, authorID, content string) (*entity.Comment, error)
	DeleteComment(ctx context.Context, postID, commentID int64, authorID string) error
	GetComments(ctx context.Context, postID int64, limit, offset int) ([]entity.CommentWithAuthor, error)
	List(ctx context.Context, brandId, skip, limit int, tags []string) (pagination.Result[entity.OrgUnit], error)
	GetPosts(ctx context.Context, skip, limit int, orgIDs []int) (pagination.Result[entity.PostWithPhotos], error)

	CreateLocation(ctx context.Context, brandID int, ownerUserID string, in dto.CreateLocationInput) (*entity.OrgUnit, error)
	UpdateLocation(ctx context.Context, locationID int, ownerUserID string, in dto.UpdateLocationInput) (*entity.OrgUnit, error)
	DeleteLocation(ctx context.Context, locationID int, ownerUserID string) error

	ListNearbyBoxes(ctx context.Context, offset, limit int, lat, lon float64, categoryID *int, orgID *int) (pagination.Result[entity.BoxWithDistance], error)

	ListLocationTags(ctx context.Context) ([]entity.LocationTag, error)
	GetVenuesByIDs(ctx context.Context, ids []int) ([]entity.OrgUnit, error)

	CreateBox(ctx context.Context, userID string, req dto.CreateBoxRequest, image *multipart.FileHeader) (*entity.Box, error)
	ReserveBox(ctx context.Context, userID string, boxID int64) (*entity.BoxReservation, error)
	Rating(ctx context.Context, id int) (float32, error)

	Create(ctx context.Context, in entity.OrgUnit) (int, error)
	UpdateOrg(ctx context.Context, id int, orgAccountID uuid.UUID, in entity.UpdateOrgUnitInput) (*entity.OrgUnit, error)
	DeleteOrg(ctx context.Context, id int, orgAccountID uuid.UUID) error
	ListNearbyVenues(ctx context.Context, lat, lon float64, skip, limit int) (pagination.Result[entity.OrgUnitWithDistance], error)
	SearchVenues(ctx context.Context, query string, skip, limit int, tags []string) (pagination.Result[entity.OrgUnit], error)

	UploadAvatar(ctx context.Context, id int, orgAccountID uuid.UUID, fileHeader *multipart.FileHeader) (*entity.OrgUnit, error)
	UploadBanner(ctx context.Context, id int, orgAccountID uuid.UUID, fileHeader *multipart.FileHeader) (*entity.OrgUnit, error)
}

func RegisterHandlers(
	r *gin.RouterGroup,
	service businessService,
	parser middleware.AccessTokenParser,
	st storage.ObjectStorage,
) {
	h := &handler{
		service: service,
		storage: st,
	}

	auth := middleware.Auth(parser)
	r.GET("/:id", h.getOrgUnit)

	orgUnits := r.Group("/org-units")
	{
		orgUnits.GET("/:id", h.get)
		orgUnits.GET("/:id/locations", h.list)
		orgUnits.GET("/:id/rating", h.rating)
		orgUnits.POST("/venues", h.getVenuesByIDs)
	}

	r.GET("/posts", h.GetPosts)
	r.GET("/posts/:id/likes", h.GetLikes)
	r.GET("/posts/:id/comments", h.GetComments)
	r.GET("/nearby-boxes", h.ListNearbyBoxes)
	r.GET("/location-tags", h.listLocationTags)
	r.GET("/venues/search", h.searchVenues)

	businessPosts := r.Group("/posts").
		Use(auth).
		Use(middleware.RequireRoles(RoleBusiness)).
		Use(common_middleware.RequireWritableAccountStatus())
	{
		businessPosts.PUT("/:id", h.UpdatePost)
		businessPosts.DELETE("/:id", h.DeletePost)
		businessPosts.POST("/:id", h.CreatePost)
	}

	orgMutations := r.Group("").
		Use(auth).
		Use(middleware.RequireRoles(RoleBusiness)).
		Use(common_middleware.RequireWritableAccountStatus())
	{
		orgMutations.POST("", h.createOrgUnit)
		orgMutations.PUT("/:id", h.updateOrgUnit)
		orgMutations.PATCH("/:id", h.updateOrgUnit)
		orgMutations.POST("/:id/avatar", h.uploadAvatar)
		orgMutations.POST("/:id/banner", h.uploadBanner)
		orgMutations.DELETE("/:id", h.deleteOrgUnit)

	}

	businessLocations := r.Group("").
		Use(auth).
		Use(middleware.RequireRoles(RoleBusiness)).
		Use(common_middleware.RequireWritableAccountStatus())
	{
		businessLocations.POST("/:id/locations", h.createLocation)
		businessLocations.PATCH("/locations/:id", h.updateLocation)
		businessLocations.DELETE("/locations/:id", h.deleteLocation)
	}

	boxes := r.Group("/boxes").
		Use(auth).
		Use(middleware.RequireRoles(RoleBusiness)).
		Use(common_middleware.RequireWritableAccountStatus())
	{
		boxes.POST("", h.CreateBox)
	}

	authenticated := r.Group("/").Use(auth)
	{
		authenticated.POST("/posts/:id/likes", h.ToggleLike)
		authenticated.POST("/posts/:id/comments", h.CreateComment)
		authenticated.PATCH("/posts/:id/comments/:comment_id", h.UpdateComment)
		authenticated.DELETE("/posts/:id/comments/:comment_id", h.DeleteComment)
	}

	r.GET("/locations/nearby", h.ListNearbyVenues)

	reservations := r.Group("/boxes").
		Use(auth)
	{
		reservations.PATCH("/:boxID/reserve", h.reserveBox)
	}
}

func (h *handler) toBrandResponse(ctx context.Context, org *entity.OrgUnit) dto.UpdateOrgResponse {
	return dto.UpdateOrgResponse{
		Id:          org.Id,
		Name:        org.Name,
		Avatar:      h.presign(ctx, org.Avatar),
		Banner:      h.presign(ctx, org.Banner),
		Description: org.Description,
	}
}

func (h *handler) presign(ctx context.Context, urlPtr *string) *string {
	if urlPtr == nil || *urlPtr == "" || h.storage == nil {
		return urlPtr
	}

	url := *urlPtr
	// Extract key from URL
	parts := strings.Split(url, "/")
	var key string
	if strings.Contains(url, "localhost") || strings.Contains(url, "127.0.0.1") {
		if len(parts) >= 5 {
			key = strings.Join(parts[4:], "/")
		}
	} else if len(parts) >= 4 {
		key = strings.Join(parts[3:], "/")
	}

	if key == "" {
		return urlPtr
	}

	presigned, err := h.storage.GetPresignedURL(ctx, key)
	if err != nil {
		return urlPtr
	}

	return &presigned
}

type errorResponse struct {
	Error string `json:"error" example:"not found"`
}

type CreateBoxResponse struct {
	ID      int64  `json:"id"`
	Message string `json:"message"`
}
