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
	"github.com/ua-academy-projects/share-bite/internal/guest/error/code"
	_ "github.com/ua-academy-projects/share-bite/pkg/notification"
)

func TestInviteCollaborator(t *testing.T) {
	t.Parallel()

	var (
		collectionID = "random-collection-uuid"
		inviterID    = "random-inviter-customer-uuid"
		inviteeID    = "random-invitee-customer-uuid"
		invitationID = "invitation-uuid"
		targetUserID = "target-user-uuid"

		validCollaboratorsCount = 4

		errRepo = errors.New("unexpected repository error")
	)

	tests := []struct {
		name  string
		input entity.InviteCollaboratorInput

		mockFn func(collectionRepo *mockCollectionRepository, customerRepo *mockCustomerRepository, tx *mockTxManager, pub *mockPublisher)

		wait           bool
		wantErr        error
		wantErrCode    code.Code
		wantErrContain string
	}{
		{
			name: "success - new invitation",
			input: entity.InviteCollaboratorInput{
				CollectionID: collectionID,
				InviterID:    inviterID,
				InviteeID:    inviteeID,
			},
			mockFn: func(collectionRepo *mockCollectionRepository, customerRepo *mockCustomerRepository, tx *mockTxManager, pub *mockPublisher) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				collectionRepo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(
						entity.Collection{
							ID:         collectionID,
							CustomerID: inviterID,
						},
						nil,
					).
					Once()

				collectionRepo.On("CheckIfCollaborator", mock.Anything, collectionID, inviteeID).
					Return(false, nil).
					Once()

				collectionRepo.On("CountCollaborators", mock.Anything, collectionID).
					Return(validCollaboratorsCount, nil).
					Once()

				collectionRepo.On("GetInvitationByInvitee", mock.Anything, collectionID, inviteeID).
					Return(
						entity.Invitation{},
						apperror.InvitationNotFoundForInvitee(collectionID, inviteeID),
					).
					Once()

				collectionRepo.On("CreateInvitation", mock.Anything, mock.MatchedBy(func(in entity.InviteCollaboratorInput) bool {
					return in.CollectionID == collectionID &&
						in.InviterID == inviterID &&
						in.InviteeID == inviteeID &&
						in.Expiry.After(time.Now())
				})).
					Return("new-invitation-uuid", nil).
					Once()

				customerRepo.On("GetByID", mock.Anything, inviteeID).
					Return(entity.Customer{ID: inviteeID, UserID: targetUserID}, nil).
					Once()

				pub.On("Publish", mock.Anything, targetUserID, mock.AnythingOfType("notification.Message")).
					Return(nil).
					Once()
			},
			wait:    true,
			wantErr: nil,
		},
		{
			name: "success - refresh invitation",
			input: entity.InviteCollaboratorInput{
				CollectionID: collectionID,
				InviterID:    inviterID,
				InviteeID:    inviteeID,
			},
			mockFn: func(collectionRepo *mockCollectionRepository, customerRepo *mockCustomerRepository, tx *mockTxManager, pub *mockPublisher) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				collectionRepo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(
						entity.Collection{
							ID:         collectionID,
							CustomerID: inviterID,
						},
						nil,
					).
					Once()

				collectionRepo.On("CheckIfCollaborator", mock.Anything, collectionID, inviteeID).
					Return(false, nil).
					Once()

				collectionRepo.On("CountCollaborators", mock.Anything, collectionID).
					Return(validCollaboratorsCount, nil).
					Once()

				collectionRepo.On("GetInvitationByInvitee", mock.Anything, collectionID, inviteeID).
					Return(
						entity.Invitation{
							ID:         invitationID,
							LastSentAt: time.Now().Add(-resendInvitationCooldown),
						},
						nil,
					).
					Once()

				collectionRepo.On("RefreshInvitation", mock.Anything, invitationID, mock.MatchedBy(func(t time.Time) bool {
					return t.After(time.Now())
				})).
					Return(nil).
					Once()

				customerRepo.On("GetByID", mock.Anything, inviteeID).
					Return(entity.Customer{ID: inviteeID, UserID: targetUserID}, nil).
					Once()

				pub.On("Publish", mock.Anything, targetUserID, mock.AnythingOfType("notification.Message")).
					Return(nil).
					Once()
			},
			wait:    true,
			wantErr: nil,
		},
		{
			name: "error - get collection for update repository fails",
			input: entity.InviteCollaboratorInput{
				CollectionID: collectionID,
				InviterID:    inviterID,
				InviteeID:    inviteeID,
			},
			mockFn: func(collectionRepo *mockCollectionRepository, customerRepo *mockCustomerRepository, tx *mockTxManager, pub *mockPublisher) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				collectionRepo.On("GetCollectionForUpdate", mock.Anything, collectionID).
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
			mockFn: func(collectionRepo *mockCollectionRepository, customerRepo *mockCustomerRepository, tx *mockTxManager, pub *mockPublisher) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				collectionRepo.On("GetCollectionForUpdate", mock.Anything, collectionID).
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
			mockFn: func(collectionRepo *mockCollectionRepository, customerRepo *mockCustomerRepository, tx *mockTxManager, pub *mockPublisher) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				collectionRepo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: inviterID,
					}, nil).
					Once()

				collectionRepo.On("CheckIfCollaborator", mock.Anything, collectionID, "not-owner-but-collaborator-uuid").
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
			mockFn: func(collectionRepo *mockCollectionRepository, customerRepo *mockCustomerRepository, tx *mockTxManager, pub *mockPublisher) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				collectionRepo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: inviterID,
					}, nil).
					Once()

				collectionRepo.On("CheckIfCollaborator", mock.Anything, collectionID, "outsider-uuid").
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
			mockFn: func(collectionRepo *mockCollectionRepository, customerRepo *mockCustomerRepository, tx *mockTxManager, pub *mockPublisher) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				collectionRepo.On("GetCollectionForUpdate", mock.Anything, collectionID).
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
			mockFn: func(collectionRepo *mockCollectionRepository, customerRepo *mockCustomerRepository, tx *mockTxManager, pub *mockPublisher) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				collectionRepo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: inviterID,
					}, nil).
					Once()

				collectionRepo.On("CheckIfCollaborator", mock.Anything, collectionID, inviteeID).
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
			mockFn: func(collectionRepo *mockCollectionRepository, customerRepo *mockCustomerRepository, tx *mockTxManager, pub *mockPublisher) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				collectionRepo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: inviterID,
					}, nil).
					Once()

				collectionRepo.On("CheckIfCollaborator", mock.Anything, collectionID, inviteeID).
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
			mockFn: func(collectionRepo *mockCollectionRepository, customerRepo *mockCustomerRepository, tx *mockTxManager, pub *mockPublisher) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				collectionRepo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: inviterID,
					}, nil).
					Once()

				collectionRepo.On("CheckIfCollaborator", mock.Anything, collectionID, inviteeID).
					Return(false, nil).
					Once()

				collectionRepo.On("CountCollaborators", mock.Anything, collectionID).
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
			mockFn: func(collectionRepo *mockCollectionRepository, customerRepo *mockCustomerRepository, tx *mockTxManager, pub *mockPublisher) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				collectionRepo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: inviterID,
					}, nil).
					Once()

				collectionRepo.On("CheckIfCollaborator", mock.Anything, collectionID, inviteeID).
					Return(false, nil).
					Once()

				collectionRepo.On("CountCollaborators", mock.Anything, collectionID).
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
			mockFn: func(collectionRepo *mockCollectionRepository, customerRepo *mockCustomerRepository, tx *mockTxManager, pub *mockPublisher) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				collectionRepo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: inviterID,
					}, nil).
					Once()

				collectionRepo.On("CheckIfCollaborator", mock.Anything, collectionID, inviteeID).
					Return(false, nil).
					Once()

				collectionRepo.On("CountCollaborators", mock.Anything, collectionID).
					Return(validCollaboratorsCount, nil).
					Once()

				collectionRepo.On("GetInvitationByInvitee", mock.Anything, collectionID, inviteeID).
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
			mockFn: func(collectionRepo *mockCollectionRepository, customerRepo *mockCustomerRepository, tx *mockTxManager, pub *mockPublisher) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				collectionRepo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: inviterID,
					}, nil).
					Once()

				collectionRepo.On("CheckIfCollaborator", mock.Anything, collectionID, inviteeID).
					Return(false, nil).
					Once()

				collectionRepo.On("CountCollaborators", mock.Anything, collectionID).
					Return(validCollaboratorsCount, nil).
					Once()

				collectionRepo.On("GetInvitationByInvitee", mock.Anything, collectionID, inviteeID).
					Return(entity.Invitation{}, apperror.InvitationNotFoundForInvitee(collectionID, inviteeID)).
					Once()

				collectionRepo.On("CreateInvitation", mock.Anything, mock.MatchedBy(func(in entity.InviteCollaboratorInput) bool {
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
			mockFn: func(collectionRepo *mockCollectionRepository, customerRepo *mockCustomerRepository, tx *mockTxManager, pub *mockPublisher) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				collectionRepo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: inviterID,
					}, nil).
					Once()

				collectionRepo.On("CheckIfCollaborator", mock.Anything, collectionID, inviteeID).
					Return(false, nil).
					Once()

				collectionRepo.On("CountCollaborators", mock.Anything, collectionID).
					Return(validCollaboratorsCount, nil).
					Once()

				collectionRepo.On("GetInvitationByInvitee", mock.Anything, collectionID, inviteeID).
					Return(entity.Invitation{
						ID:         invitationID,
						LastSentAt: time.Now(), // still cooldown
					}, nil).
					Once()
			},
			wantErrCode:    code.TooManyRequests,
			wantErrContain: "before resending this invitation",
		},
		{
			name: "error - refresh invitation repository fails",
			input: entity.InviteCollaboratorInput{
				CollectionID: collectionID,
				InviterID:    inviterID,
				InviteeID:    inviteeID,
			},
			mockFn: func(collectionRepo *mockCollectionRepository, customerRepo *mockCustomerRepository, tx *mockTxManager, pub *mockPublisher) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				collectionRepo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: inviterID,
					}, nil).
					Once()

				collectionRepo.On("CheckIfCollaborator", mock.Anything, collectionID, inviteeID).
					Return(false, nil).
					Once()

				collectionRepo.On("CountCollaborators", mock.Anything, collectionID).
					Return(validCollaboratorsCount, nil).
					Once()

				collectionRepo.On("GetInvitationByInvitee", mock.Anything, collectionID, inviteeID).
					Return(entity.Invitation{
						ID:         invitationID,
						LastSentAt: time.Now().Add(-resendInvitationCooldown),
					}, nil).
					Once()

				collectionRepo.On("RefreshInvitation", mock.Anything, invitationID, mock.MatchedBy(func(t time.Time) bool {
					return t.After(time.Now())
				})).
					Return(errRepo).
					Once()
			},
			wantErr: errRepo,
		},
		{
			name: "success main flow - get invitee customer fails in goroutine",
			input: entity.InviteCollaboratorInput{
				CollectionID: collectionID,
				InviterID:    inviterID,
				InviteeID:    inviteeID,
			},
			mockFn: func(collectionRepo *mockCollectionRepository, customerRepo *mockCustomerRepository, tx *mockTxManager, pub *mockPublisher) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				collectionRepo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: inviterID}, nil).Once()

				collectionRepo.On("CheckIfCollaborator", mock.Anything, collectionID, inviteeID).
					Return(false, nil).Once()

				collectionRepo.On("CountCollaborators", mock.Anything, collectionID).
					Return(validCollaboratorsCount, nil).Once()

				collectionRepo.On("GetInvitationByInvitee", mock.Anything, collectionID, inviteeID).
					Return(entity.Invitation{}, apperror.InvitationNotFoundForInvitee(collectionID, inviteeID)).Once()

				collectionRepo.On("CreateInvitation", mock.Anything, mock.Anything).
					Return(invitationID, nil).Once()

				customerRepo.On("GetByID", mock.Anything, inviteeID).
					Return(entity.Customer{}, errRepo).Once()
			},
			wait:    true,
			wantErr: nil,
		},
		{
			name: "success main flow - publish fails in goroutine",
			input: entity.InviteCollaboratorInput{
				CollectionID: collectionID,
				InviterID:    inviterID,
				InviteeID:    inviteeID,
			},
			mockFn: func(collectionRepo *mockCollectionRepository, customerRepo *mockCustomerRepository, tx *mockTxManager, pub *mockPublisher) {
				tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()

				collectionRepo.On("GetCollectionForUpdate", mock.Anything, collectionID).
					Return(entity.Collection{ID: collectionID, CustomerID: inviterID}, nil).Once()

				collectionRepo.On("CheckIfCollaborator", mock.Anything, collectionID, inviteeID).
					Return(false, nil).Once()

				collectionRepo.On("CountCollaborators", mock.Anything, collectionID).
					Return(validCollaboratorsCount, nil).Once()

				collectionRepo.On("GetInvitationByInvitee", mock.Anything, collectionID, inviteeID).
					Return(entity.Invitation{}, apperror.InvitationNotFoundForInvitee(collectionID, inviteeID)).Once()

				collectionRepo.On("CreateInvitation", mock.Anything, mock.Anything).
					Return(invitationID, nil).Once()

				customerRepo.On("GetByID", mock.Anything, inviteeID).
					Return(entity.Customer{ID: inviteeID, UserID: targetUserID}, nil).Once()

				pub.On("Publish", mock.Anything, targetUserID, mock.AnythingOfType("notification.Message")).
					Return(errors.New("redis timeout")).Once()
			},
			wait:    true,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			collectionRepo := new(mockCollectionRepository)
			txManager := new(mockTxManager)
			customerRepo := new(mockCustomerRepository)
			pub := new(mockPublisher)

			svc := New(collectionRepo, customerRepo, txManager, nil, WithPublisher(pub))
			tt.mockFn(collectionRepo, customerRepo, txManager, pub)

			err := svc.InviteCollaborator(context.Background(), tt.input)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else if tt.wantErrCode != "" {
				require.Error(t, err)

				var appErr *apperror.Error
				require.ErrorAs(t, err, &appErr)
				assert.Equal(t, tt.wantErrCode, appErr.Code)
				assert.Contains(t, appErr.Error(), tt.wantErrContain)
			} else {
				require.NoError(t, err)
			}

			if tt.wait {
				require.Eventually(t, func() bool {
					if len(pub.ExpectedCalls) > 0 {
						for _, call := range pub.Calls {
							if call.Method == "Publish" {
								return true
							}
						}

						return false
					}

					for _, call := range customerRepo.Calls {
						if call.Method == "GetByID" {
							return true
						}
					}

					return false
				}, time.Second, 10*time.Millisecond, "async goroutine was not called in time")
			}

			collectionRepo.AssertExpectations(t)
			customerRepo.AssertExpectations(t)
			txManager.AssertExpectations(t)
			pub.AssertExpectations(t)
		})
	}
}
