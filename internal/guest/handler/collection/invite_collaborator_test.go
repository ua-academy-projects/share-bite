package collection

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

func TestInviteCollaborator(t *testing.T) {
	t.Parallel()

	var (
		collectionID = gofakeit.UUID()
		inviterID    = gofakeit.UUID()

		inviteeID = gofakeit.UUID()

		serviceErr = errors.New("unexpected service error")
	)

	tests := []struct {
		name string

		body         any
		collectionID string
		inviterID    any

		mockFn func(s *mockCollectionService)

		wantCode int
		wantBody any
	}{
		{
			name: "success",
			body: inviteCollaboratorRequest{
				CustomerID: inviteeID,
			},
			collectionID: collectionID,
			inviterID:    inviterID,
			mockFn: func(s *mockCollectionService) {
				s.On("InviteCollaborator", mock.Anything, entity.InviteCollaboratorInput{
					CollectionID: collectionID,
					InviterID:    inviterID,
					InviteeID:    inviteeID,
				}).
					Return(nil).
					Once()
			},
			wantCode: http.StatusNoContent,
			wantBody: nil,
		},
		{
			name:         "invalid json",
			body:         "{broken-json",
			collectionID: collectionID,
			inviterID:    inviterID,
			mockFn:       func(s *mockCollectionService) {},
			wantCode:     http.StatusBadRequest,
			wantBody:     response.ErrorResponse{Message: apperror.ErrInvalidJSON.Error()},
		},
		{
			name:         "validation error - customer (invitee) id is required",
			body:         inviteCollaboratorRequest{},
			collectionID: collectionID,
			inviterID:    inviterID,
			mockFn:       func(s *mockCollectionService) {},
			wantCode:     http.StatusBadRequest,
			wantBody: response.ErrorResponse{
				Message: validationMsg,
				Details: []response.ErrorDetail{
					{Field: "customerId", Message: "This field is required"},
				},
			},
		},
		{
			name: "validation error - invalid customer (invitee) id",
			body: inviteCollaboratorRequest{
				CustomerID: "invalid-invitee-customer-id",
			},
			collectionID: collectionID,
			inviterID:    inviterID,
			mockFn:       func(s *mockCollectionService) {},
			wantCode:     http.StatusBadRequest,
			wantBody: response.ErrorResponse{
				Message: validationMsg,
				Details: []response.ErrorDetail{
					{Field: "customerId", Message: "This field must be a valid UUID"},
				},
			},
		},
		{
			name: "validation error - invalid collection id",
			body: inviteCollaboratorRequest{
				CustomerID: inviteeID,
			},
			collectionID: "invalid-collection-id",
			inviterID:    inviterID,
			mockFn:       func(s *mockCollectionService) {},
			wantCode:     http.StatusBadRequest,
			wantBody: response.ErrorResponse{
				Message: validationMsg,
				Details: []response.ErrorDetail{
					{Field: "collectionId", Message: "This field must be a valid UUID"},
				},
			},
		},
		{
			name: "missing customer (inviter) id in ctx",
			body: inviteCollaboratorRequest{
				CustomerID: inviteeID,
			},
			collectionID: collectionID,
			mockFn:       func(s *mockCollectionService) {},
			wantCode:     http.StatusInternalServerError,
			wantBody:     response.ErrorResponse{Message: internalErrMsg},
		},
		{
			name: "service unknown error",
			body: inviteCollaboratorRequest{
				CustomerID: inviteeID,
			},
			collectionID: collectionID,
			inviterID:    inviterID,
			mockFn: func(s *mockCollectionService) {
				s.On("InviteCollaborator", mock.Anything, entity.InviteCollaboratorInput{
					CollectionID: collectionID,
					InviterID:    inviterID,
					InviteeID:    inviteeID,
				}).
					Return(serviceErr).
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
			tt.mockFn(svc)

			h := &handler{service: svc}
			r := newTestRouter()

			r.POST("/collections/:collectionId/invitations", withCustomerID(tt.inviterID, h.inviteCollaborator))

			var w *httptest.ResponseRecorder
			if s, ok := tt.body.(string); ok {
				w = performRawJSONRequest(t, r, http.MethodPost, "/collections/"+tt.collectionID+"/invitations", s)
			} else {
				w = performJSONRequest(t, r, http.MethodPost, "/collections/"+tt.collectionID+"/invitations", tt.body)
			}

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
