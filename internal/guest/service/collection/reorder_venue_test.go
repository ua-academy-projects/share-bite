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

func TestReorderVenue(t *testing.T) {
	t.Parallel()

	var (
		collectionID = gofakeit.UUID()
		customerID   = gofakeit.UUID()
		venueID      = gofakeit.Int64()

		prevVenueID = gofakeit.Int64()
		nextVenueID = gofakeit.Int64()

		errRepo = errors.New("unexpected repository error")
	)

	tests := []struct {
		name string

		input entity.ReorderVenueInput

		mockFn func(repo *mockCollectionRepository, tx *mockTxManager)

		wait    bool
		wantErr error
	}{
		{
			name: "success",
			input: entity.ReorderVenueInput{
				CollectionID: collectionID,
				CustomerID:   customerID,
				VenueID:      venueID,

				PrevVenueID: &prevVenueID,
				NextVenueID: &nextVenueID,
			},
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: customerID}, nil).
					Once()

				repo.On("CheckIfVenueInCollection", mock.Anything, collectionID, venueID).
					Return(true, nil).
					Once()

				repo.On("GetCollectionVenue", mock.Anything, collectionID, prevVenueID).
					Return(entity.CollectionVenue{SortOrder: 100.0}, nil).Once()

				repo.On("GetCollectionVenue", mock.Anything, collectionID, nextVenueID).
					Return(entity.CollectionVenue{SortOrder: 200.0}, nil).Once()

				repo.On("HasVenuesBetween", mock.Anything, collectionID, 100.0, 200.0).
					Return(false, nil).Once()

				repo.On("UpdateVenueSortOrder", mock.Anything, collectionID, venueID, 150.0).
					Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name: "success - at the top",
			input: entity.ReorderVenueInput{
				CollectionID: collectionID,
				CustomerID:   customerID,
				VenueID:      venueID,
				PrevVenueID:  nil,
				NextVenueID:  &nextVenueID,
			},
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: customerID}, nil).
					Once()

				repo.On("CheckIfVenueInCollection", mock.Anything, collectionID, venueID).
					Return(true, nil).
					Once()

				repo.On("GetCollectionVenue", mock.Anything, collectionID, nextVenueID).
					Return(entity.CollectionVenue{SortOrder: 100.0}, nil).
					Once()

				repo.On("UpdateVenueSortOrder", mock.Anything, collectionID, venueID, 50.0).
					Return(nil).
					Once()
			},
			wantErr: nil,
		},
		{
			name: "success - at the bottom",
			input: entity.ReorderVenueInput{
				CollectionID: collectionID,
				CustomerID:   customerID,
				VenueID:      venueID,
				PrevVenueID:  &prevVenueID,
				NextVenueID:  nil,
			},
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: customerID}, nil).
					Once()

				repo.On("CheckIfVenueInCollection", mock.Anything, collectionID, venueID).
					Return(true, nil).
					Once()

				repo.On("GetCollectionVenue", mock.Anything, collectionID, prevVenueID).
					Return(entity.CollectionVenue{SortOrder: 100.0}, nil).
					Once()

				repo.On("UpdateVenueSortOrder", mock.Anything, collectionID, venueID, 200.0).
					Return(nil).
					Once()
			},
			wantErr: nil,
		},
		{
			name: "success - triggers rebalance",
			input: entity.ReorderVenueInput{
				CollectionID: collectionID,
				CustomerID:   customerID,
				VenueID:      venueID,
				PrevVenueID:  &prevVenueID,
				NextVenueID:  &nextVenueID,
			},
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: customerID}, nil).
					Once()

				repo.On("CheckIfVenueInCollection", mock.Anything, collectionID, venueID).
					Return(true, nil).
					Once()

				repo.On("GetCollectionVenue", mock.Anything, collectionID, prevVenueID).
					Return(entity.CollectionVenue{SortOrder: 100.0}, nil).
					Once()

				repo.On("GetCollectionVenue", mock.Anything, collectionID, nextVenueID).
					Return(entity.CollectionVenue{SortOrder: 100.0000000005}, nil).
					Once()

				repo.On("HasVenuesBetween", mock.Anything, collectionID, 100.0, 100.0000000005).
					Return(false, nil).
					Once()

				repo.On("UpdateVenueSortOrder", mock.Anything, collectionID, venueID, 100.00000000025).
					Return(nil).
					Once()

				repo.On("RebalanceCollectionSortOrders", mock.Anything, collectionID).
					Return(nil).
					Once()
			},
			wantErr: nil,
			wait:    true,
		},
		{
			name: "error - access denied",
			input: entity.ReorderVenueInput{
				CollectionID: collectionID,
				CustomerID:   customerID,
				VenueID:      venueID,
				PrevVenueID:  &prevVenueID,
			},
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: gofakeit.UUID()}, nil).
					Once()
			},
			wantErr: apperror.ErrCollectionAccessDenied,
		},
		{
			name: "error - venue not found in collection",
			input: entity.ReorderVenueInput{
				CollectionID: collectionID,
				CustomerID:   customerID,
				VenueID:      venueID,
				PrevVenueID:  &prevVenueID,
				NextVenueID:  &nextVenueID,
			},
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: customerID}, nil).
					Once()

				repo.On("CheckIfVenueInCollection", mock.Anything, collectionID, venueID).
					Return(false, nil).
					Once()
			},
			wantErr: apperror.VenueNotFoundInCollection(venueID),
		},
		{
			name: "error - both prev and next are nil",
			input: entity.ReorderVenueInput{
				CollectionID: collectionID,
				CustomerID:   customerID,
				VenueID:      venueID,
				PrevVenueID:  nil,
				NextVenueID:  nil,
			},
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: customerID}, nil).
					Once()

				repo.On("CheckIfVenueInCollection", mock.Anything, collectionID, venueID).
					Return(true, nil).
					Once()
			},
			wantErr: apperror.ErrInvalidReorderParams,
		},
		{
			name: "error - has venues between",
			input: entity.ReorderVenueInput{
				CollectionID: collectionID,
				CustomerID:   customerID,
				VenueID:      venueID,
				PrevVenueID:  &prevVenueID,
				NextVenueID:  &nextVenueID,
			},
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: customerID}, nil).
					Once()

				repo.On("CheckIfVenueInCollection", mock.Anything, collectionID, venueID).
					Return(true, nil).
					Once()

				repo.On("GetCollectionVenue", mock.Anything, collectionID, prevVenueID).
					Return(entity.CollectionVenue{SortOrder: 100.0}, nil).
					Once()

				repo.On("GetCollectionVenue", mock.Anything, collectionID, nextVenueID).
					Return(entity.CollectionVenue{SortOrder: 200.0}, nil).
					Once()

				repo.On("HasVenuesBetween", mock.Anything, collectionID, 100.0, 200.0).
					Return(true, nil).
					Once()
			},
			wantErr: apperror.ErrInvalidReorderParams,
		},
		{
			name: "error - tx returns repo fail",
			input: entity.ReorderVenueInput{
				CollectionID: collectionID,
				CustomerID:   customerID,
				VenueID:      venueID,
				PrevVenueID:  &prevVenueID,
			},
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{}, errRepo).
					Once()
			},
			wantErr: errRepo,
		},
		{
			name: "error - check if venue in collection fails",
			input: entity.ReorderVenueInput{
				CollectionID: collectionID,
				CustomerID:   customerID,
				VenueID:      venueID,
				PrevVenueID:  &prevVenueID,
			},
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: customerID}, nil).
					Once()

				repo.On("CheckIfVenueInCollection", mock.Anything, collectionID, venueID).
					Return(false, errRepo).
					Once()
			},
			wantErr: errRepo,
		},
		{
			name: "error - prev is same as venue",
			input: entity.ReorderVenueInput{
				CollectionID: collectionID,
				CustomerID:   customerID,
				VenueID:      venueID,
				PrevVenueID:  &venueID,
				NextVenueID:  &nextVenueID,
			},
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: customerID}, nil).
					Once()

				repo.On("CheckIfVenueInCollection", mock.Anything, collectionID, venueID).
					Return(true, nil).
					Once()
			},
			wantErr: apperror.ErrInvalidReorderParams,
		},
		{
			name: "error - next is same as venue",
			input: entity.ReorderVenueInput{
				CollectionID: collectionID,
				CustomerID:   customerID,
				VenueID:      venueID,
				PrevVenueID:  &prevVenueID,
				NextVenueID:  &venueID,
			},
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: customerID}, nil).
					Once()

				repo.On("CheckIfVenueInCollection", mock.Anything, collectionID, venueID).
					Return(true, nil).
					Once()
			},
			wantErr: apperror.ErrInvalidReorderParams,
		},
		{
			name: "error - prev and next are the same",
			input: entity.ReorderVenueInput{
				CollectionID: collectionID,
				CustomerID:   customerID,
				VenueID:      venueID,
				PrevVenueID:  &prevVenueID,
				NextVenueID:  &prevVenueID,
			},
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: customerID}, nil).
					Once()

				repo.On("CheckIfVenueInCollection", mock.Anything, collectionID, venueID).
					Return(true, nil).
					Once()
			},
			wantErr: apperror.ErrInvalidReorderParams,
		},
		{
			name: "error - get prev venue fails",
			input: entity.ReorderVenueInput{
				CollectionID: collectionID,
				CustomerID:   customerID,
				VenueID:      venueID,
				PrevVenueID:  &prevVenueID,
				NextVenueID:  &nextVenueID,
			},
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: customerID}, nil).
					Once()

				repo.On("CheckIfVenueInCollection", mock.Anything, collectionID, venueID).
					Return(true, nil).
					Once()

				repo.On("GetCollectionVenue", mock.Anything, collectionID, prevVenueID).
					Return(entity.CollectionVenue{}, errRepo).
					Once()
			},
			wantErr: errRepo,
		},
		{
			name: "error - get next venue fails",
			input: entity.ReorderVenueInput{
				CollectionID: collectionID,
				CustomerID:   customerID,
				VenueID:      venueID,
				PrevVenueID:  nil,
				NextVenueID:  &nextVenueID,
			},
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: customerID}, nil).
					Once()

				repo.On("CheckIfVenueInCollection", mock.Anything, collectionID, venueID).
					Return(true, nil).
					Once()

				repo.On("GetCollectionVenue", mock.Anything, collectionID, nextVenueID).
					Return(entity.CollectionVenue{}, errRepo).
					Once()
			},
			wantErr: errRepo,
		},
		{
			name: "error - order above >= order below",
			input: entity.ReorderVenueInput{
				CollectionID: collectionID,
				CustomerID:   customerID,
				VenueID:      venueID,
				PrevVenueID:  &prevVenueID,
				NextVenueID:  &nextVenueID,
			},
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: customerID}, nil).
					Once()

				repo.On("CheckIfVenueInCollection", mock.Anything, collectionID, venueID).
					Return(true, nil).
					Once()

				repo.On("GetCollectionVenue", mock.Anything, collectionID, prevVenueID).
					Return(entity.CollectionVenue{SortOrder: 200.0}, nil).
					Once()

				repo.On("GetCollectionVenue", mock.Anything, collectionID, nextVenueID).
					Return(entity.CollectionVenue{SortOrder: 100.0}, nil).
					Once()
			},
			wantErr: apperror.ErrInvalidReorderParams,
		},
		{
			name: "error - has venues between fails",
			input: entity.ReorderVenueInput{
				CollectionID: collectionID,
				CustomerID:   customerID,
				VenueID:      venueID,
				PrevVenueID:  &prevVenueID,
				NextVenueID:  &nextVenueID,
			},
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: customerID}, nil).
					Once()

				repo.On("CheckIfVenueInCollection", mock.Anything, collectionID, venueID).
					Return(true, nil).
					Once()

				repo.On("GetCollectionVenue", mock.Anything, collectionID, prevVenueID).
					Return(entity.CollectionVenue{SortOrder: 100.0}, nil).
					Once()

				repo.On("GetCollectionVenue", mock.Anything, collectionID, nextVenueID).
					Return(entity.CollectionVenue{SortOrder: 200.0}, nil).
					Once()

				repo.On("HasVenuesBetween", mock.Anything, collectionID, 100.0, 200.0).
					Return(false, errRepo).
					Once()
			},
			wantErr: errRepo,
		},
		{
			name: "error - update venue sort order fails",
			input: entity.ReorderVenueInput{
				CollectionID: collectionID,
				CustomerID:   customerID,
				VenueID:      venueID,
				PrevVenueID:  nil,
				NextVenueID:  &nextVenueID,
			},
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: customerID}, nil).
					Once()

				repo.On("CheckIfVenueInCollection", mock.Anything, collectionID, venueID).
					Return(true, nil).
					Once()

				repo.On("GetCollectionVenue", mock.Anything, collectionID, nextVenueID).
					Return(entity.CollectionVenue{SortOrder: 100.0}, nil).
					Once()

				repo.On("UpdateVenueSortOrder", mock.Anything, collectionID, venueID, 50.0).
					Return(errRepo).
					Once()
			},
			wantErr: errRepo,
		},
		{
			name: "success - rebalance fails",
			input: entity.ReorderVenueInput{
				CollectionID: collectionID,
				CustomerID:   customerID,
				VenueID:      venueID,
				PrevVenueID:  &prevVenueID,
				NextVenueID:  &nextVenueID,
			},
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: customerID}, nil).
					Once()

				repo.On("CheckIfVenueInCollection", mock.Anything, collectionID, venueID).
					Return(true, nil).
					Once()

				repo.On("GetCollectionVenue", mock.Anything, collectionID, prevVenueID).
					Return(entity.CollectionVenue{SortOrder: 100.0}, nil).
					Once()

				repo.On("GetCollectionVenue", mock.Anything, collectionID, nextVenueID).
					Return(entity.CollectionVenue{SortOrder: 100.0000000005}, nil).
					Once()

				repo.On("HasVenuesBetween", mock.Anything, collectionID, 100.0, 100.0000000005).
					Return(false, nil).
					Once()

				repo.On("UpdateVenueSortOrder", mock.Anything, collectionID, venueID, 100.00000000025).
					Return(nil).
					Once()

				repo.On("RebalanceCollectionSortOrders", mock.Anything, collectionID).
					Return(errRepo).
					Once()
			},
			wait:    true,
			wantErr: nil,
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

			err := svc.ReorderVenue(context.Background(), tt.input)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				require.NoError(t, err)
			}

			if tt.wait {
				time.Sleep(time.Millisecond * 100)
			}

			repo.AssertExpectations(t)
			txManager.AssertExpectations(t)
		})
	}
}
