package collection

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

func TestAcceptInvitation(t *testing.T) {
	t.Parallel()

	var (
		invitationID = "random-invation-uuid"
		inviteeID    = "random-invitee-customer-uuid"
		collectionID = "random-collection-uuid"

		validCollaboratorsCount = 4
		validExpiresAt          = time.Now().Add(invitationTTL)

		errRepo = errors.New("unexpected repository error")
	)

	tests := []struct {
		name string

		invitationID string
		inviteeID    string

		mockFn func(repo *mockCollectionRepository, tx *mockTxManager)

		wantErr error
	}{
		{
			name:         "success",
			invitationID: invitationID,
			inviteeID:    inviteeID,
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				repo.On("GetInvitation", mock.Anything, invitationID).
					Return(
						entity.Invitation{
							ID:           invitationID,
							InviteeID:    inviteeID,
							CollectionID: collectionID,
						},
						nil,
					).
					Once()

				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(
						entity.Collection{
							ID: collectionID,
						},
						nil,
					).
					Once()

				repo.On("GetInvitationForUpdate", mock.Anything, invitationID).
					Return(
						entity.Invitation{
							ID:           invitationID,
							InviteeID:    inviteeID,
							CollectionID: collectionID,
							Status:       entity.PendingInvitationStatus,
							ExpiresAt:    validExpiresAt,
						},
						nil,
					).
					Once()

				repo.On("CountCollaborators", mock.Anything, collectionID).
					Return(validCollaboratorsCount, nil).
					Once()

				repo.On("CreateCollaborator", mock.Anything, collectionID, inviteeID).
					Return(nil).
					Once()

				repo.On("UpdateInvitationStatus", mock.Anything, invitationID, entity.AcceptedInvitationStatus).
					Return(nil).
					Once()
			},
			wantErr: nil,
		},
		{
			name:         "error - invitation not found",
			invitationID: invitationID,
			inviteeID:    inviteeID,
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				repo.On("GetInvitation", mock.Anything, invitationID).
					Return(
						entity.Invitation{},
						apperror.InvitationNotFoundID(invitationID),
					).
					Once()
			},
			wantErr: apperror.InvitationNotFoundID(invitationID),
		},
		{
			name:         "error - get invitation repository fails",
			invitationID: invitationID,
			inviteeID:    inviteeID,
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				repo.On("GetInvitation", mock.Anything, invitationID).
					Return(
						entity.Invitation{},
						errRepo,
					).
					Once()
			},
			wantErr: errRepo,
		},
		{
			name:         "error - access denied (not found)",
			invitationID: invitationID,
			inviteeID:    inviteeID,
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				repo.On("GetInvitation", mock.Anything, invitationID).
					Return(
						entity.Invitation{
							ID:        invitationID,
							InviteeID: "not-invitee-customer-id",
						},
						nil,
					).
					Once()
			},
			wantErr: apperror.InvitationNotFoundID(invitationID),
		},
		{
			name:         "error - collection not found",
			invitationID: invitationID,
			inviteeID:    inviteeID,
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				repo.On("GetInvitation", mock.Anything, invitationID).
					Return(
						entity.Invitation{
							ID:           invitationID,
							InviteeID:    inviteeID,
							CollectionID: collectionID,
						},
						nil,
					).
					Once()

				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(
						entity.Collection{},
						apperror.CollectionNotFoundID(collectionID),
					).
					Once()
			},
			wantErr: apperror.CollectionNotFoundID(collectionID),
		},
		{
			name:         "error - get collection for update repository fails",
			invitationID: invitationID,
			inviteeID:    inviteeID,
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				repo.On("GetInvitation", mock.Anything, invitationID).
					Return(
						entity.Invitation{
							ID:           invitationID,
							InviteeID:    inviteeID,
							CollectionID: collectionID,
						},
						nil,
					).
					Once()

				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(
						entity.Collection{},
						errRepo,
					).
					Once()
			},
			wantErr: errRepo,
		},
		{
			name:         "error - get invitation for update repository fails",
			invitationID: invitationID,
			inviteeID:    inviteeID,
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				repo.On("GetInvitation", mock.Anything, invitationID).
					Return(
						entity.Invitation{
							ID:           invitationID,
							InviteeID:    inviteeID,
							CollectionID: collectionID,
						},
						nil,
					).
					Once()

				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(
						entity.Collection{
							ID: collectionID,
						},
						nil,
					).
					Once()

				repo.On("GetInvitationForUpdate", mock.Anything, invitationID).
					Return(entity.Invitation{}, errRepo).
					Once()
			},
			wantErr: errRepo,
		},
		{
			name:         "error - invitation already expired",
			invitationID: invitationID,
			inviteeID:    inviteeID,
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				repo.On("GetInvitation", mock.Anything, invitationID).
					Return(
						entity.Invitation{
							ID:           invitationID,
							InviteeID:    inviteeID,
							CollectionID: collectionID,
						},
						nil,
					).
					Once()

				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(
						entity.Collection{
							ID: collectionID,
						},
						nil,
					).
					Once()

				repo.On("GetInvitationForUpdate", mock.Anything, invitationID).
					Return(
						entity.Invitation{
							ID:           invitationID,
							InviteeID:    inviteeID,
							CollectionID: collectionID,
							Status:       entity.PendingInvitationStatus,
							ExpiresAt:    time.Now().Add(-10 * time.Hour),
						},
						nil,
					).
					Once()
			},
			wantErr: apperror.ErrInvitationExpired,
		},
		{
			name:         "error - invitation already processed (accepted)",
			invitationID: invitationID,
			inviteeID:    inviteeID,
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				repo.On("GetInvitation", mock.Anything, invitationID).
					Return(
						entity.Invitation{
							ID:           invitationID,
							InviteeID:    inviteeID,
							CollectionID: collectionID,
						},
						nil,
					).
					Once()

				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(
						entity.Collection{
							ID: collectionID,
						},
						nil,
					).
					Once()

				repo.On("GetInvitationForUpdate", mock.Anything, invitationID).
					Return(
						entity.Invitation{
							ID:           invitationID,
							InviteeID:    inviteeID,
							CollectionID: collectionID,
							Status:       entity.AcceptedInvitationStatus,
							ExpiresAt:    validExpiresAt,
						},
						nil,
					).
					Once()
			},
			wantErr: apperror.ErrInvitationAlreadyProcessed,
		},
		{
			name:         "error - invitation already processed (declined)",
			invitationID: invitationID,
			inviteeID:    inviteeID,
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				repo.On("GetInvitation", mock.Anything, invitationID).
					Return(
						entity.Invitation{
							ID:           invitationID,
							InviteeID:    inviteeID,
							CollectionID: collectionID,
						},
						nil,
					).
					Once()

				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(
						entity.Collection{
							ID: collectionID,
						},
						nil,
					).
					Once()

				repo.On("GetInvitationForUpdate", mock.Anything, invitationID).
					Return(
						entity.Invitation{
							ID:           invitationID,
							InviteeID:    inviteeID,
							CollectionID: collectionID,
							Status:       entity.DeclinedInvitationStatus,
							ExpiresAt:    validExpiresAt,
						},
						nil,
					).
					Once()
			},
			wantErr: apperror.ErrInvitationAlreadyProcessed,
		},
		{
			name:         "error - collection has reached the limit of collaborators",
			invitationID: invitationID,
			inviteeID:    inviteeID,
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				repo.On("GetInvitation", mock.Anything, invitationID).
					Return(
						entity.Invitation{
							ID:           invitationID,
							InviteeID:    inviteeID,
							CollectionID: collectionID,
						},
						nil,
					).
					Once()

				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(
						entity.Collection{
							ID: collectionID,
						},
						nil,
					).
					Once()

				repo.On("GetInvitationForUpdate", mock.Anything, invitationID).
					Return(
						entity.Invitation{
							ID:           invitationID,
							InviteeID:    inviteeID,
							CollectionID: collectionID,
							Status:       entity.PendingInvitationStatus,
							ExpiresAt:    validExpiresAt,
						},
						nil,
					).
					Once()

				repo.On("CountCollaborators", mock.Anything, collectionID).
					Return(maxCollaboratorsPerCollection, nil).
					Once()
			},
			wantErr: apperror.CollectionCollaboratorsLimitReached(maxCollaboratorsPerCollection),
		},
		{
			name:         "error - count collaborators repository fails",
			invitationID: invitationID,
			inviteeID:    inviteeID,
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				repo.On("GetInvitation", mock.Anything, invitationID).
					Return(
						entity.Invitation{
							ID:           invitationID,
							InviteeID:    inviteeID,
							CollectionID: collectionID,
						},
						nil,
					).
					Once()

				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(
						entity.Collection{
							ID: collectionID,
						},
						nil,
					).
					Once()

				repo.On("GetInvitationForUpdate", mock.Anything, invitationID).
					Return(
						entity.Invitation{
							ID:           invitationID,
							InviteeID:    inviteeID,
							CollectionID: collectionID,
							Status:       entity.PendingInvitationStatus,
							ExpiresAt:    validExpiresAt,
						},
						nil,
					).
					Once()

				repo.On("CountCollaborators", mock.Anything, collectionID).
					Return(0, errRepo).
					Once()
			},
			wantErr: errRepo,
		},
		{
			name:         "error - customer is already a collaborator",
			invitationID: invitationID,
			inviteeID:    inviteeID,
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				repo.On("GetInvitation", mock.Anything, invitationID).
					Return(
						entity.Invitation{
							ID:           invitationID,
							InviteeID:    inviteeID,
							CollectionID: collectionID,
						},
						nil,
					).
					Once()

				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(
						entity.Collection{
							ID: collectionID,
						},
						nil,
					).
					Once()

				repo.On("GetInvitationForUpdate", mock.Anything, invitationID).
					Return(
						entity.Invitation{
							ID:           invitationID,
							InviteeID:    inviteeID,
							CollectionID: collectionID,
							Status:       entity.PendingInvitationStatus,
							ExpiresAt:    validExpiresAt,
						},
						nil,
					).
					Once()

				repo.On("CountCollaborators", mock.Anything, collectionID).
					Return(validCollaboratorsCount, nil).
					Once()

				repo.On("CreateCollaborator", mock.Anything, collectionID, inviteeID).
					Return(apperror.CustomerAlreadyCollaborator(inviteeID)).
					Once()
			},
			wantErr: apperror.CustomerAlreadyCollaborator(inviteeID),
		},
		{
			name:         "error - create collaborator repository fails",
			invitationID: invitationID,
			inviteeID:    inviteeID,
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				repo.On("GetInvitation", mock.Anything, invitationID).
					Return(
						entity.Invitation{
							ID:           invitationID,
							InviteeID:    inviteeID,
							CollectionID: collectionID,
						},
						nil,
					).
					Once()

				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(
						entity.Collection{
							ID: collectionID,
						},
						nil,
					).
					Once()

				repo.On("GetInvitationForUpdate", mock.Anything, invitationID).
					Return(
						entity.Invitation{
							ID:           invitationID,
							InviteeID:    inviteeID,
							CollectionID: collectionID,
							Status:       entity.PendingInvitationStatus,
							ExpiresAt:    validExpiresAt,
						},
						nil,
					).
					Once()

				repo.On("CountCollaborators", mock.Anything, collectionID).
					Return(validCollaboratorsCount, nil).
					Once()

				repo.On("CreateCollaborator", mock.Anything, collectionID, inviteeID).
					Return(errRepo).
					Once()
			},
			wantErr: errRepo,
		},
		{
			name:         "error - update invitation status repository fails",
			invitationID: invitationID,
			inviteeID:    inviteeID,
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				repo.On("GetInvitation", mock.Anything, invitationID).
					Return(
						entity.Invitation{
							ID:           invitationID,
							InviteeID:    inviteeID,
							CollectionID: collectionID,
						},
						nil,
					).
					Once()

				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(
						entity.Collection{
							ID: collectionID,
						},
						nil,
					).
					Once()

				repo.On("GetInvitationForUpdate", mock.Anything, invitationID).
					Return(
						entity.Invitation{
							ID:           invitationID,
							InviteeID:    inviteeID,
							CollectionID: collectionID,
							Status:       entity.PendingInvitationStatus,
							ExpiresAt:    validExpiresAt,
						},
						nil,
					).
					Once()

				repo.On("CountCollaborators", mock.Anything, collectionID).
					Return(validCollaboratorsCount, nil).
					Once()

				repo.On("CreateCollaborator", mock.Anything, collectionID, inviteeID).
					Return(nil).
					Once()

				repo.On("UpdateInvitationStatus", mock.Anything, invitationID, entity.AcceptedInvitationStatus).
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
			svc := New(repo, txManager, nil)
			tt.mockFn(repo, txManager)

			err := svc.AcceptInvitation(context.Background(), tt.invitationID, tt.inviteeID)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}

			repo.AssertExpectations(t)
			txManager.AssertExpectations(t)
		})
	}
}
