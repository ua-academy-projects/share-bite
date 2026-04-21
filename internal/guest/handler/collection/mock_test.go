package collection

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

type mockCollectionService struct {
	mock.Mock
}

func (m *mockCollectionService) CreateCollection(ctx context.Context, in entity.CreateCollectionInput) (entity.Collection, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(entity.Collection), args.Error(1)
}

func (m *mockCollectionService) UpdateCollection(ctx context.Context, in entity.UpdateCollectionInput) (entity.Collection, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(entity.Collection), args.Error(1)
}

func (m *mockCollectionService) DeleteCollection(ctx context.Context, collectionID string, customerID string) error {
	args := m.Called(ctx, collectionID, customerID)
	return args.Error(0)
}

func (m *mockCollectionService) GetCollection(ctx context.Context, collectionID string, customerID *string) (entity.Collection, error) {
	args := m.Called(ctx, collectionID, customerID)
	return args.Get(0).(entity.Collection), args.Error(1)
}

func (m *mockCollectionService) ListCustomerCollections(ctx context.Context, in entity.ListCustomerCollectionsInput) (entity.ListCustomerCollectionsOutput, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(entity.ListCustomerCollectionsOutput), args.Error(1)
}

// collection venues

func (m *mockCollectionService) AddVenue(ctx context.Context, collectionID string, customerID string, venueID int64) error {
	args := m.Called(ctx, collectionID, customerID, venueID)
	return args.Error(0)
}

func (m *mockCollectionService) RemoveVenue(ctx context.Context, collectionID string, customerID string, venueID int64) error {
	args := m.Called(ctx, collectionID, customerID, venueID)
	return args.Error(0)
}

func (m *mockCollectionService) ReorderVenue(ctx context.Context, in entity.ReorderVenueInput) error {
	args := m.Called(ctx, in)
	return args.Error(0)
}

func (m *mockCollectionService) ListVenues(ctx context.Context, collectionID string, customerID *string) ([]entity.EnrichedVenueItem, error) {
	args := m.Called(ctx, collectionID, customerID)
	return args.Get(0).([]entity.EnrichedVenueItem), args.Error(1)
}

func (m *mockCollectionService) ListCollaborators(ctx context.Context, collectionID string, customerID *string) ([]entity.Collaborator, error) {
	args := m.Called(ctx, collectionID, customerID)
	return args.Get(0).([]entity.Collaborator), args.Error(1)
}

func (m *mockCollectionService) AddCollaborator(ctx context.Context, in entity.AddCollaboratorInput) error {
	args := m.Called(ctx, in)
	return args.Error(0)
}

func (m *mockCollectionService) RemoveCollaborator(ctx context.Context, in entity.RemoveCollaboratorInput) error {
	args := m.Called(ctx, in)
	return args.Error(0)
}
