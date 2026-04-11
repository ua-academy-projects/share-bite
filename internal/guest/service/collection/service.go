package collection

import (
	"context"
	"time"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database"
)

type collectionRepository interface {
	CreateCollection(ctx context.Context, in entity.CreateCollectionInput) (entity.Collection, error)
	DeleteCollection(ctx context.Context, collectionID string) error
	UpdateCollection(ctx context.Context, in entity.UpdateCollectionInput) (entity.Collection, error)

	ListCustomerCollections(ctx context.Context, customerID string, cursorTime time.Time, cursorID string, limit int) ([]entity.Collection, error)
	GetCollection(ctx context.Context, collectionID string) (entity.Collection, error)
	GetCollectionForUpdate(ctx context.Context, collectionID string) (entity.Collection, error)

	//

	CountVenues(ctx context.Context, collectionID string) (int, error)
	GetMaxSortOrder(ctx context.Context, collectionID string) (float64, error)
	CheckIfVenueInCollection(ctx context.Context, collectionID string, venueID int64) (bool, error)
	GetCollectionVenue(ctx context.Context, collectionID string, venueID int64) (entity.CollectionVenue, error)

	ListCollectionVenues(ctx context.Context, collectionID string) ([]entity.CollectionVenue, error)

	AddVenue(ctx context.Context, collectionID string, venueID int64, sortOrder float64) error
	RemoveVenue(ctx context.Context, collectionID string, venueID int64) error
	UpdateVenueSortOrder(ctx context.Context, collectionID string, venueID int64, sortOrder float64) error
	RebalanceCollectionSortOrders(ctx context.Context, collectionID string) error
	HasVenuesBetween(ctx context.Context, collectionID string, venueID int64, lower float64, upper float64) (bool, error)
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

func canAccessCollection(collection entity.Collection, customerID *string) bool {
	if collection.IsPublic {
		return true
	}
	return customerID != nil && collection.CustomerID == *customerID
}
