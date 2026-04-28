package collection

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

func TestListCollaborators(t *testing.T) {
	t.Parallel()

	var (
		collectionID = "random-collection-uuid"
		ownerID      = "random-owner-customer-uuid"
		customerID   = "random-customer-uuid"

		collaborators = []entity.Collaborator{
			{CustomerID: "collaborator-uuid-1"},
			{CustomerID: "collaborator-uuid-2"},
		}

		errRepo = errors.New("unexpected repository error")
	)

	tests := []struct {
		name string

		collectionID string
		customerID   *string

		mockFn func(repo *mockCollectionRepository)

		wantErr  error
		wantResp []entity.Collaborator
	}{
		{
			name:         "success - public collection (unauthenticated)",
			collectionID: collectionID,
			customerID:   nil,
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:       collectionID,
						IsPublic: true,
					}, nil).
					Once()

				repo.On("ListCollaborators", mock.Anything, collectionID).
					Return(collaborators, nil).
					Once()
			},
			wantErr:  nil,
			wantResp: collaborators,
		},
		{
			name:         "success - public collection (authenticated)",
			collectionID: collectionID,
			customerID:   &customerID,
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:       collectionID,
						IsPublic: true,
					}, nil).
					Once()

				repo.On("ListCollaborators", mock.Anything, collectionID).
					Return(collaborators, nil).
					Once()
			},
			wantErr:  nil,
			wantResp: collaborators,
		},
		{
			name:         "success - private collection (owner)",
			collectionID: collectionID,
			customerID:   &ownerID,
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: ownerID,
						IsPublic:   false,
					}, nil).
					Once()

				// CheckIfCollaborator is not called because ownerID == collection.CustomerID
				repo.On("ListCollaborators", mock.Anything, collectionID).
					Return(collaborators, nil).
					Once()
			},
			wantErr:  nil,
			wantResp: collaborators,
		},
		{
			name:         "success - private collection (collaborator)",
			collectionID: collectionID,
			customerID:   &customerID,
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: ownerID,
						IsPublic:   false,
					}, nil).
					Once()

				repo.On("CheckIfCollaborator", mock.Anything, collectionID, customerID).
					Return(true, nil).
					Once()

				repo.On("ListCollaborators", mock.Anything, collectionID).
					Return(collaborators, nil).
					Once()
			},
			wantErr:  nil,
			wantResp: collaborators,
		},
		{
			name:         "error - get collection repository fails",
			collectionID: collectionID,
			customerID:   &customerID,
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{}, errRepo).
					Once()
			},
			wantErr:  errRepo,
			wantResp: nil,
		},
		{
			name:         "error - collection not found",
			collectionID: collectionID,
			customerID:   &customerID,
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{}, apperror.CollectionNotFoundID(collectionID)).
					Once()
			},
			wantErr:  apperror.CollectionNotFoundID(collectionID),
			wantResp: nil,
		},
		{
			name:         "error - private collection (unauthenticated)",
			collectionID: collectionID,
			customerID:   nil,
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: ownerID,
						IsPublic:   false,
					}, nil).
					Once()
			},
			wantErr:  apperror.CollectionNotFoundID(collectionID),
			wantResp: nil,
		},
		{
			name:         "error - private collection (outsider)",
			collectionID: collectionID,
			customerID:   &customerID,
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: ownerID,
						IsPublic:   false,
					}, nil).
					Once()

				repo.On("CheckIfCollaborator", mock.Anything, collectionID, customerID).
					Return(false, nil).
					Once()
			},
			wantErr:  apperror.CollectionNotFoundID(collectionID),
			wantResp: nil,
		},
		{
			name:         "error - check if collaborator repository fails",
			collectionID: collectionID,
			customerID:   &customerID,
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: ownerID,
						IsPublic:   false,
					}, nil).
					Once()

				repo.On("CheckIfCollaborator", mock.Anything, collectionID, customerID).
					Return(false, errRepo).
					Once()
			},
			wantErr:  errRepo,
			wantResp: nil,
		},
		{
			name:         "error - list collaborators repository fails",
			collectionID: collectionID,
			customerID:   &ownerID,
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: ownerID,
						IsPublic:   false,
					}, nil).
					Once()

				repo.On("ListCollaborators", mock.Anything, collectionID).
					Return([]entity.Collaborator{}, errRepo).
					Once()
			},
			wantErr:  errRepo,
			wantResp: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := new(mockCollectionRepository)
			svc := New(repo, nil, nil)
			tt.mockFn(repo)

			resp, err := svc.ListCollaborators(context.Background(), tt.collectionID, tt.customerID)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantResp, resp)
			}

			repo.AssertExpectations(t)
		})
	}
}
