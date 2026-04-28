package collection

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database"
)

type mockCollectionRepository struct {
	mock.Mock
}

func (m *mockCollectionRepository) CreateCollection(ctx context.Context, in entity.CreateCollectionInput) (entity.Collection, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(entity.Collection), args.Error(1)
}

func (m *mockCollectionRepository) DeleteCollection(ctx context.Context, collectionID string) error {
	args := m.Called(ctx, collectionID)
	return args.Error(0)
}

func (m *mockCollectionRepository) UpdateCollection(ctx context.Context, in entity.UpdateCollectionInput) (entity.Collection, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(entity.Collection), args.Error(1)
}

func (m *mockCollectionRepository) ListCustomerCollections(
	ctx context.Context,
	customerID string,
	cursorTime time.Time,
	cursorID string,
	limit int,
) ([]entity.Collection, error) {
	args := m.Called(ctx, customerID, cursorTime, cursorID, limit)
	return args.Get(0).([]entity.Collection), args.Error(1)
}

func (m *mockCollectionRepository) GetCollection(ctx context.Context, collectionID string) (entity.Collection, error) {
	args := m.Called(ctx, collectionID)
	return args.Get(0).(entity.Collection), args.Error(1)
}

func (m *mockCollectionRepository) GetCollectionForUpdate(ctx context.Context, collectionID string) (entity.Collection, error) {
	args := m.Called(ctx, collectionID)
	return args.Get(0).(entity.Collection), args.Error(1)
}

//

func (m *mockCollectionRepository) CountVenues(ctx context.Context, collectionID string) (int, error) {
	args := m.Called(ctx, collectionID)
	return args.Int(0), args.Error(1)
}

func (m *mockCollectionRepository) GetMaxSortOrder(ctx context.Context, collectionID string) (float64, error) {
	args := m.Called(ctx, collectionID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *mockCollectionRepository) CheckIfVenueInCollection(ctx context.Context, collectionID string, venueID int64) (bool, error) {
	args := m.Called(ctx, collectionID, venueID)
	return args.Bool(0), args.Error(1)
}

func (m *mockCollectionRepository) GetCollectionVenue(ctx context.Context, collectionID string, venueID int64) (entity.CollectionVenue, error) {
	args := m.Called(ctx, collectionID, venueID)
	return args.Get(0).(entity.CollectionVenue), args.Error(1)
}

func (m *mockCollectionRepository) ListCollectionVenues(ctx context.Context, collectionID string) ([]entity.CollectionVenue, error) {
	args := m.Called(ctx, collectionID)
	return args.Get(0).([]entity.CollectionVenue), args.Error(1)
}

func (m *mockCollectionRepository) AddVenue(ctx context.Context, collectionID string, venueID int64, sortOrder float64) error {
	args := m.Called(ctx, collectionID, venueID, sortOrder)
	return args.Error(0)
}

func (m *mockCollectionRepository) RemoveVenue(ctx context.Context, collectionID string, venueID int64) error {
	args := m.Called(ctx, collectionID, venueID)
	return args.Error(0)
}

func (m *mockCollectionRepository) UpdateVenueSortOrder(ctx context.Context, collectionID string, venueID int64, sortOrder float64) error {
	args := m.Called(ctx, collectionID, venueID, sortOrder)
	return args.Error(0)
}

func (m *mockCollectionRepository) RebalanceCollectionSortOrders(ctx context.Context, collectionID string) error {
	args := m.Called(ctx, collectionID)
	return args.Error(0)
}

func (m *mockCollectionRepository) HasVenuesBetween(ctx context.Context, collectionID string, venueID int64, lower float64, upper float64) (bool, error) {
	args := m.Called(ctx, collectionID, venueID, lower, upper)
	return args.Bool(0), args.Error(1)
}

//

func (m *mockCollectionRepository) ListCollaborators(ctx context.Context, collectionID string) ([]entity.Collaborator, error) {
	args := m.Called(ctx, collectionID)
	return args.Get(0).([]entity.Collaborator), args.Error(1)
}

func (m *mockCollectionRepository) CountCollaborators(ctx context.Context, collectionID string) (int, error) {
	args := m.Called(ctx, collectionID)
	return args.Int(0), args.Error(1)
}

func (m *mockCollectionRepository) CreateCollaborator(ctx context.Context, collectionID string, customerID string) error {
	args := m.Called(ctx, collectionID, customerID)
	return args.Error(0)
}

func (m *mockCollectionRepository) DeleteCollaborator(ctx context.Context, collectionID string, customerID string) error {
	args := m.Called(ctx, collectionID, customerID)
	return args.Error(0)
}

func (m *mockCollectionRepository) CheckIfCollaborator(ctx context.Context, collectionID string, customerID string) (bool, error) {
	args := m.Called(ctx, collectionID, customerID)
	return args.Bool(0), args.Error(1)
}

func (m *mockCollectionRepository) GetInvitation(ctx context.Context, invitationID string) (entity.Invitation, error) {
	args := m.Called(ctx, invitationID)
	return args.Get(0).(entity.Invitation), args.Error(1)
}

func (m *mockCollectionRepository) GetInvitationForUpdate(ctx context.Context, invitationID string) (entity.Invitation, error) {
	args := m.Called(ctx, invitationID)
	return args.Get(0).(entity.Invitation), args.Error(1)
}

func (m *mockCollectionRepository) GetInvitationByInvitee(ctx context.Context, collectionID string, inviteeID string) (entity.Invitation, error) {
	args := m.Called(ctx, collectionID, inviteeID)
	return args.Get(0).(entity.Invitation), args.Error(1)
}

func (m *mockCollectionRepository) ListInvitations(ctx context.Context, in entity.ListInvitationsInput) ([]entity.EnrichedInvitation, error) {
	args := m.Called(ctx, in)
	return args.Get(0).([]entity.EnrichedInvitation), args.Error(1)
}

func (m *mockCollectionRepository) CreateInvitation(ctx context.Context, in entity.InviteCollaboratorInput) (string, error) {
	args := m.Called(ctx, in)
	return args.String(0), args.Error(1)
}

func (m *mockCollectionRepository) UpdateInvitationStatus(ctx context.Context, invitationID string, status entity.InvitationStatus) error {
	args := m.Called(ctx, invitationID, status)
	return args.Error(0)
}

func (m *mockCollectionRepository) RefreshInvitation(ctx context.Context, invitationID string, newExpiry time.Time) error {
	args := m.Called(ctx, invitationID, newExpiry)
	return args.Error(0)
}

type mockBusinessClient struct {
	mock.Mock
}

func (m *mockBusinessClient) ListVenuesByIDs(ctx context.Context, venueIDs []int64) (map[int64]entity.Venue, error) {
	args := m.Called(ctx, venueIDs)
	if v := args.Get(0); v != nil {
		return v.(map[int64]entity.Venue), args.Error(1)
	}

	return nil, args.Error(1)
}

type mockTxManager struct {
	mock.Mock
}

func (m *mockTxManager) ReadCommitted(ctx context.Context, fn database.Handler) error {
	args := m.Called(ctx, fn)
	if err := args.Error(0); err != nil {
		return err
	}

	return fn(ctx)
}
