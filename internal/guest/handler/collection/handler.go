package collection

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
)

type handler struct {
	service collectionService
	storage objectStorage
}

type objectStorage interface {
	BuildURL(key string) string
}

type collectionService interface {
	CreateCollection(ctx context.Context, in entity.CreateCollectionInput) (entity.Collection, error)
	UpdateCollection(ctx context.Context, in entity.UpdateCollectionInput) (entity.Collection, error)
	DeleteCollection(ctx context.Context, collectionID string, customerID string) error

	AddVenue(ctx context.Context, collectionID string, customerID string, venueID int64) error
	RemoveVenue(ctx context.Context, collectionID string, customerID string, venueID int64) error
	ReorderVenue(ctx context.Context, in entity.ReorderVenueInput) error

	AddCollaborator(ctx context.Context, in entity.AddCollaboratorInput) error
	RemoveCollaborator(ctx context.Context, in entity.RemoveCollaboratorInput) error

	GetCollection(ctx context.Context, collectionID string, customerID *string) (entity.Collection, error)
	ListCustomerCollections(ctx context.Context, in entity.ListCustomerCollectionsInput) (entity.ListCustomerCollectionsOutput, error)
	ListVenues(ctx context.Context, collectionID string, customerID *string) ([]entity.EnrichedVenueItem, error)
	ListCollaborators(ctx context.Context, collectionID string, customerID *string) ([]entity.Collaborator, error)
}

func RegisterHandlers(
	r *gin.RouterGroup,
	service collectionService,
	authMiddleware gin.HandlerFunc,
	optionalAuthMiddleware gin.HandlerFunc,
	customerMiddleware gin.HandlerFunc,
	st objectStorage,
) {
	h := &handler{
		service: service,
		storage: st,
	}

	// OPTIONAL PROTECTION:
	optional := r.Group("/").Use(optionalAuthMiddleware, customerMiddleware)

	optional.GET("/:collectionId", h.getCollection)
	optional.GET("/:collectionId/venues", h.listVenues)
	optional.GET("/:collectionId/collaborators", h.listCollaborators)

	// TODO: add search for collections

	// PROTECTED:
	protected := r.Group("/").Use(authMiddleware, middleware.RequireRoles("user"), customerMiddleware)

	protected.POST("/", h.createCollection)
	protected.GET("/me", h.listMyCollections)
	protected.PATCH("/:collectionId", h.updateCollection)
	protected.DELETE("/:collectionId", h.deleteCollection)

	protected.POST("/:collectionId/venues/:venueId", h.addVenue)
	protected.DELETE("/:collectionId/venues/:venueId", h.removeVenue)
	protected.POST("/:collectionId/venues/:venueId/reorder", h.reorderVenue)

	protected.POST("/:collectionId/collaborators", h.addCollaborator)
	protected.DELETE("/:collectionId/collaborators/:customerId", h.removeCollaborator)
}

type collectionResponse struct {
	ID string `json:"id"`

	Name        string  `json:"name"`
	Description *string `json:"description"`
	IsPublic    bool    `json:"isPublic"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type enrichedVenueItemResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`

	Description *string `json:"description"`
	AvatarURL   *string `json:"avatarUrl"`
	BannerURL   *string `json:"bannerUrl"`

	SortOrder float64   `json:"sortOrder"`
	AddedAt   time.Time `json:"addedAt"`
}

func collectionToResponse(collection entity.Collection) collectionResponse {
	return collectionResponse{
		ID: collection.ID,

		Name:        collection.Name,
		Description: collection.Description,
		IsPublic:    collection.IsPublic,

		CreatedAt: collection.CreatedAt,
		UpdatedAt: collection.UpdatedAt,
	}
}

func enrichedVenueItemToResponse(item entity.EnrichedVenueItem) enrichedVenueItemResponse {
	return enrichedVenueItemResponse{
		ID:          item.VenueItem.ID,
		Name:        item.VenueItem.Name,
		Description: item.VenueItem.Description,
		AvatarURL:   item.VenueItem.AvatarURL,
		BannerURL:   item.VenueItem.BannerURL,

		SortOrder: item.SortOrder,
		AddedAt:   item.AddedAt,
	}
}

type collaboratorResponse struct {
	CustomerID string  `json:"customerId"`
	UserName   string  `json:"userName"`
	AvatarURL  *string `json:"avatarUrl"`

	AddedAt time.Time `json:"addedAt"`
}

func (h *handler) collaboratorToResponse(collaborator entity.Collaborator) collaboratorResponse {
	var avatarURL *string
	if collaborator.AvatarObjectKey != nil && h.storage != nil {
		url := h.storage.BuildURL(*collaborator.AvatarObjectKey)
		avatarURL = &url
	}

	return collaboratorResponse{
		CustomerID: collaborator.CustomerID,
		UserName:   collaborator.UserName,
		AvatarURL:  avatarURL,
		AddedAt:    collaborator.AddedAt,
	}
}
