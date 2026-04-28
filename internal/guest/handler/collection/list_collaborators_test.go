package collection

import (
	"errors"

	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

func TestCollaborators(t *testing.T) {
	t.Parallel()

	var (
		collectionID = gofakeit.UUID()
		customerID   = gofakeit.UUID()

		collaboratorCustomerID1 = "customer-uuid-1"
		collaboratorCustomerID2 = "customer-uuid-2"

		collaboratorUserName1 = "customer-1"
		collaboratorUserName2 = "customer-2"

		collaborators = []entity.Collaborator{
			{
				CollectionID: collectionID,
				CustomerID:   collaboratorCustomerID1,
				UserName:     collaboratorUserName1,
			},
			{
				CollectionID: collectionID,
				CustomerID:   collaboratorCustomerID2,
				UserName:     collaboratorUserName2,
			},
		}

		serviceErr = errors.New("unexpected service error")
	)

	tests := []struct {
		name string

		collectionID string
		customerID   any

		mockFn func(s *mockCollectionService)

		wantCode int
		wantBody any
	}{
		{
			name:         "success (authenticated)",
			collectionID: collectionID,
			customerID:   customerID,
			mockFn: func(s *mockCollectionService) {
				s.On("ListCollaborators", mock.Anything, collectionID, &customerID).
					Return(
						collaborators,
						nil,
					).
					Once()
			},
			wantCode: http.StatusOK,
			wantBody: listCollaboratorsResponse{
				Collaborators: []collaboratorResponse{
					{
						CustomerID: collaboratorCustomerID1,
						UserName:   collaboratorUserName1,
					},
					{
						CustomerID: collaboratorCustomerID2,
						UserName:   collaboratorUserName2,
					},
				},
			},
		},
		{
			name:         "success (unauthenticated)",
			collectionID: collectionID,
			customerID:   nil,
			mockFn: func(s *mockCollectionService) {
				s.On("ListCollaborators", mock.Anything, collectionID, (*string)(nil)).
					Return(
						collaborators,
						nil,
					).
					Once()
			},
			wantCode: http.StatusOK,
			wantBody: listCollaboratorsResponse{
				Collaborators: []collaboratorResponse{
					{
						CustomerID: collaboratorCustomerID1,
						UserName:   collaboratorUserName1,
					},
					{
						CustomerID: collaboratorCustomerID2,
						UserName:   collaboratorUserName2,
					},
				},
			},
		},
		{
			name:         "validation error - invalid collectionId uuid",
			collectionID: "invalid-collection-uuid",
			customerID:   customerID,
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
			name:         "optional customer id invalid type in context",
			collectionID: collectionID,
			customerID:   123,
			mockFn:       func(s *mockCollectionService) {},
			wantCode:     http.StatusInternalServerError,
			wantBody:     response.ErrorResponse{Message: internalErrMsg},
		},
		{
			name:         "service returns unknown error",
			collectionID: collectionID,
			customerID:   nil,
			mockFn: func(s *mockCollectionService) {
				s.On("ListCollaborators", mock.Anything, collectionID, (*string)(nil)).
					Return([]entity.Collaborator{}, serviceErr).
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
			h := &handler{service: svc}
			tt.mockFn(svc)

			r := newTestRouter()
			r.GET("/collections/:collectionId/collaborators", withCustomerID(tt.customerID, h.listCollaborators))

			w := performRequest(r, http.MethodGet, "/collections/"+tt.collectionID+"/collaborators")
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
