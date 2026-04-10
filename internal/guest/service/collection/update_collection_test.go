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

func TestUpdateCollection(t *testing.T) {
	var (
		collectionID = gofakeit.UUID()
		customerID   = gofakeit.UUID()

		name        = gofakeit.ProductName()
		description = gofakeit.ProductDescription()

		errRepo = errors.New("unexpected repository error")
	)

	tests := []struct {
		name           string
		input          entity.UpdateCollectionInput
		mockFn         func(repo *mockCollectionRepository)
		wantCollection entity.Collection
		wantErr        error
	}{
		{
			name: "success",
			input: entity.UpdateCollectionInput{
				CollectionID: collectionID,
				CustomerID:   customerID,
				Name:         strPtr(name),
				Description:  strPtr(description),
				IsPublic:     boolPtr(true),
			},
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: customerID}, nil).Once()

				repo.On("UpdateCollection", mock.Anything, entity.UpdateCollectionInput{
					CollectionID: collectionID,
					CustomerID:   customerID,
					Name:         strPtr(name),
					Description:  strPtr(description),
					IsPublic:     boolPtr(true),
				}).
					Return(entity.Collection{
						ID:          collectionID,
						Name:        name,
						Description: &description,
						IsPublic:    true,
						CustomerID:  customerID,
					}, nil).Once()
			},
			wantCollection: entity.Collection{
				ID:          collectionID,
				Name:        name,
				Description: &description,
				IsPublic:    true,
				CustomerID:  customerID,
			},
			wantErr: nil,
		},
		{
			name: "error - get collection fails",
			input: entity.UpdateCollectionInput{
				CollectionID: collectionID,
				CustomerID:   customerID,
				Name:         strPtr(name),
				Description:  strPtr(description),
				IsPublic:     boolPtr(true),
			},
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{}, errRepo).Once()
			},
			wantCollection: entity.Collection{},
			wantErr:        errRepo,
		},
		{
			name: "error - access denied",
			input: entity.UpdateCollectionInput{
				CollectionID: collectionID,
				CustomerID:   customerID,
				Name:         strPtr(name),
				Description:  strPtr(description),
				IsPublic:     boolPtr(true),
			},
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: gofakeit.UUID()}, nil).Once()
			},
			wantCollection: entity.Collection{},
			wantErr:        apperror.ErrCollectionAccessDenied,
		},
		{
			name: "error - update collection repo fails",
			input: entity.UpdateCollectionInput{
				CollectionID: collectionID,
				CustomerID:   customerID,
				Name:         strPtr(name),
				Description:  strPtr(description),
				IsPublic:     boolPtr(true),
			},
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: customerID}, nil).Once()

				repo.On("UpdateCollection", mock.Anything, entity.UpdateCollectionInput{
					CollectionID: collectionID,
					CustomerID:   customerID,
					Name:         strPtr(name),
					Description:  strPtr(description),
					IsPublic:     boolPtr(true),
				}).
					Return(entity.Collection{}, errRepo).Once()
			},
			wantCollection: entity.Collection{},
			wantErr:        errRepo,
		},
		{
			name: "error - empty update",
			input: entity.UpdateCollectionInput{
				CollectionID: collectionID,
				CustomerID:   customerID,
				Name:         nil,
				Description:  nil,
				IsPublic:     nil,
			},
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: customerID}, nil).Once()

				repo.On("UpdateCollection", mock.Anything, entity.UpdateCollectionInput{
					CollectionID: collectionID,
					CustomerID:   customerID,
					Name:         nil,
					Description:  nil,
					IsPublic:     nil,
				}).
					Return(entity.Collection{}, apperror.ErrEmptyUpdate).Once()
			},
			wantCollection: entity.Collection{},
			wantErr:        apperror.ErrEmptyUpdate,
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

			collection, err := svc.UpdateCollection(context.Background(), tt.input)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.wantCollection, collection)
			repo.AssertExpectations(t)
		})
	}
}
