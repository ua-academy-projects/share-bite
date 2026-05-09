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

func TestRemoveCollaborator(t *testing.T) {
	t.Parallel()

	var (
		customerID       = gofakeit.UUID()
		targetCustomerID = gofakeit.UUID()
		collectionID     = gofakeit.UUID()

		serviceErr = errors.New("unexpected service error")
	)

	tests := []struct {
		name string

		collectionID     string
		targetCustomerID string
		customerID       any

		mockFn func(s *mockCollectionService)

		wantCode int
		wantBody any
	}{
		{
			name:             "success",
			collectionID:     collectionID,
			targetCustomerID: targetCustomerID,
			customerID:       customerID,
			mockFn: func(s *mockCollectionService) {
				s.On("RemoveCollaborator", mock.Anything, entity.RemoveCollaboratorInput{
					CollectionID:     collectionID,
					CustomerID:       customerID,
					TargetCustomerID: targetCustomerID,
				}).
					Return(nil).
					Once()
			},
			wantCode: http.StatusNoContent,
			wantBody: nil,
		},
		{
			name:             "validation error - invalid collection id",
			collectionID:     "invalid-collection-id",
			targetCustomerID: targetCustomerID,
			customerID:       customerID,
			mockFn:           func(s *mockCollectionService) {},
			wantCode:         http.StatusBadRequest,
			wantBody: response.ErrorResponse{
				Message: validationMsg,
				Details: []response.ErrorDetail{
					{Field: "collectionId", Message: "This field must be a valid UUID"},
				},
			},
		},
		{
			name:             "validation error - invalid target customer id",
			collectionID:     collectionID,
			targetCustomerID: "invalid-target-customer-id",
			customerID:       customerID,
			mockFn:           func(s *mockCollectionService) {},
			wantCode:         http.StatusBadRequest,
			wantBody: response.ErrorResponse{
				Message: validationMsg,
				Details: []response.ErrorDetail{
					{Field: "customerId", Message: "This field must be a valid UUID"},
				},
			},
		},
		{
			name:             "missing customer id in ctx",
			collectionID:     collectionID,
			targetCustomerID: targetCustomerID,
			customerID:       nil,
			mockFn:           func(s *mockCollectionService) {},
			wantCode:         http.StatusInternalServerError,
			wantBody:         response.ErrorResponse{Message: internalErrMsg},
		},
		{
			name:             "service unknown error",
			collectionID:     collectionID,
			targetCustomerID: targetCustomerID,
			customerID:       customerID,
			mockFn: func(s *mockCollectionService) {
				s.On("RemoveCollaborator", mock.Anything, entity.RemoveCollaboratorInput{
					CollectionID:     collectionID,
					CustomerID:       customerID,
					TargetCustomerID: targetCustomerID,
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
			h := &handler{
				service: svc,
			}
			tt.mockFn(svc)

			r := newTestRouter()
			r.DELETE("/collections/:collectionId/collaborators/:customerId", withCustomerID(tt.customerID, h.removeCollaborator))

			w := performRequest(r, http.MethodDelete, "/collections/"+tt.collectionID+"/collaborators/"+tt.targetCustomerID)
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
