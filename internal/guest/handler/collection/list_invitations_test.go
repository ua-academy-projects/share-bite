package collection

import (
	"encoding/base64"
	"errors"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

func TestListInvitations(t *testing.T) {
	t.Parallel()

	var (
		customerID   = gofakeit.UUID()
		collectionID = gofakeit.UUID()
		inviterID    = gofakeit.UUID()
		inviteeID    = gofakeit.UUID()

		invitationId1 = gofakeit.UUID()
		invitationId2 = gofakeit.UUID()
		invitationId3 = gofakeit.UUID()

		validCursorID = gofakeit.UUID()
		validToken    = base64.RawURLEncoding.EncodeToString([]byte(validCursorID))

		statusPending = string(entity.PendingInvitationStatus)
		limitValid    = 10
		limitMax      = 100

		invitations = []entity.EnrichedInvitation{
			{ID: invitationId1},
			{ID: invitationId2},
			{ID: invitationId3},
		}

		serviceErr = errors.New("unexpected service error")
	)

	tests := []struct {
		name string

		query      string
		customerID any
		mockFn     func(s *mockCollectionService)

		wantCode int
		wantBody any
	}{
		{
			name:       "success - all possible params",
			query:      "?collectionId=" + collectionID + "&status=" + statusPending + "&pageSize=10&pageToken=" + validToken,
			customerID: customerID,
			mockFn: func(s *mockCollectionService) {
				status := entity.InvitationStatus(statusPending)
				s.On("ListInvitations", mock.Anything, entity.ListInvitationsInput{
					CollectionID: &collectionID,
					Status:       &status,
					CursorID:     validCursorID,
					CallerID:     customerID,
					Limit:        limitValid + 1,
				}).
					Return(
						entity.ListInvitationsOutput{
							Invitations:  invitations,
							NextCursorID: invitationId3,
						},
						nil,
					).
					Once()
			},
			wantCode: http.StatusOK,
			wantBody: listInvitationsResponse{
				Invitations: []invitationResponse{
					{ID: invitationId1},
					{ID: invitationId2},
					{ID: invitationId3},
				},
				NextPageToken: base64.RawURLEncoding.EncodeToString([]byte(invitationId3)),
			},
		},
		{
			name:       "success - pageSize is exactly at max limit (100)",
			query:      "?inviterId=" + customerID + "&pageSize=100",
			customerID: customerID,
			mockFn: func(s *mockCollectionService) {
				s.On("ListInvitations", mock.Anything, entity.ListInvitationsInput{
					InviterID: &customerID,
					CallerID:  customerID,
					Limit:     limitMax + 1,
				}).
					Return(
						entity.ListInvitationsOutput{
							Invitations:  []entity.EnrichedInvitation{{ID: invitationId1}},
							NextCursorID: "",
						},
						nil,
					).
					Once()
			},
			wantCode: http.StatusOK,
			wantBody: listInvitationsResponse{
				Invitations:   []invitationResponse{{ID: invitationId1}},
				NextPageToken: "",
			},
		},
		{
			name:       "success - pageSize is 0 (fallback to default limit)",
			query:      "?inviterId=" + customerID + "&pageSize=0",
			customerID: customerID,
			mockFn: func(s *mockCollectionService) {
				s.On("ListInvitations", mock.Anything, entity.ListInvitationsInput{
					InviterID: &customerID,
					CallerID:  customerID,
					Limit:     defaultInvitationsLimit + 1,
				}).
					Return(
						entity.ListInvitationsOutput{
							Invitations:  []entity.EnrichedInvitation{{ID: invitationId1}},
							NextCursorID: "",
						},
						nil,
					).
					Once()
			},
			wantCode: http.StatusOK,
			wantBody: listInvitationsResponse{
				Invitations:   []invitationResponse{{ID: invitationId1}},
				NextPageToken: "",
			},
		},
		{
			name:       "success - by inviter id",
			query:      "?inviterId=" + customerID,
			customerID: customerID,
			mockFn: func(s *mockCollectionService) {
				s.On("ListInvitations", mock.Anything, entity.ListInvitationsInput{
					InviterID: &customerID,
					CallerID:  customerID,
					Limit:     defaultInvitationsLimit + 1,
				}).
					Return(
						entity.ListInvitationsOutput{
							Invitations:  []entity.EnrichedInvitation{{ID: invitationId1}},
							NextCursorID: "",
						},
						nil,
					).
					Once()
			},
			wantCode: http.StatusOK,
			wantBody: listInvitationsResponse{
				Invitations:   []invitationResponse{{ID: invitationId1}},
				NextPageToken: "",
			},
		},
		{
			name:       "success - by invitee id",
			query:      "?inviteeId=" + customerID,
			customerID: customerID,
			mockFn: func(s *mockCollectionService) {
				s.On("ListInvitations", mock.Anything, entity.ListInvitationsInput{
					InviteeID: &customerID,
					CallerID:  customerID,
					Limit:     defaultInvitationsLimit + 1,
				}).
					Return(
						entity.ListInvitationsOutput{
							Invitations:  []entity.EnrichedInvitation{{ID: invitationId2}},
							NextCursorID: "",
						},
						nil,
					).
					Once()
			},
			wantCode: http.StatusOK,
			wantBody: listInvitationsResponse{
				Invitations:   []invitationResponse{{ID: invitationId2}},
				NextPageToken: "",
			},
		},
		{
			name:       "validation error - missing all required ids",
			query:      "?status=" + statusPending,
			customerID: customerID,
			mockFn:     func(s *mockCollectionService) {},
			wantCode:   http.StatusBadRequest,
			wantBody: response.ErrorResponse{
				Message: validationMsg,
				Details: []response.ErrorDetail{
					{Field: "collectionId", Message: "This field is required if all of InviterID, InviteeID are missing"},
					{Field: "inviterId", Message: "This field is required if all of InviteeID, CollectionID are missing"},
					{Field: "inviteeId", Message: "This field is required if all of InviterID, CollectionID are missing"},
				},
			},
		},
		{
			name:       "validation error - invalid collectionId format",
			query:      "?collectionId=invalid-uuid",
			customerID: customerID,
			mockFn:     func(s *mockCollectionService) {},
			wantCode:   http.StatusBadRequest,
			wantBody: response.ErrorResponse{
				Message: validationMsg,
				Details: []response.ErrorDetail{
					{Field: "collectionId", Message: "This field must be a valid UUID"},
				},
			},
		},
		{
			name:       "validation error - invalid inviterId format",
			query:      "?inviterId=invalid-uuid",
			customerID: customerID,
			mockFn:     func(s *mockCollectionService) {},
			wantCode:   http.StatusBadRequest,
			wantBody: response.ErrorResponse{
				Message: validationMsg,
				Details: []response.ErrorDetail{
					{Field: "inviterId", Message: "This field must be a valid UUID"},
				},
			},
		},
		{
			name:       "validation error - invalid inviteeId format",
			query:      "?inviteeId=invalid-uuid",
			customerID: customerID,
			mockFn:     func(s *mockCollectionService) {},
			wantCode:   http.StatusBadRequest,
			wantBody: response.ErrorResponse{
				Message: validationMsg,
				Details: []response.ErrorDetail{
					{Field: "inviteeId", Message: "This field must be a valid UUID"},
				},
			},
		},
		{
			name:       "validation error - invalid status",
			query:      "?collectionId=" + collectionID + "&status=unknown",
			customerID: customerID,
			mockFn:     func(s *mockCollectionService) {},
			wantCode:   http.StatusBadRequest,
			wantBody: response.ErrorResponse{
				Message: validationMsg,
				Details: []response.ErrorDetail{
					{Field: "status", Message: "This field is invalid"},
				},
			},
		},
		{
			name:       "validation error - pageSize exceeds limit",
			query:      "?collectionId=" + collectionID + "&pageSize=101",
			customerID: customerID,
			mockFn:     func(s *mockCollectionService) {},
			wantCode:   http.StatusBadRequest,
			wantBody: response.ErrorResponse{
				Message: validationMsg,
				Details: []response.ErrorDetail{
					{Field: "pageSize", Message: "This field must be less than or equal to 100"},
				},
			},
		},
		{
			name:       "error - invalid pageToken base64",
			query:      "?collectionId=" + collectionID + "&pageToken=***",
			customerID: customerID,
			mockFn:     func(s *mockCollectionService) {},
			wantCode:   http.StatusBadRequest,
			wantBody:   response.ErrorResponse{Message: apperror.ErrInvalidPageToken.Error()},
		},
		{
			name:       "error - listing others outbound",
			query:      "?inviterId=" + inviterID,
			customerID: customerID,
			mockFn:     func(s *mockCollectionService) {},
			wantCode:   http.StatusForbidden,
			wantBody:   response.ErrorResponse{Message: apperror.ErrCannotListOthersOutboundInvites.Error()},
		},
		{
			name:       "error - listing others inbound",
			query:      "?inviteeId=" + inviteeID,
			customerID: customerID,
			mockFn:     func(s *mockCollectionService) {},
			wantCode:   http.StatusForbidden,
			wantBody:   response.ErrorResponse{Message: apperror.ErrCannotListOthersInboundInvites.Error()},
		},
		{
			name:       "missing context customer id",
			query:      "?collectionId=" + collectionID,
			customerID: nil,
			mockFn:     func(s *mockCollectionService) {},
			wantCode:   http.StatusInternalServerError,
			wantBody:   response.ErrorResponse{Message: internalErrMsg},
		},
		{
			name:       "service error",
			query:      "?collectionId=" + collectionID,
			customerID: customerID,
			mockFn: func(s *mockCollectionService) {
				s.On("ListInvitations", mock.Anything, mock.Anything).
					Return(entity.ListInvitationsOutput{}, serviceErr).
					Once()
			},
			wantCode: http.StatusInternalServerError,
			wantBody: response.ErrorResponse{Message: internalErrMsg},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := new(mockCollectionService)
			h := &handler{
				service: svc,
			}
			tt.mockFn(svc)

			r := newTestRouter()
			r.GET("/collections/invitations", withCustomerID(tt.customerID, h.listInvitations))

			w := performRequest(r, http.MethodGet, "/collections/invitations"+tt.query)
			require.Equal(t, tt.wantCode, w.Code)

			if tt.wantBody != nil {
				assertJSONBody(t, tt.wantBody, w.Body.String())
			} else {
				require.Empty(t, w.Body.String())
			}

			svc.AssertExpectations(t)
		})
	}
}
