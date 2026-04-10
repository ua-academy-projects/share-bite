package collection

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

func TestListVenues(t *testing.T) {
	t.Parallel()

	var (
		collectionID = gofakeit.UUID()
		customerID   = gofakeit.UUID()

		venueID1 = gofakeit.Int64()
		venueID2 = gofakeit.Int64()

		errRepo = errors.New("unexpected database error")
		now     = time.Now()
	)

	collectionVenue1 := entity.CollectionVenue{CollectionID: collectionID, VenueID: venueID1, SortOrder: 100.0, AddedAt: now}
	collectionVenue2 := entity.CollectionVenue{CollectionID: collectionID, VenueID: venueID2, SortOrder: 200.0, AddedAt: now}

	venue1 := entity.Venue{
		ID:          venueID1,
		Name:        gofakeit.ProductName(),
		Description: strPtr(gofakeit.ProductDescription()),
		AvatarURL:   strPtr(gofakeit.URL()),
		BannerURL:   strPtr(gofakeit.URL()),
	}
	venue2 := entity.Venue{
		ID:          venueID2,
		Name:        gofakeit.ProductName(),
		Description: strPtr(gofakeit.ProductDescription()),
		AvatarURL:   strPtr(gofakeit.URL()),
		BannerURL:   strPtr(gofakeit.URL()),
	}

	enrichedVenue1 := entity.EnrichedVenueItem{VenueItem: venue1, SortOrder: 100.0, AddedAt: now}
	enrichedVenue2 := entity.EnrichedVenueItem{VenueItem: venue2, SortOrder: 200.0, AddedAt: now}

	tests := []struct {
		name string

		collectionID string
		customerID   *string

		mockFn func(repo *mockCollectionRepository, businessClient *mockBusinessClient)

		wantRes []entity.EnrichedVenueItem
		wantErr error
	}{
		{
			name:         "success",
			collectionID: collectionID,
			customerID:   &customerID,
			mockFn: func(repo *mockCollectionRepository, client *mockBusinessClient) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: customerID, IsPublic: false}, nil).Once()

				repo.On("ListCollectionVenues", mock.Anything, collectionID).
					Return([]entity.CollectionVenue{collectionVenue1, collectionVenue2}, nil).Once()

				client.On("ListVenuesByIDs", mock.Anything, []int64{venueID1, venueID2}).
					Return(map[int64]entity.Venue{venueID1: venue1, venueID2: venue2}, nil).Once()
			},
			wantRes: []entity.EnrichedVenueItem{enrichedVenue1, enrichedVenue2},
			wantErr: nil,
		},
		{
			name:         "success - venue deleted in another service",
			collectionID: collectionID,
			customerID:   &customerID,
			mockFn: func(repo *mockCollectionRepository, client *mockBusinessClient) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: customerID, IsPublic: false}, nil).Once()

				repo.On("ListCollectionVenues", mock.Anything, collectionID).
					Return([]entity.CollectionVenue{collectionVenue1, collectionVenue2}, nil).Once()

				client.On("ListVenuesByIDs", mock.Anything, []int64{venueID1, venueID2}).
					Return(map[int64]entity.Venue{venueID1: venue1}, nil).Once()
			},
			wantRes: []entity.EnrichedVenueItem{enrichedVenue1},
			wantErr: nil,
		},
		{
			name:         "success - empty collection venues",
			collectionID: collectionID,
			customerID:   &customerID,
			mockFn: func(repo *mockCollectionRepository, client *mockBusinessClient) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: customerID, IsPublic: false}, nil).Once()

				repo.On("ListCollectionVenues", mock.Anything, collectionID).
					Return([]entity.CollectionVenue{}, nil).Once()
			},
			wantRes: nil,
			wantErr: nil,
		},
		{
			name:         "error - get collection fails",
			collectionID: collectionID,
			customerID:   &customerID,
			mockFn: func(repo *mockCollectionRepository, client *mockBusinessClient) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{}, errRepo).Once()
			},
			wantRes: nil,
			wantErr: errRepo,
		},
		{
			name:         "error - access denied",
			collectionID: collectionID,
			customerID:   &customerID,
			mockFn: func(repo *mockCollectionRepository, client *mockBusinessClient) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: gofakeit.UUID(), IsPublic: false}, nil).Once()
			},
			wantRes: nil,
			wantErr: apperror.CollectionNotFoundID(collectionID),
		},
		{
			name:         "error - list collection venues fails",
			collectionID: collectionID,
			customerID:   &customerID,
			mockFn: func(repo *mockCollectionRepository, client *mockBusinessClient) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: customerID, IsPublic: false}, nil).Once()

				repo.On("ListCollectionVenues", mock.Anything, collectionID).
					Return([]entity.CollectionVenue(nil), errRepo).Once()
			},
			wantRes: nil,
			wantErr: errRepo,
		},
		{
			name:         "error - list venues by ids fails",
			collectionID: collectionID,
			customerID:   &customerID,
			mockFn: func(repo *mockCollectionRepository, client *mockBusinessClient) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: customerID, IsPublic: false}, nil).Once()

				repo.On("ListCollectionVenues", mock.Anything, collectionID).
					Return([]entity.CollectionVenue{collectionVenue1}, nil).Once()

				client.On("ListVenuesByIDs", mock.Anything, []int64{venueID1}).
					Return(map[int64]entity.Venue(nil), errRepo).Once()
			},
			wantRes: nil,
			wantErr: errRepo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := new(mockCollectionRepository)
			txManager := new(mockTxManager)
			businessClient := new(mockBusinessClient)
			svc := New(repo, txManager, businessClient)

			tt.mockFn(repo, businessClient)

			venues, err := svc.ListVenues(context.Background(), tt.collectionID, tt.customerID)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.wantRes, venues)

			repo.AssertExpectations(t)
			businessClient.AssertExpectations(t)
		})
	}
}
