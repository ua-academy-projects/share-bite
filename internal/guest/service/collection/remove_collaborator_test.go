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

func TestRemoveCollaborator(t *testing.T) {
	t.Parallel()

	var (
		collectionID     = "random-collection-uuid"
		customerID       = "random-collection-owner-customer-uuid"
		targetCustomerID = "random-target-customer-uuid"

		collectionName = gofakeit.BeerName()

		errRepo = errors.New("unexpected repository error")
	)

	tests := []struct {
		name  string
		input entity.RemoveCollaboratorInput

		mockFn func(repo *mockCollectionRepository)

		wantErr error
	}{
		{
			name: "success",
			input: entity.RemoveCollaboratorInput{
				CollectionID:     collectionID,
				CustomerID:       customerID,
				TargetCustomerID: targetCustomerID,
			},
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(
						entity.Collection{
							ID:         collectionID,
							CustomerID: customerID,
							Name:       collectionName,
							IsPublic:   false,
						},
						nil,
					).Once()

				repo.On("DeleteCollaborator", mock.Anything, collectionID, targetCustomerID).
					Return(nil).
					Once()
			},
			wantErr: nil,
		},
		{
			name: "error - get collection resource not found",
			input: entity.RemoveCollaboratorInput{
				CollectionID:     "unknown-uuid-not-found-i-promise",
				CustomerID:       customerID,
				TargetCustomerID: targetCustomerID,
			},
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, "unknown-uuid-not-found-i-promise").
					Return(entity.Collection{}, apperror.CollectionNotFoundID("unknown-uuid-not-found-i-promise")).
					Once()
			},
			wantErr: apperror.CollectionNotFoundID("unknown-uuid-not-found-i-promise"),
		},
		{
			name: "error - get collection repository fails",
			input: entity.RemoveCollaboratorInput{
				CollectionID:     collectionID,
				CustomerID:       customerID,
				TargetCustomerID: targetCustomerID,
			},
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{}, errRepo).
					Once()
			},
			wantErr: errRepo,
		},
		{
			name: "error - collection access denied (try out as collaborator)",
			input: entity.RemoveCollaboratorInput{
				CollectionID:     collectionID,
				CustomerID:       "not-owner-but-collaborator-uuid",
				TargetCustomerID: targetCustomerID,
			},
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: customerID,
					}, nil).
					Once()

				repo.On("CheckIfCollaborator", mock.Anything, collectionID, "not-owner-but-collaborator-uuid").
					Return(true, nil).
					Once()
			},
			wantErr: apperror.ErrCollectionAccessDenied,
		},
		{
			name: "error - check if collaborator repository fails",
			input: entity.RemoveCollaboratorInput{
				CollectionID:     collectionID,
				CustomerID:       "not-owner-but-collaborator-uuid",
				TargetCustomerID: targetCustomerID,
			},
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: customerID,
					}, nil).
					Once()

				repo.On("CheckIfCollaborator", mock.Anything, collectionID, "not-owner-but-collaborator-uuid").
					Return(false, errRepo).
					Once()
			},
			wantErr: errRepo,
		},
		{
			name: "error - collection access denied (hidden)",
			input: entity.RemoveCollaboratorInput{
				CollectionID:     collectionID,
				CustomerID:       "random-customer-uuid",
				TargetCustomerID: targetCustomerID,
			},
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: customerID,
					}, nil).
					Once()

				repo.On("CheckIfCollaborator", mock.Anything, collectionID, "random-customer-uuid").
					Return(false, nil).
					Once()
			},
			wantErr: apperror.CollectionNotFoundID(collectionID),
		},
		{
			name: "error - delete collaborator not found",
			input: entity.RemoveCollaboratorInput{
				CollectionID:     collectionID,
				CustomerID:       customerID,
				TargetCustomerID: "not-collaborator-uuid",
			},
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: customerID,
					}, nil).
					Once()

				repo.On("DeleteCollaborator", mock.Anything, collectionID, "not-collaborator-uuid").
					Return(apperror.CollaboratorNotFound("not-collaborator-uuid")).
					Once()
			},
			wantErr: apperror.CollaboratorNotFound("not-collaborator-uuid"),
		},
		{
			name: "error - delete collaborator repository fails",
			input: entity.RemoveCollaboratorInput{
				CollectionID:     collectionID,
				CustomerID:       customerID,
				TargetCustomerID: targetCustomerID,
			},
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: customerID,
					}, nil).
					Once()

				repo.On("DeleteCollaborator", mock.Anything, collectionID, targetCustomerID).
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
			svc := New(repo, nil, nil)
			tt.mockFn(repo)

			err := svc.RemoveCollaborator(context.Background(), tt.input)
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
