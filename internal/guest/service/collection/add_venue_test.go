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

func TestAddVenue(t *testing.T) {
	t.Parallel()

	var (
		collectionID = gofakeit.UUID()
		customerID   = gofakeit.UUID()
		venueID      = gofakeit.Int64()

		errRepo = errors.New("unexpected database error")
	)

	tests := []struct {
		name string

		collectionID string
		customerID   string
		venueID      int64

		mockFn func(repo *mockCollectionRepository, tx *mockTxManager)

		wantErr error
	}{
		{
			name:         "success - as owner",
			collectionID: collectionID,
			customerID:   customerID,
			venueID:      venueID,
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).Return(
					entity.Collection{
						ID:         collectionID,
						CustomerID: customerID,
					}, nil).Once()

				repo.On("CountVenues", mock.Anything, collectionID).Return(5, nil).Once()

				repo.On("GetMaxSortOrder", mock.Anything, collectionID).Return(500.0, nil).Once()

				repo.On("AddVenue", mock.Anything, collectionID, venueID, 600.0).Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name:         "success - as collaborator",
			collectionID: collectionID,
			customerID:   customerID,
			venueID:      venueID,
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).Return(
					entity.Collection{
						ID:         collectionID,
						CustomerID: gofakeit.UUID(),
					}, nil).Once()

				repo.On("CheckIfCollaborator", mock.Anything, collectionID, customerID).
					Return(true, nil).
					Once()

				repo.On("CountVenues", mock.Anything, collectionID).Return(5, nil).Once()

				repo.On("GetMaxSortOrder", mock.Anything, collectionID).Return(500.0, nil).Once()

				repo.On("AddVenue", mock.Anything, collectionID, venueID, 600.0).Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name:         "error - collection not found (outsider)",
			collectionID: collectionID,
			customerID:   customerID,
			venueID:      venueID,
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).Return(
					entity.Collection{
						ID:         collectionID,
						CustomerID: gofakeit.UUID(),
					}, nil).Once()

				repo.On("CheckIfCollaborator", mock.Anything, collectionID, customerID).
					Return(false, nil).
					Once()
			},
			wantErr: apperror.CollectionNotFoundID(collectionID),
		},
		{
			name:         "error - collection full",
			collectionID: collectionID,
			customerID:   customerID,
			venueID:      venueID,
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).Return(
					entity.Collection{
						ID:         collectionID,
						CustomerID: customerID,
					}, nil).Once()

				repo.On("CountVenues", mock.Anything, collectionID).Return(maxVenuesPerCollection, nil).Once()
			},
			wantErr: apperror.CollectionVenuesLimitReached(maxVenuesPerCollection),
		},
		{
			name:         "error - get collection repo fails",
			collectionID: collectionID,
			customerID:   customerID,
			venueID:      venueID,
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{}, errRepo)
			},
			wantErr: errRepo,
		},
		{
			name:         "error - count collection venues repo fails",
			collectionID: collectionID,
			customerID:   customerID,
			venueID:      venueID,
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).Return(
					entity.Collection{
						ID:         collectionID,
						CustomerID: customerID,
					}, nil).Once()

				repo.On("CountVenues", mock.Anything, collectionID).Return(0, errRepo).Once()
			},
			wantErr: errRepo,
		},
		{
			name:         "error - get max sort order repo fails",
			collectionID: collectionID,
			customerID:   customerID,
			venueID:      venueID,
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).Return(
					entity.Collection{
						ID:         collectionID,
						CustomerID: customerID,
					}, nil).Once()

				repo.On("CountVenues", mock.Anything, collectionID).Return(10, nil).Once()

				repo.On("GetMaxSortOrder", mock.Anything, collectionID).Return(0.0, errRepo).Once()
			},
			wantErr: errRepo,
		},
		{
			name:         "error - add venue repo fails",
			collectionID: collectionID,
			customerID:   customerID,
			venueID:      venueID,
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).Return(
					entity.Collection{
						ID:         collectionID,
						CustomerID: customerID,
					}, nil).Once()

				repo.On("CountVenues", mock.Anything, collectionID).Return(0, nil).Once()

				repo.On("GetMaxSortOrder", mock.Anything, collectionID).Return(0.0, nil).Once()

				repo.On("AddVenue", mock.Anything, collectionID, venueID, 100.0).Return(errRepo).Once()
			},
			wantErr: errRepo,
		},
		{
			name:         "error - venue already in collection",
			collectionID: collectionID,
			customerID:   customerID,
			venueID:      venueID,
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).Return(
					entity.Collection{ID: collectionID, CustomerID: customerID}, nil).Once()

				repo.On("CountVenues", mock.Anything, collectionID).Return(0, nil).Once()
				repo.On("GetMaxSortOrder", mock.Anything, collectionID).Return(0.0, nil).Once()

				repo.On("AddVenue", mock.Anything, collectionID, venueID, 100.0).
					Return(apperror.ErrVenueAlreadyInCollection).Once()
			},
			wantErr: apperror.ErrVenueAlreadyInCollection,
		},
		{
			name:         "error - check if collaborator repo fails",
			collectionID: collectionID,
			customerID:   customerID,
			venueID:      venueID,
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).Return(
					entity.Collection{ID: collectionID, CustomerID: gofakeit.UUID()}, nil).Once()

				repo.On("CheckIfCollaborator", mock.Anything, collectionID, customerID).
					Return(false, errRepo).Once()
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

			tt.mockFn(repo, txManager)

			err := svc.AddVenue(context.Background(), tt.collectionID, tt.customerID, tt.venueID)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}

			repo.AssertExpectations(t)
			txManager.AssertExpectations(t)
			businessClient.AssertExpectations(t)
		})
	}
}
