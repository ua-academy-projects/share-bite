package collection

import (
	"errors"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

func TestDeclineInvitation(t *testing.T) {
	t.Parallel()

	var (
		customerID   = gofakeit.UUID()
		invitationID = gofakeit.UUID()

		serviceErr = errors.New("unexpected service error")
	)

	tests := []struct {
		name string

		invitationID string
		customerID   any

		mockFn func(s *mockCollectionService)

		wantCode int
		wantBody any
	}{
		{
			name:         "success",
			invitationID: invitationID,
			customerID:   customerID,
			mockFn: func(s *mockCollectionService) {
				s.On("DeclineInvitation", mock.Anything, invitationID, customerID).
					Return(nil).
					Once()
			},
			wantCode: http.StatusNoContent,
			wantBody: nil,
		},
		{
			name:         "validation error - invalid invitation id",
			invitationID: "invalid-invitation-id",
			customerID:   customerID,
			mockFn:       func(s *mockCollectionService) {},
			wantCode:     http.StatusBadRequest,
			wantBody: response.ErrorResponse{
				Message: validationMsg,
				Details: []response.ErrorDetail{
					{Field: "invitationId", Message: "This field must be a valid UUID"},
				},
			},
		},
		{
			name:         "missing customer id in ctx",
			invitationID: invitationID,
			customerID:   nil,
			mockFn:       func(s *mockCollectionService) {},
			wantCode:     http.StatusInternalServerError,
			wantBody:     response.ErrorResponse{Message: internalErrMsg},
		},
		{
			name:         "service unknown error",
			invitationID: invitationID,
			customerID:   customerID,
			mockFn: func(s *mockCollectionService) {
				s.On("DeclineInvitation", mock.Anything, invitationID, customerID).
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
			h := &handler{
				service: svc,
			}
			tt.mockFn(svc)

			r := newTestRouter()
			r.POST("/collections/invitations/:invitationId/decline", withCustomerID(tt.customerID, h.declineInvitation))

			w := performRequest(r, http.MethodPost, "/collections/invitations/"+tt.invitationID+"/decline")
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
