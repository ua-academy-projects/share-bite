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

func TestDeleteCollection(t *testing.T) {
	t.Parallel()

	var (
		collectionID = gofakeit.UUID()
		customerID   = gofakeit.UUID()

		errRepo = errors.New("repository error")
	)

	tests := []struct {
		name string

		collectionID string
		customerID   string

		mockFn func(repo *mockCollectionRepository)

		wantErr error
	}{
		{
			name:         "success",
			collectionID: collectionID,
			customerID:   customerID,
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(
						entity.Collection{
							ID:         collectionID,
							CustomerID: customerID,
						}, nil).Once()

				repo.On("DeleteCollection", mock.Anything, collectionID).Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name:         "error - get collection repo fails",
			collectionID: collectionID,
			customerID:   customerID,
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{}, errRepo).Once()
			},
			wantErr: errRepo,
		},
		{
			name:         "error - collection access denied (collaborator cannot delete)",
			collectionID: collectionID,
			customerID:   customerID,
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).Return(
					entity.Collection{ID: collectionID, CustomerID: gofakeit.UUID()},
					nil,
				).Once()

				repo.On("CheckIfCollaborator", mock.Anything, collectionID, customerID).
					Return(true, nil).Once()
			},
			wantErr: apperror.ErrCollectionAccessDenied,
		},
		{
			name:         "error - delete collection repo fails",
			collectionID: collectionID,
			customerID:   customerID,
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: customerID,
					}, nil).Once()

				repo.On("DeleteCollection", mock.Anything, collectionID).
					Return(errRepo).Once()
			},
			wantErr: errRepo,
		},
		{
			name:         "error - collection not found (outsider cannot delete)",
			collectionID: collectionID,
			customerID:   customerID,
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
			name:         "error - check if collaborator repo fails",
			collectionID: collectionID,
			customerID:   customerID,
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: gofakeit.UUID()}, nil).
					Once()

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

			tt.mockFn(repo)

			err := svc.DeleteCollection(context.Background(), tt.collectionID, tt.customerID)

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
