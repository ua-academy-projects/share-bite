package collection

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
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
			name:         "error - collection access denied",
			collectionID: collectionID,
			customerID:   customerID,
			venueID:      venueID,
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: gofakeit.UUID()}, nil).
					Once()
			},
			wantErr: apperror.ErrCollectionAccessDenied,
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
				assert.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				require.NoError(t, err)
			}

			repo.AssertExpectations(t)
		})
	}
}
