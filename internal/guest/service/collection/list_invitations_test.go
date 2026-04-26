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

func TestListInvitations(t *testing.T) {
	t.Parallel()

	var (
		collectionID = "random-collection-uuid"
		inviterID    = "random-inviter-customer-uuid"
		inviteeID    = "random-invitee-customer-uuid"
		outsiderID   = "random-outsider-uuid"

		errRepo = errors.New("unexpected repository error")
	)

	pendingStatus := entity.InvitationStatus(entity.PendingInvitationStatus)

	tests := []struct {
		name  string
		input entity.ListInvitationsInput

		mockFn func(repo *mockCollectionRepository)

		wantErr  error
		wantResp entity.ListInvitationsOutput
	}{
		{
			name: "success - with next page",
			input: entity.ListInvitationsInput{
				CollectionID: &collectionID,
				CallerID:     inviterID,
				Limit:        3,
			},
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: inviterID,
					}, nil).
					Once()

				repo.On("ListInvitations", mock.Anything, entity.ListInvitationsInput{
					CollectionID: &collectionID,
					CallerID:     inviterID,
					Limit:        3,
				}).
					Return([]entity.EnrichedInvitation{
						{ID: "invitation-uuid-1"},
						{ID: "invitation-uuid-2"},
						{ID: "invitation-uuid-3"},
					}, nil).
					Once()
			},
			wantErr: nil,
			wantResp: entity.ListInvitationsOutput{
				Invitations: []entity.EnrichedInvitation{
					{ID: "invitation-uuid-1"},
					{ID: "invitation-uuid-2"},
				},
				NextCursorID: "invitation-uuid-2",
			},
		},
		{
			name: "success - last page (no next cursor)",
			input: entity.ListInvitationsInput{
				CollectionID: &collectionID,
				CallerID:     inviterID,
				Limit:        3,
			},
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: inviterID,
					}, nil).
					Once()

				repo.On("ListInvitations", mock.Anything, entity.ListInvitationsInput{
					CollectionID: &collectionID,
					CallerID:     inviterID,
					Limit:        3,
				}).
					Return([]entity.EnrichedInvitation{
						{ID: "invitation-uuid-1"},
					}, nil).
					Once()
			},
			wantErr: nil,
			wantResp: entity.ListInvitationsOutput{
				Invitations: []entity.EnrichedInvitation{
					{ID: "invitation-uuid-1"},
				},
				NextCursorID: "",
			},
		},
		{
			name: "success - with status filter",
			input: entity.ListInvitationsInput{
				CollectionID: &collectionID,
				Status:       &pendingStatus,
				CallerID:     inviterID,
				Limit:        3,
			},
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: inviterID,
					}, nil).
					Once()

				repo.On("ListInvitations", mock.Anything, entity.ListInvitationsInput{
					CollectionID: &collectionID,
					Status:       &pendingStatus,
					CallerID:     inviterID,
					Limit:        3,
				}).
					Return([]entity.EnrichedInvitation{
						{ID: "invitation-uuid-1"},
					}, nil).
					Once()
			},
			wantErr: nil,
			wantResp: entity.ListInvitationsOutput{
				Invitations: []entity.EnrichedInvitation{
					{ID: "invitation-uuid-1"},
				},
				NextCursorID: "",
			},
		},
		{
			name: "success - filter by invitee (no collection check)",
			input: entity.ListInvitationsInput{
				InviteeID: &inviteeID,
				CallerID:  inviteeID,
				Limit:     3,
			},
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("ListInvitations", mock.Anything, entity.ListInvitationsInput{
					InviteeID: &inviteeID,
					CallerID:  inviteeID,
					Limit:     3,
				}).
					Return([]entity.EnrichedInvitation{
						{ID: "invitation-uuid-1"},
					}, nil).
					Once()
			},
			wantErr: nil,
			wantResp: entity.ListInvitationsOutput{
				Invitations: []entity.EnrichedInvitation{
					{ID: "invitation-uuid-1"},
				},
				NextCursorID: "",
			},
		},
		{
			name: "error - collection not found",
			input: entity.ListInvitationsInput{
				CollectionID: &collectionID,
				CallerID:     inviterID,
				Limit:        3,
			},
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{}, apperror.CollectionNotFoundID(collectionID)).
					Once()
			},
			wantErr:  apperror.CollectionNotFoundID(collectionID),
			wantResp: entity.ListInvitationsOutput{},
		},
		{
			name: "error - get collection repository fails",
			input: entity.ListInvitationsInput{
				CollectionID: &collectionID,
				CallerID:     inviterID,
				Limit:        3,
			},
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{}, errRepo).
					Once()
			},
			wantErr:  errRepo,
			wantResp: entity.ListInvitationsOutput{},
		},
		{
			name: "error - caller is collaborator (not owner)",
			input: entity.ListInvitationsInput{
				CollectionID: &collectionID,
				CallerID:     "collaborator-uuid",
				Limit:        3,
			},
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: inviterID,
					}, nil).
					Once()

				repo.On("CheckIfCollaborator", mock.Anything, collectionID, "collaborator-uuid").
					Return(true, nil).
					Once()
			},
			wantErr:  apperror.ErrCollectionAccessDenied,
			wantResp: entity.ListInvitationsOutput{},
		},
		{
			name: "error - caller is outsider",
			input: entity.ListInvitationsInput{
				CollectionID: &collectionID,
				CallerID:     outsiderID,
				Limit:        3,
			},
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: inviterID,
					}, nil).
					Once()

				repo.On("CheckIfCollaborator", mock.Anything, collectionID, outsiderID).
					Return(false, nil).
					Once()
			},
			wantErr:  apperror.CollectionNotFoundID(collectionID),
			wantResp: entity.ListInvitationsOutput{},
		},
		{
			name: "error - check if collaborator repository fails",
			input: entity.ListInvitationsInput{
				CollectionID: &collectionID,
				CallerID:     outsiderID,
				Limit:        3,
			},
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: inviterID,
					}, nil).
					Once()

				repo.On("CheckIfCollaborator", mock.Anything, collectionID, outsiderID).
					Return(false, errRepo).
					Once()
			},
			wantErr:  errRepo,
			wantResp: entity.ListInvitationsOutput{},
		},
		{
			name: "error - list invitations repository fails (with collection filter)",
			input: entity.ListInvitationsInput{
				CollectionID: &collectionID,
				CallerID:     inviterID,
				Limit:        3,
			},
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("GetCollection", mock.Anything, collectionID).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: inviterID,
					}, nil).
					Once()

				repo.On("ListInvitations", mock.Anything, entity.ListInvitationsInput{
					CollectionID: &collectionID,
					CallerID:     inviterID,
					Limit:        3,
				}).
					Return([]entity.EnrichedInvitation(nil), errRepo).
					Once()
			},
			wantErr:  errRepo,
			wantResp: entity.ListInvitationsOutput{},
		},
		{
			name: "error - list invitations repository fails (without collection filter)",
			input: entity.ListInvitationsInput{
				CallerID: inviterID,
				Limit:    3,
			},
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("ListInvitations", mock.Anything, entity.ListInvitationsInput{
					CallerID: inviterID,
					Limit:    3,
				}).
					Return([]entity.EnrichedInvitation(nil), errRepo).
					Once()
			},
			wantErr:  errRepo,
			wantResp: entity.ListInvitationsOutput{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := new(mockCollectionRepository)
			svc := New(repo, nil, nil)
			tt.mockFn(repo)

			resp, err := svc.ListInvitations(context.Background(), tt.input)
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
