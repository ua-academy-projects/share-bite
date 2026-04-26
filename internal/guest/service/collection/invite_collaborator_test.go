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

func TestInviteCollaborator(t *testing.T) {
	t.Parallel()

	var (
		collectionID = "random-collection-uuid"
		inviterID    = "random-inviter-customer-uuid"
		inviteeID    = "random-invitee-customer-uuid"
		invitationID = "invitation-uuid"

		expiry = time.Now().UTC().Add(invitationTTL)

		validCollaboratorsCount = 4

		errRepo = errors.New("unexpected repository error")
		_       = errRepo
		_       = expiry
	)

	tests := []struct {
		name  string
		input entity.InviteCollaboratorInput

		mockFn func(repo *mockCollectionRepository, tx *mockTxManager)

		wantErr error
	}{
		{
			name: "success - new invitation",
			input: entity.InviteCollaboratorInput{
				CollectionID: collectionID,
				InviterID:    inviterID,
				InviteeID:    inviteeID,
			},
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(
						entity.Collection{
							ID:         collectionID,
							CustomerID: inviterID,
						},
						nil,
					).
					Once()

				repo.On("CheckIfCollaborator", mock.Anything, collectionID, inviteeID).
					Return(false, nil).
					Once()

				repo.On("CountCollaborators", mock.Anything, collectionID).
					Return(validCollaboratorsCount, nil).
					Once()

				repo.On("GetInvitationByInvitee", mock.Anything, collectionID, inviteeID).
					Return(
						entity.Invitation{},
						apperror.InvitationNotFoundForInvitee(collectionID, inviteeID),
					).
					Once()

				repo.On("CreateInvitation", mock.Anything, mock.MatchedBy(func(in entity.InviteCollaboratorInput) bool {
					return in.CollectionID == collectionID &&
						in.InviterID == inviterID &&
						in.InviteeID == inviteeID &&
						in.Expiry.After(time.Now())
				})).
					Return("new-invitation-uuid", nil).
					Once()
			},
			wantErr: nil,
		},
		{
			name: "success - refresh invitation",
			input: entity.InviteCollaboratorInput{
				CollectionID: collectionID,
				InviterID:    inviterID,
				InviteeID:    inviteeID,
			},
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(
						entity.Collection{
							ID:         collectionID,
							CustomerID: inviterID,
						},
						nil,
					).
					Once()

				repo.On("CheckIfCollaborator", mock.Anything, collectionID, inviteeID).
					Return(false, nil).
					Once()

				repo.On("CountCollaborators", mock.Anything, collectionID).
					Return(validCollaboratorsCount, nil).
					Once()

				repo.On("GetInvitationByInvitee", mock.Anything, collectionID, inviteeID).
					Return(
						entity.Invitation{
							ID:         invitationID,
							LastSentAt: time.Now().Add(-resendInvitationCooldown),
						},
						nil,
					).
					Once()

				repo.On("RefreshInvitation", mock.Anything, invitationID, mock.MatchedBy(func(t time.Time) bool {
					return t.After(time.Now())
				})).
					Return(nil).
					Once()
			},
			wantErr: nil,
		},
		{
			name: "error - get collection for update repository fails",
			input: entity.InviteCollaboratorInput{
				CollectionID: collectionID,
				InviterID:    inviterID,
				InviteeID:    inviteeID,
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
			name: "error - collection not found",
			input: entity.InviteCollaboratorInput{
				CollectionID: collectionID,
				InviterID:    inviterID,
				InviteeID:    inviteeID,
			},
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{}, apperror.CollectionNotFoundID(collectionID)).
					Once()
			},
			wantErr: apperror.CollectionNotFoundID(collectionID),
		},
		{
			name: "error - inviter is not the owner (collaborator)",
			input: entity.InviteCollaboratorInput{
				CollectionID: collectionID,
				InviterID:    "not-owner-but-collaborator-uuid",
				InviteeID:    inviteeID,
			},
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: inviterID,
					}, nil).
					Once()

				repo.On("CheckIfCollaborator", mock.Anything, collectionID, "not-owner-but-collaborator-uuid").
					Return(true, nil).
					Once()
			},
			wantErr: apperror.ErrCollectionAccessDenied,
		},
		{
			name: "error - inviter is not the owner (outsider)",
			input: entity.InviteCollaboratorInput{
				CollectionID: collectionID,
				InviterID:    "outsider-uuid",
				InviteeID:    inviteeID,
			},
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: inviterID,
					}, nil).
					Once()

				repo.On("CheckIfCollaborator", mock.Anything, collectionID, "outsider-uuid").
					Return(false, nil).
					Once()
			},
			wantErr: apperror.CollectionNotFoundID(collectionID),
		},
		{
			name: "error - invitee is the owner",
			input: entity.InviteCollaboratorInput{
				CollectionID: collectionID,
				InviterID:    inviterID,
				InviteeID:    inviterID,
			},
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: inviterID,
					}, nil).
					Once()

				// CheckIfCollaborator is not called because the check
				// `in.InviteeID == collection.CustomerID` short-circuits before it
			},
			wantErr: apperror.CustomerAlreadyCollaborator(inviterID),
		},
		{
			name: "error - invitee is already a collaborator",
			input: entity.InviteCollaboratorInput{
				CollectionID: collectionID,
				InviterID:    inviterID,
				InviteeID:    inviteeID,
			},
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: inviterID,
					}, nil).
					Once()

				repo.On("CheckIfCollaborator", mock.Anything, collectionID, inviteeID).
					Return(true, nil).
					Once()
			},
			wantErr: apperror.CustomerAlreadyCollaborator(inviteeID),
		},
		{
			name: "error - check if collaborator repository fails",
			input: entity.InviteCollaboratorInput{
				CollectionID: collectionID,
				InviterID:    inviterID,
				InviteeID:    inviteeID,
			},
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: inviterID,
					}, nil).
					Once()

				repo.On("CheckIfCollaborator", mock.Anything, collectionID, inviteeID).
					Return(false, errRepo).
					Once()
			},
			wantErr: errRepo,
		},
		{
			name: "error - collaborators limit reached",
			input: entity.InviteCollaboratorInput{
				CollectionID: collectionID,
				InviterID:    inviterID,
				InviteeID:    inviteeID,
			},
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: inviterID,
					}, nil).
					Once()

				repo.On("CheckIfCollaborator", mock.Anything, collectionID, inviteeID).
					Return(false, nil).
					Once()

				repo.On("CountCollaborators", mock.Anything, collectionID).
					Return(maxCollaboratorsPerCollection, nil).
					Once()
			},
			wantErr: apperror.CollectionCollaboratorsLimitReached(maxCollaboratorsPerCollection),
		},
		{
			name: "error - count collaborators repository fails",
			input: entity.InviteCollaboratorInput{
				CollectionID: collectionID,
				InviterID:    inviterID,
				InviteeID:    inviteeID,
			},
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: inviterID,
					}, nil).
					Once()

				repo.On("CheckIfCollaborator", mock.Anything, collectionID, inviteeID).
					Return(false, nil).
					Once()

				repo.On("CountCollaborators", mock.Anything, collectionID).
					Return(0, errRepo).
					Once()
			},
			wantErr: errRepo,
		},
		{
			name: "error - get invitation by invitee repository fails",
			input: entity.InviteCollaboratorInput{
				CollectionID: collectionID,
				InviterID:    inviterID,
				InviteeID:    inviteeID,
			},
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: inviterID,
					}, nil).
					Once()

				repo.On("CheckIfCollaborator", mock.Anything, collectionID, inviteeID).
					Return(false, nil).
					Once()

				repo.On("CountCollaborators", mock.Anything, collectionID).
					Return(validCollaboratorsCount, nil).
					Once()

				repo.On("GetInvitationByInvitee", mock.Anything, collectionID, inviteeID).
					Return(entity.Invitation{}, errRepo).
					Once()
			},
			wantErr: errRepo,
		},
		{
			name: "error - create invitation repository fails",
			input: entity.InviteCollaboratorInput{
				CollectionID: collectionID,
				InviterID:    inviterID,
				InviteeID:    inviteeID,
			},
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: inviterID,
					}, nil).
					Once()

				repo.On("CheckIfCollaborator", mock.Anything, collectionID, inviteeID).
					Return(false, nil).
					Once()

				repo.On("CountCollaborators", mock.Anything, collectionID).
					Return(validCollaboratorsCount, nil).
					Once()

				repo.On("GetInvitationByInvitee", mock.Anything, collectionID, inviteeID).
					Return(entity.Invitation{}, apperror.InvitationNotFoundForInvitee(collectionID, inviteeID)).
					Once()

				repo.On("CreateInvitation", mock.Anything, mock.MatchedBy(func(in entity.InviteCollaboratorInput) bool {
					return in.CollectionID == collectionID &&
						in.InviterID == inviterID &&
						in.InviteeID == inviteeID &&
						in.Expiry.After(time.Now())
				})).
					Return("", errRepo).
					Once()
			},
			wantErr: errRepo,
		},
		{
			name: "error - invitation cooldown",
			input: entity.InviteCollaboratorInput{
				CollectionID: collectionID,
				InviterID:    inviterID,
				InviteeID:    inviteeID,
			},
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: inviterID,
					}, nil).
					Once()

				repo.On("CheckIfCollaborator", mock.Anything, collectionID, inviteeID).
					Return(false, nil).
					Once()

				repo.On("CountCollaborators", mock.Anything, collectionID).
					Return(validCollaboratorsCount, nil).
					Once()

				repo.On("GetInvitationByInvitee", mock.Anything, collectionID, inviteeID).
					Return(entity.Invitation{
						ID:         invitationID,
						LastSentAt: time.Now(), // still cooldown
					}, nil).
					Once()
			},
			wantErr: apperror.InvitationCooldown(resendInvitationCooldown),
		},
		{
			name: "error - refresh invitation repository fails",
			input: entity.InviteCollaboratorInput{
				CollectionID: collectionID,
				InviterID:    inviterID,
				InviteeID:    inviteeID,
			},
			mockFn: func(repo *mockCollectionRepository, tx *mockTxManager) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				repo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: inviterID,
					}, nil).
					Once()

				repo.On("CheckIfCollaborator", mock.Anything, collectionID, inviteeID).
					Return(false, nil).
					Once()

				repo.On("CountCollaborators", mock.Anything, collectionID).
					Return(validCollaboratorsCount, nil).
					Once()

				repo.On("GetInvitationByInvitee", mock.Anything, collectionID, inviteeID).
					Return(entity.Invitation{
						ID:         invitationID,
						LastSentAt: time.Now().Add(-resendInvitationCooldown),
					}, nil).
					Once()

				repo.On("RefreshInvitation", mock.Anything, invitationID, mock.MatchedBy(func(t time.Time) bool {
					return t.After(time.Now())
				})).
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

			err := svc.InviteCollaborator(context.Background(), tt.input)
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
