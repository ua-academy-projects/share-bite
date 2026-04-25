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

	GetCollection(ctx context.Context, collectionID string, customerID *string) (entity.Collection, error)
	ListCustomerCollections(ctx context.Context, in entity.ListCustomerCollectionsInput) (entity.ListCustomerCollectionsOutput, error)
	ListVenues(ctx context.Context, collectionID string, customerID *string) ([]entity.EnrichedVenueItem, error)
	ListCollaborators(ctx context.Context, collectionID string, customerID *string) ([]entity.Collaborator, error)

	InviteCollaborator(ctx context.Context, in entity.InviteCollaboratorInput) error
	AcceptInvitation(ctx context.Context, invitationID string, customerID string) error
	DeclineInvitation(ctx context.Context, invitationID string, customerID string) error

	ListInvitations(ctx context.Context, in entity.ListInvitationsInput) (entity.ListInvitationsOutput, error)

	RemoveCollaborator(ctx context.Context, in entity.RemoveCollaboratorInput) error
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
	optional := r.Group("/").Use(optionalAuthMiddleware, customerMiddleware)

	optional.GET("/:collectionId", h.getCollection)
	optional.GET("/:collectionId/venues", h.listVenues)

	optional.GET("/:collectionId/collaborators", h.listCollaborators)

	protected := r.Group("/").Use(authMiddleware, middleware.RequireRoles("user"), customerMiddleware)

	// collections
	protected.POST("/", h.createCollection)
	protected.GET("/me", h.listMyCollections)
	protected.PATCH("/:collectionId", h.updateCollection)
	protected.DELETE("/:collectionId", h.deleteCollection)

	// TODO: add search for collections

	// venues
	protected.POST("/:collectionId/venues/:venueId", h.addVenue)
	protected.DELETE("/:collectionId/venues/:venueId", h.removeVenue)
	protected.POST("/:collectionId/venues/:venueId/reorder", h.reorderVenue)

	// invitations and collaborators
	protected.POST("/:collectionId/invitations", h.inviteCollaborator)
	protected.DELETE("/:collectionId/collaborators/:customerId", h.removeCollaborator)

	protected.GET("/invitations", h.listInvitations)
	protected.POST("/invitations/:invitationId/accept", h.acceptInvitation)
	protected.POST("/invitations/:invitationId/decline", h.declineInvitation)
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

type invitationResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`

	CollectionID   string `json:"collectionId"`
	CollectionName string `json:"collectionName"`

	Inviter customerInfoResponse `json:"inviter"`
	Invitee customerInfoResponse `json:"invitee"`

	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `json:"createdAt"`
}

type customerInfoResponse struct {
	ID        string  `json:"id"`
	UserName  string  `json:"userName"`
	AvatarURL *string `json:"avatarUrl"`
}

func (h *handler) invitationToResponse(invitation entity.EnrichedInvitation) invitationResponse {
	var (
		inviterAvatarURL *string
		inviteeAvatarURL *string
	)

	if invitation.InviterAvatarObjectKey != nil {
		url := h.storage.BuildURL(*invitation.InviterAvatarObjectKey)
		inviterAvatarURL = &url
	}
	if invitation.InviteeAvatarObjectKey != nil {
		url := h.storage.BuildURL(*invitation.InviteeAvatarObjectKey)
		inviteeAvatarURL = &url
	}

	return invitationResponse{
		ID:     invitation.ID,
		Status: string(invitation.Status),

		CollectionID:   invitation.CollectionID,
		CollectionName: invitation.CollectionName,

		Inviter: customerInfoResponse{
			ID:        invitation.InviterID,
			UserName:  invitation.InviterUserName,
			AvatarURL: inviterAvatarURL,
		},
		Invitee: customerInfoResponse{
			ID:        invitation.InviteeID,
			UserName:  invitation.InviteeUserName,
			AvatarURL: inviteeAvatarURL,
		},

		CreatedAt: invitation.CreatedAt,
		ExpiresAt: invitation.ExpiresAt,
	}
}
