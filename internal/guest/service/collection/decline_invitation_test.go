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

func TestDeclineInvitation(t *testing.T) {
	t.Parallel()

	var (
		invitationID = "random-invation-uuid"
		inviteeID    = "random-invitee-customer-uuid"

		errRepo = errors.New("unexpected repository error")
	)

	tests := []struct {
		name string

		invitationID string
		inviteeID    string

		mockFn func(repo *mockCollectionRepository)

		wantErr error
	}{
		{
			name:         "success",
			invitationID: invitationID,
			inviteeID:    inviteeID,
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetInvitation", mock.Anything, invitationID).
					Return(
						entity.Invitation{
							ID:        invitationID,
							InviteeID: inviteeID,
							Status:    entity.PendingInvitationStatus,
						},
						nil,
					).
					Once()

				repo.On("UpdateInvitationStatus", mock.Anything, invitationID, entity.DeclinedInvitationStatus).
					Return(nil).
					Once()
			},
			wantErr: nil,
		},
		{
			name:         "error - invitation not found",
			invitationID: invitationID,
			inviteeID:    inviteeID,
			mockFn: func(repo *mockCollectionRepository) {
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
			mockFn: func(repo *mockCollectionRepository) {
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
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetInvitation", mock.Anything, invitationID).
					Return(
						entity.Invitation{
							ID:        invitationID,
							InviteeID: "not-invitee-customer-id",
							Status:    entity.PendingInvitationStatus,
						},
						nil,
					).
					Once()
			},
			wantErr: apperror.InvitationNotFoundID(invitationID),
		},
		{
			name:         "error - invitation already processed (accepted)",
			invitationID: invitationID,
			inviteeID:    inviteeID,
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetInvitation", mock.Anything, invitationID).
					Return(
						entity.Invitation{
							ID:        invitationID,
							InviteeID: inviteeID,
							Status:    entity.AcceptedInvitationStatus,
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
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetInvitation", mock.Anything, invitationID).
					Return(
						entity.Invitation{
							ID:        invitationID,
							InviteeID: inviteeID,
							Status:    entity.DeclinedInvitationStatus,
						},
						nil,
					).
					Once()
			},
			wantErr: apperror.ErrInvitationAlreadyProcessed,
		},
		{
			name:         "error - update invitation status repository fails",
			invitationID: invitationID,
			inviteeID:    inviteeID,
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetInvitation", mock.Anything, invitationID).
					Return(
						entity.Invitation{
							ID:        invitationID,
							InviteeID: inviteeID,
							Status:    entity.PendingInvitationStatus,
						},
						nil,
					).
					Once()

				repo.On("UpdateInvitationStatus", mock.Anything, invitationID, entity.DeclinedInvitationStatus).
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

			err := svc.DeclineInvitation(context.Background(), tt.invitationID, tt.inviteeID)
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
