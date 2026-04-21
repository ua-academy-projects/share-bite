package collection

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

func TestRemoveVenue(t *testing.T) {
	t.Parallel()

	var (
		collectionID = gofakeit.UUID()
		customerID   = gofakeit.UUID()
		venueID      = gofakeit.Int64()

		errRepo = errors.New("unexpected repo error")
	)

	tests := []struct {
		name string

		collectionID string
		customerID   string
		venueID      int64

		mockFn func(repo *mockCollectionRepository)

		wantErr error
	}{
		{
			name:         "success",
			collectionID: collectionID,
			customerID:   customerID,
			venueID:      venueID,
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: customerID}, nil).
					Once()

				repo.On("CheckIfVenueInCollection", mock.Anything, collectionID, venueID).
					Return(true, nil).
					Once()

				repo.On("RemoveVenue", mock.Anything, collectionID, venueID).
					Return(nil).
					Once()
			},
			wantErr: nil,
		},
		{
			name:         "success - as collaborator",
			collectionID: collectionID,
			customerID:   customerID,
			venueID:      venueID,
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: gofakeit.UUID()}, nil).
					Once()

				repo.On("CheckIfCollaborator", mock.Anything, collectionID, customerID).
					Return(true, nil).Once()

				repo.On("CheckIfVenueInCollection", mock.Anything, collectionID, venueID).
					Return(true, nil).Once()

				repo.On("RemoveVenue", mock.Anything, collectionID, venueID).
					Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name:         "error - check if collaborator repo fails",
			collectionID: collectionID,
			customerID:   customerID,
			venueID:      venueID,
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: "other-owner-id"}, nil).Once()
				repo.On("CheckIfCollaborator", mock.Anything, collectionID, customerID).
					Return(false, errRepo).Once()
			},
			wantErr: errRepo,
		},
		{
			name:         "error - get collection repository fails",
			collectionID: collectionID,
			customerID:   customerID,
			venueID:      venueID,
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{}, errRepo).
					Once()
			},
			wantErr: errRepo,
		},
		{
			name:         "error - get collection not found",
			collectionID: collectionID,
			customerID:   customerID,
			venueID:      venueID,
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{}, apperror.CollectionNotFoundID(collectionID)).
					Once()
			},
			wantErr: apperror.CollectionNotFoundID(collectionID),
		},
		{
			name:         "error - collection not found (outsider)",
			collectionID: collectionID,
			customerID:   customerID,
			venueID:      venueID,
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: gofakeit.UUID()}, nil).
					Once()

				repo.On("CheckIfCollaborator", mock.Anything, collectionID, customerID).
					Return(false, nil).Once()
			},
			wantErr: apperror.CollectionNotFoundID(collectionID),
		},
		{
			name:         "error - collection venue doesn't exist",
			collectionID: collectionID,
			customerID:   customerID,
			venueID:      venueID,
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: customerID}, nil).
					Once()

				repo.On("CheckIfVenueInCollection", mock.Anything, collectionID, venueID).
					Return(false, nil).
					Once()
			},
			wantErr: apperror.VenueNotFoundInCollection(venueID),
		},
		{
			name:         "error - check if venue in collection repository fails",
			collectionID: collectionID,
			customerID:   customerID,
			venueID:      venueID,
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: customerID}, nil).
					Once()

				repo.On("CheckIfVenueInCollection", mock.Anything, collectionID, venueID).
					Return(false, errRepo).
					Once()
			},
			wantErr: errRepo,
		},
		{
			name:         "error - remove venue repository fails",
			collectionID: collectionID,
			customerID:   customerID,
			venueID:      venueID,
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: customerID}, nil).
					Once()

				repo.On("CheckIfVenueInCollection", mock.Anything, collectionID, venueID).
					Return(true, nil).
					Once()

				repo.On("RemoveVenue", mock.Anything, collectionID, venueID).
					Return(errRepo).
					Once()
			},
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

			tt.mockFn(repo)

			err := svc.RemoveVenue(context.Background(), tt.collectionID, tt.customerID, tt.venueID)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}

			repo.AssertExpectations(t)
		})
	}
}
