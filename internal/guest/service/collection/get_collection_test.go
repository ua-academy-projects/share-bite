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

func TestGetCollection(t *testing.T) {
	t.Parallel()
	var (
		collectionID = gofakeit.UUID()
		customerID   = gofakeit.UUID()

		errRepo = errors.New("unexpected database error")
	)

	tests := []struct {
		name string

		collectionID string
		customerID   *string

		mockFn func(repo *mockCollectionRepository)

		wantCollection entity.Collection
		wantErr        error
	}{
		{
			name:         "success",
			collectionID: collectionID,
			customerID:   &customerID,
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: customerID, IsPublic: false}, nil).
					Once()
			},
			wantCollection: entity.Collection{
				ID:         collectionID,
				CustomerID: customerID,
				IsPublic:   false,
			},
			wantErr: nil,
		},
		{
			name:         "success - public collection accessed by non-owner",
			collectionID: collectionID,
			customerID:   strPtr(gofakeit.UUID()),
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: customerID, IsPublic: true}, nil).
					Once()
			},
			wantCollection: entity.Collection{
				ID:         collectionID,
				CustomerID: customerID,
				IsPublic:   true,
			},
			wantErr: nil,
		},
		{
			name:         "success - public collection accessed unauthenticated",
			collectionID: collectionID,
			customerID:   nil,
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: customerID, IsPublic: true}, nil).
					Once()
			},
			wantCollection: entity.Collection{
				ID:         collectionID,
				CustomerID: customerID,
				IsPublic:   true,
			},
			wantErr: nil,
		},
		{
			name:         "error - repository fails",
			collectionID: collectionID,
			customerID:   &customerID,
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{}, errRepo).Once()
			},
			wantCollection: entity.Collection{},
			wantErr:        errRepo,
		},
		{
			name:         "error - access denied",
			collectionID: collectionID,
			customerID:   &customerID,
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: gofakeit.UUID(),
						IsPublic:   false,
					}, nil).Once()
			},
			wantCollection: entity.Collection{},
			wantErr:        apperror.CollectionNotFoundID(collectionID),
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

			collection, err := svc.GetCollection(context.Background(), tt.collectionID, tt.customerID)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.wantCollection, collection)
			repo.AssertExpectations(t)
		})
	}
}
