package collection

import (
	"context"
	"fmt"
	"time"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/pkg/database"
)

const (
	maxVenuesPerCollection = 100
	sortOrderGap           = 100.0

	rebalanceGapLimit = 1e-9

	maxCollaboratorsPerCollection = 20

	resendInvitationCooldown = time.Hour * 1
	invitationTTL            = time.Hour * 24 * 7 // 7 days
)

type collectionRepository interface {
	CreateCollection(ctx context.Context, in entity.CreateCollectionInput) (entity.Collection, error)
	DeleteCollection(ctx context.Context, collectionID string) error
	UpdateCollection(ctx context.Context, in entity.UpdateCollectionInput) (entity.Collection, error)

	GetCollection(ctx context.Context, collectionID string) (entity.Collection, error)
	GetCollectionForUpdate(ctx context.Context, collectionID string) (entity.Collection, error)
	ListCustomerCollections(ctx context.Context, customerID string, cursorTime time.Time, cursorID string, limit int) ([]entity.Collection, error)

	AddVenue(ctx context.Context, collectionID string, venueID int64, sortOrder float64) error
	RemoveVenue(ctx context.Context, collectionID string, venueID int64) error
	UpdateVenueSortOrder(ctx context.Context, collectionID string, venueID int64, sortOrder float64) error
	RebalanceCollectionSortOrders(ctx context.Context, collectionID string) error

	HasVenuesBetween(ctx context.Context, collectionID string, venueID int64, lower float64, upper float64) (bool, error)
	CheckIfVenueInCollection(ctx context.Context, collectionID string, venueID int64) (bool, error)
	CountVenues(ctx context.Context, collectionID string) (int, error)
	CountCollaborators(ctx context.Context, collectionID string) (int, error)

	GetCollectionVenue(ctx context.Context, collectionID string, venueID int64) (entity.CollectionVenue, error)
	ListCollectionVenues(ctx context.Context, collectionID string) ([]entity.CollectionVenue, error)
	GetMaxSortOrder(ctx context.Context, collectionID string) (float64, error)

	CreateCollaborator(ctx context.Context, collectionID string, customerID string) error
	DeleteCollaborator(ctx context.Context, collectionID string, customerID string) error
	CheckIfCollaborator(ctx context.Context, collectionID string, customerID string) (bool, error)

	ListCollaborators(ctx context.Context, collectionID string) ([]entity.Collaborator, error)

	CreateInvitation(ctx context.Context, in entity.InviteCollaboratorInput) (string, error)
	UpdateInvitationStatus(ctx context.Context, invitationID string, status entity.InvitationStatus) error
	RefreshInvitation(ctx context.Context, invitationID string, newExpiry time.Time) error

	GetInvitation(ctx context.Context, invitationID string) (entity.Invitation, error)
	GetInvitationForUpdate(ctx context.Context, invitationID string) (entity.Invitation, error)
	GetInvitationByInvitee(ctx context.Context, collectionID string, inviteeID string) (entity.Invitation, error)
	ListInvitations(ctx context.Context, in entity.ListInvitationsInput) ([]entity.EnrichedInvitation, error)
}

type businessClient interface {
	ListVenuesByIDs(ctx context.Context, venueIDs []int64) (map[int64]entity.Venue, error)
}

type service struct {
	collectionRepo collectionRepository

	txManager database.TxManager

	businessClient businessClient
}

func New(
	collectionRepo collectionRepository,
	txManager database.TxManager,
	businessClient businessClient,
) *service {
	return &service{
		collectionRepo: collectionRepo,
		txManager:      txManager,
		businessClient: businessClient,
	}
}

// requireOwner returns nil if the requesting user is the owner.
// Returns AccessDenied for collaborators, NotFound for outsiders.
func (s *service) requireOwner(
	ctx context.Context,
	collectionID string,
	customerID string,
	ownerID string,
) error {
	if customerID == ownerID {
		return nil
	}

	isCollaborator, err := s.collectionRepo.CheckIfCollaborator(ctx, collectionID, customerID)
	if err != nil {
		return fmt.Errorf("check if customer is a collaborator: %w", err)
	}
	if isCollaborator {
		return apperror.ErrCollectionAccessDenied
	}

	return apperror.CollectionNotFoundID(collectionID)
}

// requireCollaborator returns nil if the requesting user is the owner or a collaborator.
// Returns NotFound for outsiders.
func (s *service) requireCollaborator(
	ctx context.Context,
	collectionID string,
	customerID string,
	ownerID string,
) error {
	if customerID == ownerID {
		return nil
	}

	isCollaborator, err := s.collectionRepo.CheckIfCollaborator(ctx, collectionID, customerID)
	if err != nil {
		return fmt.Errorf("check if customer is a collaborator: %w", err)
	}
	if !isCollaborator {
		return apperror.CollectionNotFoundID(collectionID)
	}

	return nil
}
